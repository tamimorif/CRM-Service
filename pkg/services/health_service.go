package services

import (
	"context"
	"runtime"
	"time"

	"gorm.io/gorm"
)

// HealthService defines the interface for health operations
type HealthService interface {
	CheckHealth(ctx context.Context) HealthStatus
	CheckDatabase(ctx context.Context) DatabaseHealth
}

type healthService struct {
	db *gorm.DB
}

// NewHealthService creates a new health service
func NewHealthService(db *gorm.DB) HealthService {
	return &healthService{db: db}
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Database  DatabaseHealth    `json:"database"`
	System    SystemInfo        `json:"system"`
	Services  map[string]string `json:"services"`
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

func (s *healthService) CheckHealth(ctx context.Context) HealthStatus {
	dbHealth := s.CheckDatabase(ctx)

	sysInfo := SystemInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
	}

	services := map[string]string{
		"auth_service": "healthy", // Placeholder
	}

	status := "healthy"
	if dbHealth.Status != "healthy" {
		status = "unhealthy"
	}

	return HealthStatus{
		Status:    status,
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Database:  dbHealth,
		System:    sysInfo,
		Services:  services,
	}
}

func (s *healthService) CheckDatabase(ctx context.Context) DatabaseHealth {
	start := time.Now()

	sqlDB, err := s.db.DB()
	if err != nil {
		return DatabaseHealth{
			Status:       "unhealthy",
			ResponseTime: time.Since(start),
		}
	}

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
