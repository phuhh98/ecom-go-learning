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
	userService    *UserService
	addressService *AddressService
}

func NewOrderService(orderRepo repository.OrderRepository, productService *ProductService, userService *UserService, addressService *AddressService) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		userService:    userService,
		addressService: addressService,
	}
}

func (s *OrderService) ListOrders(ctx context.Context, page int, perPage int) ([]models.Order, error) {
	orders, err := s.orderRepo.List(ctx, page, perPage)
	if err != nil {
		return nil, appError.NewServerError("error listing orders", err)
	}
	return orders, nil
}

func (s *OrderService) ListOrdersByUser(ctx context.Context, userID uint, page int, perPage int) ([]models.Order, error) {
	orders, err := s.orderRepo.ListByUser(ctx, userID, page, perPage)
	if err != nil {
		return nil, appError.NewServerError("error listing user orders", err)
	}
	return orders, nil
}

func (s *OrderService) Create(ctx context.Context, userID uint, createOrderDTO dtos.CreateOrderDTO) (*models.Order, error) {
	// Error handling is consistent with typed errors
	_, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, appError.NewNotFoundError("user not found", err)
	}

	// validate address exists and belongs to the user
	address, err := s.addressService.GetByID(ctx, createOrderDTO.AddressID)
	if err != nil {
		return nil, appError.NewNotFoundError("address not found", err)
	}

	// Authorization check using typed error
	if address.UserID != userID {
		return nil, appError.NewUnauthorizedError("address does not belong to this user")
	}

	order := &models.Order{
		UserID:    userID,
		Status:    "created",
		Items:     []*models.OrderItem{},
		AddressID: createOrderDTO.AddressID,
		Address:   *address,
	}
	// Check list of product whether products exist
	for _, item := range createOrderDTO.Items {
		product, err := s.productService.GetByID(ctx, item.ProductID)
		if err != nil {
			// Pass through typed errors from the product service
			// This maintains the specific error type (NotFound vs Server error)
			return nil, err
		}

		order.Items = append(order.Items, &models.OrderItem{
			Product:  *product,
			Quantity: item.Quantity,
		})
	}

	// Server error wrapping is consistent
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, appError.NewServerError("error creating order", err)
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, appError.NewNotFoundError("order not found", err)
	}
	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id uint, updateOrderDTO dtos.UpdateOrderDTO) (*models.Order, error) {
	// First get the existing order
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, appError.NewNotFoundError("order not found", err)
	}

	// Update the order status
	order.Status = updateOrderDTO.Status

	// Update the order in the database
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, appError.NewServerError("error updating order", err)
	}

	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, id uint) error {
	// Check if order exists first
	_, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return appError.NewNotFoundError("order not found", err)
	}

	// Delete the order
	if err := s.orderRepo.Delete(ctx, id); err != nil {
		return appError.NewServerError("error deleting order", err)
	}

	return nil
}
