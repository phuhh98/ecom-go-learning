package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CategoryRepo implements the CategoryRepository interface using PostgreSQL/GORM
type CategoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepo creates a new category repository
func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

// TODO: Template - need to update
// Create adds a new category to the database
func (r *CategoryRepo) Create(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Create(category)
	if result.Error != nil {
		// Check for unique constraint violation
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepo) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).Preload(clause.Associations).Preload("SubCategories").First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &category, nil
}

// GetByName retrieves a category by name - unique
func (r *CategoryRepo) GetByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).Where("name = ?", name).Preload(clause.Associations).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepo) Update(ctx context.Context, category *models.Category) error {
	// Use `Association().Replace` to synchronize the relationship.
	if err := r.db.WithContext(ctx).Model(category).Association("Products").Replace(category.Products); err != nil {
		return err // Handle the error appropriately
	}

	if err := r.db.WithContext(ctx).Model(category).Association("SubCategories").Replace(category.SubCategories); err != nil {
		return err // Handle the error appropriately
	}

	// Then, update other attributes
	if err := r.db.WithContext(ctx).Model(category).Select("*").Omit("Products", "SubCategories").Updates(category).Error; err != nil {
		return err // Handle the error appropriately
	}

	return nil
}

// Delete removes a category from the database
func (r *CategoryRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List retrieves categorys with pagination
func (r *CategoryRepo) List(ctx context.Context, offset, limit int) ([]*models.Category, error) {
	var categorys []*models.Category
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Preload(clause.Associations).Find(&categorys)
	if result.Error != nil {
		return nil, result.Error
	}
	return categorys, nil
}

// Count returns the total number of categorys
func (r *CategoryRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Category{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
