package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ItemRepo implements the ItemRepository interface using PostgreSQL/GORM
type ItemRepo struct {
	db *gorm.DB
}

// NewItemRepo creates a new item repository
func NewItemRepo(db *gorm.DB) *ItemRepo {
	return &ItemRepo{
		db: db,
	}
}

// Create adds a new item to the database
func (r *ItemRepo) Create(ctx context.Context, item *models.Item) error {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		// Check for unique constraint violation
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves an item by ID
func (r *ItemRepo) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	result := r.db.WithContext(ctx).Preload(clause.Associations).First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &item, nil
}

// Update updates an existing item
func (r *ItemRepo) Update(ctx context.Context, item *models.Item) error {
	result := r.db.WithContext(ctx).Save(item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	return nil
}

// Delete removes an item from the database
func (r *ItemRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List retrieves items with pagination
func (r *ItemRepo) List(ctx context.Context, offset, limit int, productId *uint) ([]*models.Item, error) {
	var items []*models.Item
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit)
	if productId != nil {
		query = query.Where("product_id = ?", *productId)
	}
	result := query.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

// Count returns the total number of items
func (r *ItemRepo) Count(ctx context.Context, productId *uint) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Item{})
	if productId != nil {
		query = query.Where("product_id = ?", *productId)
	}
	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
