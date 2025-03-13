package service

import (
	"context"
	"ecom-go/internal/dtos"
	"errors"

	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"

	"github.com/davecgh/go-spew/spew"
)

// ItemService handles business logic related to items
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService creates a new item service
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

// // Create creates a new item
// func (s *ItemService) Create(ctx context.Context, createItemDTO dtos.CreateItemDTO) (*models.Item, error) {
// 	item := &models.Item{
// 		ProductID: createItemDTO.ProductId,
// 		Status:      createItemDTO.Status,
// 	}

// 	if err := s.repo.Create(ctx, item); err != nil {
// 		return nil, appError.NewServerError("error creating item", err)
// 	}

// 	return item, nil
// }

func (s *ItemService) CreateItems(ctx context.Context, createItemDtos dtos.CreateItemsDTO) ([]*models.Item, error) {

	var items []*models.Item
	for i := 0; i < createItemDtos.Quantity; i++ {
		item := &models.Item{
			ProductID: createItemDtos.ProductId,
			Status:    createItemDtos.Status,
		}
		items = append(items, item)
	}

	for _, item := range items {
		err := s.repo.Create(ctx, item)
		if err != nil {
			return nil, appError.NewServerError("error creating item", err)
		}
	}

	return items, nil

}

// GetByID retrieves an item by ID
func (s *ItemService) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("item not found")
		}
		return nil, appError.NewServerError("error retrieving item", err)
	}
	return item, nil
}

// Update updates an existing item
func (s *ItemService) Update(ctx context.Context, updateItemDTO dtos.UpdateItemDto) (*models.Item, error) {
	item, err := s.repo.GetByID(ctx, updateItemDTO.ItemId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("item not found")
		}
		return nil, appError.NewServerError("error retrieving item", err)
	}

	if updateItemDTO.Status != item.Status {
		item.Status = updateItemDTO.Status
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, appError.NewServerError("error updating item", err)
	}

	return item, nil
}

// Delete removes an item
func (s *ItemService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return appError.NewNotFoundError("item not found")
		}
		spew.Dump(err)
		return appError.NewServerError("error deleting item", err)
	}
	return nil
}

// List retrieves Items with pagination
func (s *ItemService) ListItemByProductId(ctx context.Context, page, pageSize int, productId uint) ([]*models.Item, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	// Get items
	addresss, err := s.repo.List(ctx, offset, pageSize, &productId)
	if err != nil {
		return nil, 0, appError.NewServerError("error retrieving items", err)
	}

	// Get total count
	total, err := s.repo.Count(ctx, &productId)
	if err != nil {
		return nil, 0, appError.NewServerError("error counting items", err)
	}

	return addresss, total, nil
}
