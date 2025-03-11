package repository

import (
	"context"

	"ecom-go/internal/models"
)

type ProductRepository interface {
	// Add new Product to db
	Create(ctx context.Context, user *models.Product) error

	// Get Product by Id
	GetByID(ctx context.Context, id uint) (*models.Product, error)

	GetByCode(ctx context.Context, code string) (*models.Product, error)

	// Update updates an existing Product
	Update(ctx context.Context, user *models.Product) error

	// Remove a Product from the database
	Delete(ctx context.Context, id uint) error

	// List retrieves Products with pagination
	List(ctx context.Context, offset, limit int) ([]*models.Product, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
