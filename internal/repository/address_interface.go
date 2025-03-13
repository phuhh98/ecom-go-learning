package repository

import (
	"context"

	"ecom-go/internal/models"
)

type AddressRepository interface {
	// Add new Address to db
	Create(ctx context.Context, user *models.Address) error

	// Get Address by Id
	GetByID(ctx context.Context, id uint) (*models.Address, error)

	// Update updates an existing Address
	Update(ctx context.Context, user *models.Address) error

	// Remove a Address from the database
	Delete(ctx context.Context, id uint) error

	// List retrieves Addresss with pagination
	List(ctx context.Context, offset, limit int, userID *uint) ([]*models.Address, error)

	// Count returns the total number of users
	Count(ctx context.Context, userID *uint) (int64, error)
}
