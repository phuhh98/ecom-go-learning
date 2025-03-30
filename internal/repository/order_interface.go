package repository

import "ecom-go/internal/models"

type OrderRepository interface {
	ListOrders(page int, perPage int) ([]models.Order, error)
	CreateOrder(order *models.Order) error
	GetOrderById(id uint) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id uint) error
	GetOrderItems(orderID uint) ([]*models.Item, error)
}
