package repository

import (
	"context"

	"ecom-go/internal/models"
)

type CategoryRepository interface {
	// Add new Category to db
	Create(ctx context.Context, user *models.Category) error

	// Get Category by Id
	GetByID(ctx context.Context, id uint) (*models.Category, error)

	GetByName(ctx context.Context, name string) (*models.Category, error)

	// Update updates an existing Category
	Update(ctx context.Context, user *models.Category) error

	// Remove a Category from the database
	Delete(ctx context.Context, id uint) error

	// List retrieves Categorys with pagination
	List(ctx context.Context, offset, limit int) ([]*models.Category, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
