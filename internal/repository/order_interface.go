package repository

import (
	"context"
	"ecom-go/internal/models"
)

type OrderRepository interface {
	List(ctx context.Context, page int, perPage int) ([]models.Order, error)
	ListByUser(ctx context.Context, userID uint, page int, perPage int) ([]models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id uint) error
}
