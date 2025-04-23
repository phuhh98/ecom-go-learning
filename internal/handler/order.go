package handler

import (
	"ecom-go/internal/dtos"
	"ecom-go/internal/service"
	"net/http"
	"strconv"

	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	orders := router.Group("/orders")
	orders.Use(authMiddleware)
	{
		orders.GET("", h.ListOrders)
		orders.POST("", h.CreateOrder)
		orders.GET("/:orderId", h.GetOrderById)
		orders.PUT("/:orderId", h.UpdateOrder)
		orders.DELETE("/:orderId", h.DeleteOrder)
	}
}

// ListOrders handles retrieving orders with pagination
// @Summary List orders
// @Description Get a paginated list of orders (users see only their own orders, admins see all)
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dtos.DOCOrderListResponse
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /orders [get]
// @Security BearerAuth
func (h *OrderHandler) ListOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid page number", err))
		return
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid per_page number", err))
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	// Get user's role
	role, _ := c.Get("role")
	userRole, _ := role.(string)

	// Get orders based on role
	var orders interface{}
	var serviceErr error

	// If admin, can see all orders, otherwise only user's orders
	if userRole == "admin" {
		orders, serviceErr = h.orderService.ListOrders(c.Request.Context(), page, perPage)
	} else {
		orders, serviceErr = h.orderService.ListOrdersByUser(c.Request.Context(), userId.(uint), page, perPage)
	}

	if serviceErr != nil {
		c.Error(serviceErr)
		return
	}

	response.Success(c, http.StatusOK, orders)
}

// CreateOrder handles order creation
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dtos.CreateOrderDTO true "Order information"
// @Success 201 {object} dtos.DOCResponseWrapper{data=dtos.DOCOrderResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /orders [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var createOrderDTO dtos.CreateOrderDTO
	if err := c.ShouldBindJSON(&createOrderDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	// Get user ID from context (set by the auth middleware)
	userId, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	// Call service with context, user ID, and the DTO
	order, err := h.orderService.Create(c.Request.Context(), userId.(uint), createOrderDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, order)
}

// GetOrderById handles retrieving an order by ID
// @Summary Get order by ID
// @Description Retrieve an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path int true "Order ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCOrderResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /orders/{orderId} [get]
// @Security BearerAuth
func (h *OrderHandler) GetOrderById(c *gin.Context) {
	idStr := c.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid order ID", err))
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	// Get user's role
	role, _ := c.Get("role")
	userRole, _ := role.(string)

	order, err := h.orderService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	// Ensure the user can access this order (either it's theirs or they're an admin)
	if order.UserID != userId.(uint) && userRole != "admin" {
		c.Error(errors.NewForbiddenError("you don't have permission to access this order"))
		return
	}

	response.Success(c, http.StatusOK, order)
}

// UpdateOrder handles updating an order
// @Summary Update order
// @Description Update an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path int true "Order ID"
// @Param order body dtos.UpdateOrderDTO true "Updated order information"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCOrderResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /orders/{orderId} [put]
// @Security BearerAuth
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid order ID", err))
		return
	}

	var updateOrderDto dtos.UpdateOrderDTO
	if err := c.ShouldBindJSON(&updateOrderDto); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	// Get user's role
	role, _ := c.Get("role")
	userRole, _ := role.(string)

	// First check if the order exists and belongs to the user
	order, err := h.orderService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	// Ensure the user can update this order (either it's theirs or they're an admin)
	if order.UserID != userId.(uint) && userRole != "admin" {
		c.Error(errors.NewForbiddenError("you don't have permission to update this order"))
		return
	}

	updatedOrder, err := h.orderService.Update(c.Request.Context(), uint(id), updateOrderDto)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, updatedOrder)
}

// DeleteOrder handles deleting an order
// @Summary Delete order
// @Description Delete an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path int true "Order ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=object{message=string}}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /orders/{orderId} [delete]
// @Security BearerAuth
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Error(errors.NewBadRequestError("Invalid order ID", err))
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	// Get user's role
	role, _ := c.Get("role")
	userRole, _ := role.(string)

	// First check if the order exists and belongs to the user
	order, err := h.orderService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	// Ensure the user can delete this order (either it's theirs or they're an admin)
	if order.UserID != userId.(uint) && userRole != "admin" {
		c.Error(errors.NewForbiddenError("you don't have permission to delete this order"))
		return
	}

	if err := h.orderService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
