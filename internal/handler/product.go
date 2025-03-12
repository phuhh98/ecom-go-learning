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
	}
}

// Create handles product creation
func (h *ProductHandler) Create(c *gin.Context) {
	var createProductDTO dtos.CreateProductDTO
	if err := c.ShouldBindJSON(&createProductDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.Error(c, errors.NewValidationError(&ve))
			return
		}

		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.Create(c.Request.Context(), createProductDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, product)
}

// GetByID handles retrieving a product by ID
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

// Update handles updating a product
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	var updateProductDTO dtos.UpdateProductDTO
	if err := c.ShouldBindJSON(&updateProductDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.Error(c, errors.NewValidationError(&ve))
			return
		}

		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	product, err := h.productService.Update(c.Request.Context(), uint(id), updateProductDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

// Delete handles deleting a product
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product ID"))
		return
	}

	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
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
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, products, page, pageSize, total)
}
