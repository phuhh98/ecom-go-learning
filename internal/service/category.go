package service

import (
	"context"
	"ecom-go/internal/dtos"
	"errors"

	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

// CategoryService handles business logic related to Category
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, createCategoryDTO dtos.CreateCategoryDTO) (*models.Category, error) {
	// Check if category with the same already exists
	_, err := s.repo.GetByName(ctx, createCategoryDTO.Name)
	if err == nil {
		return nil, appError.NewBadRequestError("category name already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, appError.NewServerError("error checking existing category", err)
	}

	category := &models.Category{
		Name:        createCategoryDTO.Name,
		Description: createCategoryDTO.Description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, appError.NewServerError("error creating category", err)
	}

	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("category not found")
		}
		return nil, appError.NewServerError("error retrieving category", err)
	}
	return category, nil
}

func (s *CategoryService) GetByName(ctx context.Context, name string) (*models.Category, error) {
	category, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("category not found")
		}
		return nil, appError.NewServerError("error retrieving user", err)
	}
	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id uint, updateCategoryDTO dtos.UpdateCategoryDTO) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("category not found")
		}
		return nil, appError.NewServerError("error retrieving category", err)
	}

	// Check if category name is being changed and is already in use
	if updateCategoryDTO.Name != "" && updateCategoryDTO.Name != category.Name {
		existingCategory, err := s.repo.GetByName(ctx, updateCategoryDTO.Name)
		if err == nil && existingCategory.ID != id {
			return nil, appError.NewBadRequestError("category code already exists")
		} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewServerError("error checking existing category code", err)
		}
		category.Name = updateCategoryDTO.Name
	}

	// Update fields
	if updateCategoryDTO.Name != "" {
		category.Name = updateCategoryDTO.Name
	}
	if updateCategoryDTO.Description != "" {
		category.Description = updateCategoryDTO.Description
	}

	// Save to database
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, appError.NewServerError("error updating category", err)
	}

	return category, nil
}

// Delete removes a category
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return appError.NewNotFoundError("category not found")
		}
		return appError.NewServerError("error deleting category", err)
	}
	return nil
}

// List retrieves categorys with pagination
func (s *CategoryService) List(ctx context.Context, page, pageSize int) ([]*models.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get categorys
	categorys, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, appError.NewServerError("error retrieving categorys", err)
	}

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, appError.NewServerError("error counting categorys", err)
	}

	return categorys, total, nil
}
