package services

import (
	"github.com/glebarez/sqlite"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&models.Assignment{},
		&models.AssignmentSubmission{},
		&models.Group{},
		&models.Course{},
		&models.Teacher{},
		&models.Student{},
		&models.Waitlist{},
		&models.Parent{},
		&models.ParentStudent{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}
