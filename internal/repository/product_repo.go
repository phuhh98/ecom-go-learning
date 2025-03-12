package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductRepo implements the ProductRepository interface using PostgreSQL/GORM
type ProductRepo struct {
	db *gorm.DB
}

// NewProductRepo creates a new product repository
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// TODO: Template - need to update
// Create adds a new product to the database
func (r *ProductRepo) Create(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		// Check for unique constraint violation
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepo) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	result := r.db.WithContext(ctx).Preload(clause.Associations).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &product, nil
}

// GetByCode retrieves a product by product code
func (r *ProductRepo) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	var product models.Product
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepo) Update(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	return nil
}

// Delete removes a product from the database
func (r *ProductRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List retrieves products with pagination
func (r *ProductRepo) List(ctx context.Context, offset, limit int) ([]*models.Product, error) {
	var products []*models.Product
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// Count returns the total number of products
func (r *ProductRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
