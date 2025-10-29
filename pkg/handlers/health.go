package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"net/http"
	"runtime"
	"time"
)

type HealthStatus struct {
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      string            `json:"uptime"`
	Database    DatabaseHealth    `json:"database"`
	System      SystemInfo        `json:"system"`
	Services    map[string]string `json:"services"`
}

type DatabaseHealth struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
}

type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
}

var startTime = time.Now()

func (h *handler) HealthCheck(c *gin.Context) {
	// Check database health
	dbHealth := h.checkDatabaseHealth()

	// Get system info
	sysInfo := SystemInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
	}

	// Check external services
	services := map[string]string{
		"auth_service": "healthy", // You can implement actual health check for auth service
	}

	// Overall status
	status := "healthy"
	if dbHealth.Status != "healthy" {
		status = "unhealthy"
	}

	health := HealthStatus{
		Status:    status,
		Version:   "1.0.0", // You can get this from build info or env
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Database:  dbHealth,
		System:    sysInfo,
		Services:  services,
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

func (h *handler) checkDatabaseHealth() DatabaseHealth {
	start := time.Now()
	
	// Get underlying sql.DB from GORM
	sqlDB, err := h.DB.DB()
	if err != nil {
		return DatabaseHealth{
			Status:       "unhealthy",
			ResponseTime: time.Since(start),
		}
	}

	// Ping the database
	if err := sqlDB.Ping(); err != nil {
		return DatabaseHealth{
			Status:       "unhealthy", 
			ResponseTime: time.Since(start),
		}
	}

	return DatabaseHealth{
		Status:       "healthy",
		ResponseTime: time.Since(start),
	}
}

func (h *handler) ReadinessProbe(c *gin.Context) {
	// Check if the application is ready to serve requests
	dbHealth := h.checkDatabaseHealth()
	
	if dbHealth.Status == "healthy" {
		helpers.SuccessResponse(c, gin.H{"ready": true}, "Application is ready")
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready": false,
			"reason": "Database is not ready",
		})
	}
}

func (h *handler) LivenessProbe(c *gin.Context) {
	// Simple liveness check - if we can respond, we're alive
	helpers.SuccessResponse(c, gin.H{"alive": true}, "Application is alive")
}