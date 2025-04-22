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
// @Summary Create user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.CreateUserDTO true "User registration information"
// @Success 201 {object} dtos.DOCResponseWrapper{data=dtos.DOCUserDetailResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 409 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError} "Email already exists"
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /users [post]
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
// @Summary Get user by ID
// @Description Retrieve a user's information by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCUserDetailResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /users/{id} [get]
// @Security BearerAuth
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
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dtos.UpdateUserDTO true "Updated user information"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCUserDetailResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /users/{id} [put]
// @Security BearerAuth
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
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.DOCSuccessMessageResponse
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 403 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /users/{id} [delete]
// @Security BearerAuth
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
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dtos.DOCPaginatedResponse{data=[]dtos.DOCUserListItem}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /users [get]
// @Security BearerAuth
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
