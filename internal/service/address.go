package service

import (
	"context"
	"ecom-go/internal/dtos"
	"errors"

	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

// AddressService handles business logic related to addresss
type AddressService struct {
	repo repository.AddressRepository
}

// NewAddressService creates a new address service
func NewAddressService(repo repository.AddressRepository) *AddressService {
	return &AddressService{
		repo: repo,
	}
}

// Create creates a new address
func (s *AddressService) Create(ctx context.Context, createAddressDTO dtos.CreateAddressDTO) (*models.Address, error) {
	// Create new address
	address := &models.Address{
		UserID:      createAddressDTO.UserID,
		FullName:    createAddressDTO.FullName,
		PhoneNumber: createAddressDTO.PhoneNumber,
		AddressLine: createAddressDTO.AddressLine,
	}

	// Save to database
	if err := s.repo.Create(ctx, address); err != nil {
		return nil, appError.NewServerError("error creating address", err)
	}

	return address, nil
}

// GetByID retrieves a address by ID
func (s *AddressService) GetByID(ctx context.Context, id uint) (*models.Address, error) {
	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("address not found")
		}
		return nil, appError.NewServerError("error retrieving address", err)
	}
	return address, nil
}

// Update updates an existing address
func (s *AddressService) Update(ctx context.Context, id uint, updateAddressDTO dtos.UpdateAddressDTO) (*models.Address, error) {
	// Get existing address
	address, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("address not found")
		}
		return nil, appError.NewServerError("error retrieving address", err)
	}

	// Update fields
	if updateAddressDTO.FullName != "" {
		address.FullName = updateAddressDTO.FullName
	}
	if updateAddressDTO.PhoneNumber != address.PhoneNumber {
		address.PhoneNumber = updateAddressDTO.PhoneNumber
	}
	if updateAddressDTO.AddressLine != "" {
		address.AddressLine = updateAddressDTO.AddressLine
	}

	// Save to database
	if err := s.repo.Update(ctx, address); err != nil {
		return nil, appError.NewServerError("error updating address", err)
	}

	return address, nil
}

// Delete removes a address
func (s *AddressService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return appError.NewNotFoundError("address not found")
		}
		return appError.NewServerError("error deleting address", err)
	}
	return nil
}

// // List retrieves addresss with pagination
// func (s *AddressService) List(ctx context.Context, page, pageSize int) ([]*models.Address, int64, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 {
// 		pageSize = 10
// 	}

// 	offset := (page - 1) * pageSize

// 	// Get addresss
// 	addresss, err := s.repo.List(ctx, offset, pageSize)
// 	if err != nil {
// 		return nil, 0, appError.NewServerError("error retrieving addresss", err)
// 	}

// 	// Get total count
// 	total, err := s.repo.Count(ctx)
// 	if err != nil {
// 		return nil, 0, appError.NewServerError("error counting addresss", err)
// 	}

// 	return addresss, total, nil
// }
