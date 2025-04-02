package repository

import (
	"context"
	"ecom-go/internal/models"
)

type OrderRepository interface {
	// List(ctx context.Context, page int, perPage int) ([]models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	// GetById(ctx context.Context, id uint) (*models.Order, error)
	// Update( ctx context.Context, order *models.Order) error
	// Delete(ctx context.Context, id uint) error
	// GetItems(ctx context.Context, orderID uint) ([]*models.Item, error)
}
