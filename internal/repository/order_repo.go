package repository

import (
	"context"
	"ecom-go/internal/models"

	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// func (r *OrderRepo) ListOrders(page int, perPage int) ([]models.Order, error) {
// 	var orders []models.Order
// 	offset := (page - 1) * perPage
// 	result := r.db.Offset(offset).Limit(perPage).Find(&orders)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	// For each order, load its items
// 	for i := range orders {
// 		items, err := r.GetOrderItems(orders[i].ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		orders[i].Items = items
// 	}

// 	return orders, nil
// }

func (r *OrderRepo) Create(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	return result.Error
}

// func (r *OrderRepo) GetOrderById(id uint) (*models.Order, error) {
// 	var order models.Order
// 	result := r.db.First(&order, id)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &order, nil
// }

// func (r *OrderRepo) UpdateOrder(order *models.Order) error {
// 	result := r.db.Save(order)
// 	return result.Error
// }

// func (r *OrderRepo) DeleteOrder(id uint) error {
// 	result := r.db.Delete(&models.Order{}, id)
// 	return result.Error
// }

// func (r *OrderRepo) GetOrderItems(orderID uint) ([]*models.Item, error) {
// 	var items []*models.Item
// 	// Use preloading to get items with their associated products
// 	result := r.db.Preload("Product").Where("order_id = ?", orderID).Find(&items)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return items, nil
// }
