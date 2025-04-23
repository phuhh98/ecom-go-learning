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

func (r *OrderRepo) List(ctx context.Context, page int, perPage int) ([]models.Order, error) {
	var orders []models.Order
	offset := (page - 1) * perPage

	result := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("Address").
		Offset(offset).
		Limit(perPage).
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *OrderRepo) ListByUser(ctx context.Context, userID uint, page int, perPage int) ([]models.Order, error) {
	var orders []models.Order
	offset := (page - 1) * perPage

	result := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("Address").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(perPage).
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *OrderRepo) Create(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	return result.Error
}

func (r *OrderRepo) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	result := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("Address").
		First(&order, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepo) Update(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Save(order)
	return result.Error
}

func (r *OrderRepo) Delete(ctx context.Context, id uint) error {
	// Handling transaction to ensure OrderItems are deleted
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete associated OrderItems first
		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}

		// Then delete the order
		return tx.Delete(&models.Order{}, id).Error
	})
}
