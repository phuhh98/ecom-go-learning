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
		addresses.GET("/:id", h.GetByID)
		addresses.PUT("/:id", h.Update)
		addresses.DELETE("/:id", h.Delete)

		addresses.GET("/user/:userId", h.ListByUserID)
	}
}

// Create handles address creation
func (h *AddressHandler) Create(c *gin.Context) {
	var createAddressDTO dtos.CreateAddressDTO
	if err := c.ShouldBindJSON(&createAddressDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	address, err := h.addressService.Create(c.Request.Context(), createAddressDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, address)
}

// GetByID handles retrieving a address by ID
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

	address, err := h.addressService.Update(c.Request.Context(), uint(id), updateAddressDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// Delete handles deleting a address
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid address ID"))
		return
	}

	if err := h.addressService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

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
