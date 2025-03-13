package repository

import (
	"context"

	"ecom-go/internal/models"
)

type ItemRepository interface {
	// Add new item to db
	Create(ctx context.Context, item *models.Item) error

	// Get Item by Id
	GetByID(ctx context.Context, id uint) (*models.Item, error)

	// Update updates an existing Item
	Update(ctx context.Context, item *models.Item) error

	// Remove an Item from the database
	Delete(ctx context.Context, id uint) error

	// List retrieves Items with pagination
	List(ctx context.Context, offset, limit int, productId *uint) ([]*models.Item, error)

	// Count returns the total number of Items
	Count(ctx context.Context, productId *uint) (int64, error)
}
