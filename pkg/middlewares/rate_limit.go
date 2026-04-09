package middlewares

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"golang.org/x/time/rate"
)

// ipEntry holds a rate limiter and its last access time for TTL-based cleanup
type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter holds rate limiters for each IP with TTL-based cleanup
type IPRateLimiter struct {
	ips map[string]*ipEntry
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
	ttl time.Duration
}

// NewIPRateLimiter creates a new IP rate limiter with TTL-based cleanup
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*ipEntry),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
		ttl: 10 * time.Minute, // Default TTL of 10 minutes
	}

	// Start a background routine to clean up expired entries
	go i.cleanup()

	return i
}

// AddIP creates a new rate limiter for an IP
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = &ipEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP and updates last seen
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	entry, exists := i.ips[ip]
	i.mu.RUnlock()

	if exists {
		// Update last seen time
		i.mu.Lock()
		entry.lastSeen = time.Now()
		i.mu.Unlock()
		return entry.limiter
	}

	return i.AddIP(ip)
}

// cleanup removes expired entries based on TTL to prevent memory leaks
func (i *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for ip, entry := range i.ips {
			if now.Sub(entry.lastSeen) > i.ttl {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimit middleware limits requests based on IP
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			errors.HandleError(c, errors.New(errors.ErrCodeTooManyRequests, "Rate limit exceeded"))
			c.Abort()
			return
		}
		c.Next()
	}
}
