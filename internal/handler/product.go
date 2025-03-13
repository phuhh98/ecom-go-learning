package handler

import (
	"ecom-go/internal/dtos"
	"ecom-go/internal/service"
	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Register sets up routes for the product handler
func (h *ProductHandler) Register(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
		products.POST("/:id/items", h.CreateItemsForProduct)
		products.GET("/items/:itemId", h.GetItemById)
		products.GET("/:id/items", h.GetItemsForProduct)
		products.PUT("/items/:itemId", h.UpdateItem)
		products.DELETE("/items/:itemId", h.DeleteItem)
	}
}

// Create handles product creation
func (h *ProductHandler) Create(c *gin.Context) {
	var createProductDTO dtos.CreateProductDTO
	if err := c.ShouldBindJSON(&createProductDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.Create(c.Request.Context(), createProductDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, product)
}

// GetByID handles retrieving a product by ID
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

// Update handles updating a product
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	var updateProductDTO dtos.UpdateProductDTO
	if err := c.ShouldBindJSON(&updateProductDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.Update(c.Request.Context(), uint(id), updateProductDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

// Delete handles deleting a product
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// List handles retrieving products with pagination
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	products, total, err := h.productService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, products, page, pageSize, total)
}

func (h *ProductHandler) CreateItemsForProduct(c *gin.Context) {
	var createItemDto dtos.CreateItemsDTO
	if err := c.ShouldBindJSON(&createItemDto); err != nil {

		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid request body"))
		return
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	createItemDto.ProductId = uint(productId)

	items, err := h.productService.CreateItemsForProduct(c.Request.Context(), createItemDto)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithTotalCount(c, http.StatusCreated, items, int64(createItemDto.Quantity))
}

func (h *ProductHandler) GetItemById(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid item ID"))
		return
	}

	item, err := h.productService.GetItemById(c.Request.Context(), uint(itemId))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, item)
}

func (h *ProductHandler) GetItemsForProduct(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.Error(errors.NewBadRequestError("invalid page parameter"))
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil || pageSize <= 0 {
		c.Error(errors.NewBadRequestError("invalid per_page parameter"))
		return
	}

	productId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	items, total, err := h.productService.GetItemsForProduct(c.Request.Context(), page, pageSize, uint(productId))
	if err != nil {
		c.Error(err)
		return
	}
	response.SuccessWithPagination(c, http.StatusOK, items, page, pageSize, total)
}

func (h *ProductHandler) UpdateItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid item ID"))
		return
	}
	var updateItemDto dtos.UpdateItemDto
	if err := c.ShouldBindJSON(&updateItemDto); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid request body"))
		return
	}

	updateItemDto.ItemId = uint(itemId)
	item, err := h.productService.UpdateItem(c.Request.Context(), updateItemDto)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, item)
}

func (h *ProductHandler) DeleteItem(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid item ID"))
		return
	}
	err = h.productService.DeleteItem(c.Request.Context(), uint(itemId))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "item deleted successfully")
}
