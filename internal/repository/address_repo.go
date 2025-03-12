package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
)

// AddressRepo implements the AddressRepository interface using PostgreSQL/GORM
type AddressRepo struct {
	db *gorm.DB
}

// NewAddressRepo creates a new address repository
func NewAddressRepo(db *gorm.DB) *AddressRepo {
	return &AddressRepo{
		db: db,
	}
}

// TODO: Template - need to update
// Create adds a new address to the database
func (r *AddressRepo) Create(ctx context.Context, address *models.Address) error {
	result := r.db.WithContext(ctx).Create(address)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves a address by ID
func (r *AddressRepo) GetByID(ctx context.Context, id uint) (*models.Address, error) {
	var address models.Address
	result := r.db.WithContext(ctx).First(&address, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &address, nil
}

// Update updates an existing address
func (r *AddressRepo) Update(ctx context.Context, address *models.Address) error {
	result := r.db.WithContext(ctx).Save(address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	return nil
}

// Delete removes a address from the database
func (r *AddressRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Address{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// // List retrieves addresss with pagination
// func (r *AddressRepo) List(ctx context.Context, offset, limit int) ([]*models.Address, error) {
// 	var addresss []*models.Address
// 	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&addresss)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return addresss, nil
// }

// // Count returns the total number of addresss
// func (r *AddressRepo) Count(ctx context.Context) (int64, error) {
// 	var count int64
// 	result := r.db.WithContext(ctx).Model(&models.Address{}).Count(&count)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}
// 	return count, nil
// }
