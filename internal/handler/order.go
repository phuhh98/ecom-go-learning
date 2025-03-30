package handler

import (
	"ecom-go/internal/dtos"
	"ecom-go/internal/models"
	"ecom-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *OrderHandler) ListOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid per_page number"})
		return
	}

	orders, err := h.orderService.ListOrders(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderDto dtos.OrderDTO
	if err := c.BindJSON(&orderDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		UserID: orderDto.UserID,
		Status: orderDto.Status,
	}

	if err := h.orderService.CreateOrder(&order, orderDto.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the complete order with items
	completeOrder, err := h.orderService.GetOrderById(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, completeOrder)
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	idStr := c.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrderById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	idStr := c.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var orderDto dtos.OrderDTO
	if err := c.BindJSON(&orderDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		ID:     uint(id),
		UserID: orderDto.UserID,
		Status: orderDto.Status,
	}

	if err := h.orderService.UpdateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the complete order with items
	completeOrder, err := h.orderService.GetOrderById(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, completeOrder)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("orderId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderService.DeleteOrder(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
