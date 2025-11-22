package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository is a generic interface for data access
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context, offset, limit int, query interface{}, args ...interface{}) ([]T, int64, error)
	GetDB() *gorm.DB
}

// GormRepository implements Repository using GORM
type GormRepository[T any] struct {
	db *gorm.DB
}

// NewGormRepository creates a new GormRepository
func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

// Create inserts a new entity
func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update updates an existing entity
func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete deletes an entity by ID
func (r *GormRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

// GetByID retrieves an entity by ID
func (r *GormRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetAll retrieves all entities with pagination
func (r *GormRepository[T]) GetAll(ctx context.Context, offset, limit int, query interface{}, args ...interface{}) ([]T, int64, error) {
	var entities []T
	var total int64

	db := r.db.WithContext(ctx).Model(new(T))

	if query != nil {
		db = db.Where(query, args...)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// GetDB returns the underlying GORM DB
func (r *GormRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// WithTx runs a function within a transaction
func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(fn)
}
