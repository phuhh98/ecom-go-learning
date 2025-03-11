package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepo implements the UserRepository interface using PostgreSQL/GORM
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new user repository
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// Create adds a new user to the database
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Preload(clause.Associations).Create(user)
	if result.Error != nil {
		// Check for unique constraint violation
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return result.Error
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Preload(clause.Associations).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).Preload(clause.Associations).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Preload(clause.Associations).Save(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return result.Error
	}
	return nil
}

// Delete removes a user from the database
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List retrieves users with pagination
func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Preload(clause.Associations).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Count returns the total number of users
func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
