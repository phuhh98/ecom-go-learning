package handler

import (
	"ecom-go/internal/dtos"
	"ecom-go/internal/service"
	"net/http"

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
		// orders.GET("", h.ListOrders)
		orders.POST("", h.CreateOrder)
		// orders.GET("/:orderId", h.GetOrderById)
		// orders.PUT("/:orderId", h.UpdateOrder)
		// orders.DELETE("/:orderId", h.DeleteOrder)
	}
}

// func (h *OrderHandler) ListOrders(c *gin.Context) {
// 	pageStr := c.DefaultQuery("page", "1")
// 	perPageStr := c.DefaultQuery("per_page", "10")

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
// 		return
// 	}

// 	perPage, err := strconv.Atoi(perPageStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page number"})
// 		return
// 	}

// 	orders, err := h.orderService.ListOrders(page, perPage)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, orders)
// }

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

	userId, ok := c.Get("user_id")
	if !ok {
		c.Error(errors.NewBadRequestError("user ID not found in context"))
		return
	}
	createOrderDTO.UserID = userId.(uint)

	order, err := h.orderService.Create(c.Request.Context(), createOrderDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, order)
}

// func (h *OrderHandler) GetOrderById(c *gin.Context) {
// 	idStr := c.Param("orderId")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
// 		return
// 	}

// 	order, err := h.orderService.GetOrderById(uint(id))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, order)
// }

// func (h *OrderHandler) UpdateOrder(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
// 		return
// 	}

// 	var updateOrderDto dtos.UpdateOrderDTO
// 	if err := c.ShouldBindJSON(&updateOrderDto); err != nil {
// 		if ve, ok := err.(validator.ValidationErrors); ok {
// 			c.Error(errors.NewValidationError(&ve))
// 			return
// 		}

// 		c.Error(errors.NewBadRequestError("invalid input", err))
// 		return
// 	}

// 	order, err := h.orderService.UpdateOrder(c.Request.Context(), uint(id), updateOrderDto)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	response.Success(c, http.StatusOK, order)
// }

// func (h *OrderHandler) DeleteOrder(c *gin.Context) {
// 	idStr := c.Param("orderId")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
// 		return
// 	}

// 	if err := h.orderService.DeleteOrder(uint(id)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
// }
