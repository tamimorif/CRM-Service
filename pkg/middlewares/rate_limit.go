package middlewares

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"golang.org/x/time/rate"
)

// IPRateLimiter holds rate limiters for each IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter creates a new IP rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	// Start a background routine to clean up old entries
	go i.cleanup()

	return i
}

// AddIP creates a new rate limiter for an IP if it doesn't exist
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// cleanup removes old entries to prevent memory leaks
// In a real production system, you might use Redis for this
func (i *IPRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		i.mu.Lock()
		// Simple cleanup: reset map if it gets too big
		// A better approach would be to track last access time
		if len(i.ips) > 10000 {
			i.ips = make(map[string]*rate.Limiter)
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
