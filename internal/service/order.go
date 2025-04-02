package service

import (
	"context"
	"ecom-go/internal/dtos"
	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

type OrderService struct {
	orderRepo      repository.OrderRepository
	productService *ProductService
	userService  *UserService
}

func NewOrderService(orderRepo repository.OrderRepository, productService *ProductService, userService *UserService) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		userService: userService,
	}
}

// func (s *OrderService) ListOrders(page int, perPage int) ([]models.Order, error) {
// 	return s.orderRepo.ListOrders(page, perPage)
// }

func (s *OrderService) Create(ctx context.Context, createOrderDTO dtos.CreateOrderDTO) (*models.Order, error) {
	// validate user exists
	_, err := s.userService.GetByID(ctx, createOrderDTO.UserID)
	if err != nil {
		return nil, appError.NewNotFoundError("user not found", err)
	}
	
	order := &models.Order{
		UserID: createOrderDTO.UserID,
		Status: "created",
		Items:  []*models.OrderItem{},
	}
	// Check list of product whether products exist
	for _, item := range createOrderDTO.Items {
		product, err := s.productService.GetByID(nil, item.ProductID)
		if err != nil {
			return nil, err
		}

		order.Items = append(order.Items, &models.OrderItem{
			Product:  *product,
			Quantity: item.Quantity,
		})
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, appError.NewServerError("error creating order", err)
	}

	return order, nil
}

// func (s *OrderService) GetOrderById(id uint) (*models.Order, error) {
// 	order, err := s.orderRepo.GetOrderById(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get items associated with this order
// 	items, err := s.orderRepo.GetOrderItems(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Attach items to the order
// 	order.Items = items
// 	return order, nil
// }

// func (s *OrderService) UpdateOrder(order *models.Order) error {
// 	return s.orderRepo.UpdateOrder(order)
// }

// func (s *OrderService) DeleteOrder(id uint) error {
// 	// Get items associated with this order
// 	items, err := s.orderRepo.GetOrderItems(id)
// 	if err != nil {
// 		return err
// 	}

// 	// Update items to disassociate from order and mark as available
// 	for _, item := range items {
// 		item.OrderID = nil
// 		item.Status = "available"
// 		if err := s.itemRepo.Update(nil, item); err != nil {
// 			return err
// 		}
// 	}

// 	return s.orderRepo.DeleteOrder(id)
// }
