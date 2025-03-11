package service

import (
	"context"
	"ecom-go/internal/dtos"
	"errors"

	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

// ProductService handles business logic related to Product
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// Create creates a new product
func (s *ProductService) Create(ctx context.Context, createProductDTO dtos.CreateProductDTO) (*models.Product, error) {
	// Check if user with same email already exists
	_, err := s.repo.GetByCode(ctx, createProductDTO.Code)
	if err == nil {
		return nil, appError.NewBadRequestError("product code already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, appError.NewServerError("error checking existing product", err)
	}

	product := &models.Product{
		Code:        createProductDTO.Code,
		Name:        createProductDTO.Name,
		Description: createProductDTO.Description,
		Price:       createProductDTO.Price,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, appError.NewServerError("error creating product", err)
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("product not found")
		}
		return nil, appError.NewServerError("error retrieving product", err)
	}
	return product, nil
}

func (s *ProductService) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	product, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("product not found")
		}
		return nil, appError.NewServerError("error retrieving user", err)
	}
	return product, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, updateProductDTO dtos.UpdateProductDTO) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("product not found")
		}
		return nil, appError.NewServerError("error retrieving product", err)
	}

	// Check if product code is being changed and is already in use
	if updateProductDTO.Code != "" && updateProductDTO.Code != product.Code {
		existingProduct, err := s.repo.GetByCode(ctx, updateProductDTO.Code)
		if err == nil && existingProduct.ID != id {
			return nil, appError.NewBadRequestError("product code already exists")
		} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewServerError("error checking existing product code", err)
		}
		product.Code = updateProductDTO.Code
	}

	// Update fields
	if updateProductDTO.Name != "" {
		product.Name = updateProductDTO.Name
	}
	if updateProductDTO.Description != "" {
		product.Description = updateProductDTO.Description
	}
	if updateProductDTO.Price != product.Price {
		product.Price = updateProductDTO.Price
	}

	// Save to database
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, appError.NewServerError("error updating product", err)
	}

	return product, nil
}

// Delete removes a product
func (s *ProductService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return appError.NewNotFoundError("product not found")
		}
		return appError.NewServerError("error deleting product", err)
	}
	return nil
}

// List retrieves products with pagination
func (s *ProductService) List(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get products
	products, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, appError.NewServerError("error retrieving products", err)
	}

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, appError.NewServerError("error counting products", err)
	}

	return products, total, nil
}
