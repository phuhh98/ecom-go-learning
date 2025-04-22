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
func (h *CategoryHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminValidation gin.HandlerFunc) {
	categories := router.Group("/categories")
	{
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)
	}

	adminProtected := categories.Group("/")
	adminProtected.Use(authMiddleware, adminValidation)
	{
		adminProtected.PUT("/:id", h.Update)
		adminProtected.DELETE("/:id", h.Delete)
		adminProtected.POST("", h.Create)
		adminProtected.POST("/:id/products/:product_id", h.AddProduct)
		adminProtected.DELETE("/:id/products/:product_id", h.RemoveProduct)
		adminProtected.POST("/:id/subcategories/:subcategory_id", h.AddSubcategory)
		adminProtected.DELETE("/:id/subcategories/:subcategory_id", h.RemoveSubcategory)
	}
}

// Create handles category creation
// @Summary Create category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dtos.CreateCategoryDTO true "Category information"
// @Success 201 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories [post]
// @Security BearerAuth
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
// @Summary Get category by ID
// @Description Retrieve a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id} [get]
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
// @Summary Update category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body dtos.UpdateCategoryDTO true "Updated category information"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id} [put]
// @Security BearerAuth
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
// @Summary Delete category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dtos.DOCSuccessMessageResponse
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id} [delete]
// @Security BearerAuth
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
// @Summary List categories
// @Description Get a paginated list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dtos.DOCPaginatedResponse{data=[]dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories [get]
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

// AddProduct adds a product to a category
// @Summary Add product to category
// @Description Add a product to a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param product_id path int true "Product ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id}/products/{product_id} [post]
// @Security BearerAuth
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

// RemoveProduct removes a product from a category
// @Summary Remove product from category
// @Description Remove a product from a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param product_id path int true "Product ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id}/products/{product_id} [delete]
// @Security BearerAuth
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

// AddSubcategory adds a subcategory to a category
// @Summary Add subcategory
// @Description Add a subcategory to a parent category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Parent Category ID"
// @Param subcategory_id path int true "Subcategory ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id}/subcategories/{subcategory_id} [post]
// @Security BearerAuth
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
	// Subcategory Id should not be the same as category Id
	if categoryId == subcategoryId {
		c.Error(errors.NewBadRequestError("category ID and subcategory ID should not be the same"))
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

// RemoveSubcategory removes a subcategory from a category
// @Summary Remove subcategory
// @Description Remove a subcategory from a parent category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Parent Category ID"
// @Param subcategory_id path int true "Subcategory ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCCategoryResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /categories/{id}/subcategories/{subcategory_id} [delete]
// @Security BearerAuth
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
