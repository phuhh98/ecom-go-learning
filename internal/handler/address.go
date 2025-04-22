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

// @title E-Commerce API
// @description API for the E-Commerce application
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// AddressHandler handles HTTP requests related to addresss
type AddressHandler struct {
	addressService *service.AddressService
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// Register sets up routes for the address handler
func (h *AddressHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	addresses := router.Group("/addresses")
	addresses.Use(authMiddleware)
	{
		addresses.POST("", h.Create)
		addresses.GET("", h.List)
		addresses.GET("/:id", h.GetByID)
		addresses.PUT("/:id", h.Update)
		addresses.DELETE("/:id", h.Delete)
		addresses.GET("/user/:userId", h.ListByUserID)
	}
}

// Create handles address creation
// @Summary Create address
// @Description Create a new address for the authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Param address body dtos.CreateAddressDTO true "Address information"
// @Success 201 {object} dtos.DOCResponseWrapper{data=dtos.DOCAddressResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /addresses [post]
// @Security BearerAuth
func (h *AddressHandler) Create(c *gin.Context) {
	// Extract user ID from the context (assuming it's set by the auth middleware)
	userID, ok := c.Get("user_id")

	if !ok {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	var createAddressDTO dtos.CreateAddressDTO
	if err := c.ShouldBindJSON(&createAddressDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	createAddressDTO.UserID = userID.(uint)

	address, err := h.addressService.Create(c.Request.Context(), createAddressDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, address)
}

// GetByID handles retrieving a address by ID
// @Summary Get address by ID
// @Description Retrieve an address by its ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCAddressResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 404 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /addresses/{id} [get]
// @Security BearerAuth
func (h *AddressHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid address ID"))
		return
	}

	address, err := h.addressService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// Update handles updating a address
// @Summary Update address
// @Description Update an existing address
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param address body dtos.UpdateAddressDTO true "Updated address information"
// @Success 200 {object} dtos.DOCAddressResponse
// @Failure 400 {object} dtos.DOCErrorResponse
// @Failure 401 {object} dtos.DOCErrorResponse
// @Failure 404 {object} dtos.DOCErrorResponse
// @Failure 500 {object} dtos.DOCErrorResponse
// @Router /addresses/{id} [put]
// @Security BearerAuth
func (h *AddressHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid address ID"))
		return
	}

	var updateAddressDTO dtos.UpdateAddressDTO
	if err := c.ShouldBindJSON(&updateAddressDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	// user id from context
	userID, ok := c.Get("user_id")
	if !ok {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	address, err := h.addressService.Update(c.Request.Context(), uint(id), updateAddressDTO, userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// Delete handles deleting a address
// @Summary Delete address
// @Description Delete an address by ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} dtos.DOCSuccessResponse
// @Failure 400 {object} dtos.DOCErrorResponse
// @Failure 401 {object} dtos.DOCErrorResponse
// @Failure 404 {object} dtos.DOCErrorResponse
// @Failure 500 {object} dtos.DOCErrorResponse
// @Router /addresses/{id} [delete]
// @Security BearerAuth
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid address ID"))
		return
	}

	// user id from context
	userID, ok := c.Get("user_id")
	if !ok {
		c.Error(errors.NewUnauthorizedError("user ID not found in context"))
	}

	if err := h.addressService.Delete(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// ListByUserID handles retrieving addresses by user ID
// @Summary List addresses by user ID
// @Description Get a paginated list of addresses for a specific user
// @Tags addresses
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dtos.DOCPaginatedResponse{data=[]dtos.DOCAddressResponse}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /addresses/user/{userId} [get]
// @Security BearerAuth
func (h *AddressHandler) ListByUserID(c *gin.Context) {
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

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	addresses, total, err := h.addressService.ListByUserID(c.Request.Context(), page, pageSize, uint(userID))
	if err != nil {
		c.Error(err)
		return
	}
	response.SuccessWithPagination(c, http.StatusOK, addresses, page, pageSize, total)
}

// List handles retrieving all addresses for the authenticated user
// @Summary List user addresses
// @Description Get a paginated list of addresses for the authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} dtos.DOCAddressListResponse
// @Failure 400 {object} dtos.DOCErrorResponse
// @Failure 401 {object} dtos.DOCErrorResponse
// @Failure 500 {object} dtos.DOCErrorResponse
// @Router /addresses [get]
// @Security BearerAuth
func (h *AddressHandler) List(c *gin.Context) {
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

	userID, ok := c.Get("user_id")
	if !ok {
		c.Error(errors.NewBadRequestError("invalid user ID in context"))
		return
	}

	addresses, total, err := h.addressService.ListByUserID(c.Request.Context(), page, pageSize, userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}
	response.SuccessWithPagination(c, http.StatusOK, addresses, page, pageSize, total)
}
