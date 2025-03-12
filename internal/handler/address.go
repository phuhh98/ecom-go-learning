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
func (h *AddressHandler) Register(router *gin.RouterGroup) {
	addresss := router.Group("/addresses")
	{
		addresss.POST("", h.Create)
		// addresss.GET("", h.List)
		addresss.GET("/:id", h.GetByID)
		addresss.PUT("/:id", h.Update)
		addresss.DELETE("/:id", h.Delete)
	}
}

// Create handles address creation
func (h *AddressHandler) Create(c *gin.Context) {
	var createAddressDTO dtos.CreateAddressDTO
	if err := c.ShouldBindJSON(&createAddressDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.Error(c, errors.NewValidationError(&ve))
			return
		}

		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	address, err := h.addressService.Create(c.Request.Context(), createAddressDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, address)
}

// GetByID handles retrieving a address by ID
func (h *AddressHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid address ID"))
		return
	}

	address, err := h.addressService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// Update handles updating a address
func (h *AddressHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid address ID"))
		return
	}

	var updateAddressDTO dtos.UpdateAddressDTO
	if err := c.ShouldBindJSON(&updateAddressDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.Error(c, errors.NewValidationError(&ve))
			return
		}

		response.Error(c, errors.NewBadRequestError("invalid input", err))
		return
	}

	address, err := h.addressService.Update(c.Request.Context(), uint(id), updateAddressDTO)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// Delete handles deleting a address
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid address ID"))
		return
	}

	if err := h.addressService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// // List handles retrieving addresss with pagination
// func (h *AddressHandler) List(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

// 	addresss, total, err := h.addressService.List(c.Request.Context(), page, pageSize)
// 	if err != nil {
// 		response.Error(c, err)
// 		return
// 	}

// 	response.SuccessWithPagination(c, http.StatusOK, addresss, page, pageSize, total)
// }
