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

// CategoryHandler handles HTTP requests related to categorys
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Register sets up routes for the category handler
func (h *CategoryHandler) Register(router *gin.RouterGroup,authMiddleware gin.HandlerFunc) {
	categories := router.Group("/categories")
	{
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)
	}

	protected := categories.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("", h.Create)
		protected.PUT("/:id", h.Update)
		protected.DELETE("/:id", h.Delete)
		protected.POST("/:id/products/:product_id", h.AddProduct)
		protected.DELETE("/:id/products/:product_id", h.RemoveProduct)
		protected.POST("/:id/subcategories/:subcategory_id", h.AddSubcategory)
		protected.DELETE("/:id/subcategories/:subcategory_id", h.RemoveSubcategory)
	}
}

// Create handles category creation
func (h *CategoryHandler) Create(c *gin.Context) {
	var createCategoryDTO dtos.CreateCategoryDTO
	if err := c.ShouldBindJSON(&createCategoryDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), createCategoryDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// GetByID handles retrieving a category by ID
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, category)
}

// Update handles updating a category
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}

	var updateCategoryDTO dtos.UpdateCategoryDTO
	if err := c.ShouldBindJSON(&updateCategoryDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), uint(id), updateCategoryDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, category)
}

// Delete handles deleting a category
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// List handles retrieving categorys with pagination
func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	categorys, total, err := h.categoryService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, categorys, page, pageSize, total)
}

func (h *CategoryHandler) AddProduct(c *gin.Context) {
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}

	productId, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	updatedCategory, err := h.categoryService.AddProductToCategory(c.Request.Context(), uint(categoryId), uint(productId))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) RemoveProduct(c *gin.Context) {
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}
	productId, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid product ID"))
		return
	}

	updatedCategory, err := h.categoryService.RemoveProductFromCategory(c.Request.Context(), uint(categoryId), uint(productId))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) AddSubcategory(c *gin.Context) {
	// validate category id and subcategory id
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}
	subcategoryId, err := strconv.ParseUint(c.Param("subcategory_id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid subcategory ID"))
		return
	}

	// update category subcategories with an update
	updatedCategory, err := h.categoryService.AddSubcategoryToCategory(c.Request.Context(), uint(categoryId), uint(subcategoryId))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) RemoveSubcategory(c *gin.Context) {
	// validate category id and subcategory id
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid category ID"))
		return
	}
	subcategoryId, err := strconv.ParseUint(c.Param("subcategory_id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid subcategory ID"))
		return
	}

	// update category subcategories with an update
	updatedCategory, err := h.categoryService.RemoveSubcategoryFromCategory(c.Request.Context(), uint(categoryId), uint(subcategoryId))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, updatedCategory)
}
