package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/softclub-go-0-0/crm-service/pkg/config"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DBInit initializes database connection with configuration
func DBInit(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger based on app log level
	gormLogLevel := gormlogger.Silent
	if cfg.IsDevelopment() {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Prepare statements for better performance
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database")

	// Run auto migrations (only in development)
	if cfg.IsDevelopment() {
		logger.Info("Running database auto-migrations (development mode)")
		if err := runMigrations(db); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

// runMigrations runs GORM auto-migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Teacher{},
		&models.Course{},
		&models.Student{},
		&models.Group{},
		&models.Timetable{},
		&models.Attendance{},
		&models.Grade{},
		// Phase 1: RBAC & Infrastructure
		&models.User{},
		&models.Permission{},
		&models.RolePermission{},
		&models.AuditLog{},
		&models.Session{},
	)
}

// HealthCheck checks if database is healthy
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// GetStats returns database connection pool statistics
func GetStats(db *gorm.DB) sql.DBStats {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}
	}
	return sqlDB.Stats()
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
