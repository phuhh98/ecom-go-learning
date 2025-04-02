package handler

import (
	"ecom-go/internal/dtos"
	"net/http"
	"strconv"

	"ecom-go/internal/service"
	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, tokenService *service.TokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// Register sets up routes for the user handler
func (h *UserHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	{
		users.POST("", h.Create)
	}

	protected := users.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("", h.List)
		protected.GET("/:id", h.GetByID)
		protected.PUT("/:id", h.Update)
		protected.DELETE("/:id", h.Delete)
	}

}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var createUserDTO dtos.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}
		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), createUserDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// GetByID handles retrieving a user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Update handles updating a user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	var updateUserDTO dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), updateUserDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Delete handles deleting a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// List handles retrieving users with pagination
func (h *UserHandler) List(c *gin.Context) {
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

	users, total, err := h.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, users, page, pageSize, total)
}
