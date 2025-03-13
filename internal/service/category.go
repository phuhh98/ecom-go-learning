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
	repo           repository.CategoryRepository
	productService *ProductService
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository, productService *ProductService) *CategoryService {
	return &CategoryService{
		repo:           repo,
		productService: productService,
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

func (s *CategoryService) Update(ctx context.Context, categoryID uint, updateCategoryDTO dtos.UpdateCategoryDTO) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("category not found")
		}
		return nil, appError.NewServerError("error retrieving category", err)
	}

	// Check if category name is being changed and is already in use
	if updateCategoryDTO.Name != "" && updateCategoryDTO.Name != category.Name {
		existingCategory, err := s.repo.GetByName(ctx, updateCategoryDTO.Name)
		if err == nil && existingCategory.ID != categoryID {
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

	updatedCategory, err := s.GetByID(ctx, categoryID)
	return updatedCategory, nil
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

func (s *CategoryService) AddProductToCategory(ctx context.Context, categoryID uint, productID uint) (*models.Category, error) {
	category, err := s.GetByID(ctx, categoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("category not found")
	}

	product, err := s.productService.GetByID(ctx, productID)
	if err != nil {
		return nil, appError.NewNotFoundError("product not found")
	}

	// Check if product already exists in category
	for _, p := range category.Products {
		if p.ID == productID {
			return nil, appError.NewBadRequestError("product already exists in category")
		}
	}

	category.Products = append(category.Products, product)
	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, appError.NewServerError("error updating category", err)
	}
	updatedCategory, err := s.GetByID(ctx, categoryID)
	return updatedCategory, nil
}

func (s *CategoryService) RemoveProductFromCategory(ctx context.Context, categoryID uint, productID uint) (*models.Category, error) {
	category, err := s.GetByID(ctx, categoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("category not found")
	}
	_, err = s.productService.GetByID(ctx, productID)
	if err != nil {
		return nil, appError.NewNotFoundError("product not found")
	}
	// Check if product exists in category
	var index int
	for i, p := range category.Products {
		if p.ID == productID {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, appError.NewBadRequestError("product does not exist in category")
	}

	// Remove product from category
	category.Products = append(category.Products[:index], category.Products[index+1:]...)
	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, appError.NewServerError("error updating category", err)
	}
	updatedCategory, err := s.GetByID(ctx, categoryID)
	return updatedCategory, nil
}

func (s *CategoryService) AddSubcategoryToCategory(ctx context.Context, categoryID uint, subcategoryID uint) (*models.Category, error) {
	category, err := s.GetByID(ctx, categoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("category not found")
	}

	subcategory, err := s.GetByID(ctx, subcategoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("subcategory not found")
	}

	// Check if sucategory already exists in category
	for _, sc := range category.SubCategories {
		if sc.ID == subcategoryID {
			return nil, appError.NewBadRequestError("subcategory already exists in category")
		}
	}

	// Perform an update with sub category attached to category
	category.SubCategories = append(category.SubCategories, subcategory)
	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, appError.NewServerError("error updating category", err)
	}
	updatedCategory, err := s.GetByID(ctx, categoryID)
	return updatedCategory, nil
}

func (s *CategoryService) RemoveSubcategoryFromCategory(ctx context.Context, categoryID uint, subcategoryID uint) (*models.Category, error) {
	category, err := s.GetByID(ctx, categoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("category not found")
	}

	_, err = s.GetByID(ctx, subcategoryID)
	if err != nil {
		return nil, appError.NewNotFoundError("subcategory not found")
	}

	// Perform an update with sub category removed from category
	for i, subCategory := range category.SubCategories {
		if subCategory.ID == subcategoryID {
			category.SubCategories = append(category.SubCategories[:i], category.SubCategories[i+1:]...)
			break
		}
	}

	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, appError.NewServerError("error updating category", err)
	}
	updatedCategory, err := s.GetByID(ctx, categoryID)
	return updatedCategory, nil
}
