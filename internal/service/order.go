package service

import (
	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	"errors"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	itemRepo  repository.ItemRepository
}

func NewOrderService(orderRepo repository.OrderRepository, itemRepo repository.ItemRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

func (s *OrderService) ListOrders(page int, perPage int) ([]models.Order, error) {
	return s.orderRepo.ListOrders(page, perPage)
}

func (s *OrderService) CreateOrder(order *models.Order, itemIDs []uint) error {
	// Create the order first
	if err := s.orderRepo.CreateOrder(order); err != nil {
		return err
	}

	// Associate items with the order
	for _, itemID := range itemIDs {
		item, err := s.itemRepo.GetByID(nil, itemID)
		if err != nil {
			return err
		}

		// Check if item is available
		if item.Status != "available" {
			return errors.New("item is not available")
		}

		// Update item to associate with order and mark as sold
		item.OrderID = &order.ID
		item.Status = "sold"
		if err := s.itemRepo.Update(nil, item); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) GetOrderById(id uint) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderById(id)
	if err != nil {
		return nil, err
	}

	// Get items associated with this order
	items, err := s.orderRepo.GetOrderItems(id)
	if err != nil {
		return nil, err
	}

	// Attach items to the order
	order.Items = items
	return order, nil
}

func (s *OrderService) UpdateOrder(order *models.Order) error {
	return s.orderRepo.UpdateOrder(order)
}

func (s *OrderService) DeleteOrder(id uint) error {
	// Get items associated with this order
	items, err := s.orderRepo.GetOrderItems(id)
	if err != nil {
		return err
	}

	// Update items to disassociate from order and mark as available
	for _, item := range items {
		item.OrderID = nil
		item.Status = "available"
		if err := s.itemRepo.Update(nil, item); err != nil {
			return err
		}
	}

	return s.orderRepo.DeleteOrder(id)
}
