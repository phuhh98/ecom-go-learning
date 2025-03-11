package service

import (
	"context"
	"ecom-go/internal/dtos"
	"errors"

	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

// UserService handles business logic related to users
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, createUserDTO dtos.CreateUserDTO) (*models.User, error) {
	// Check if user with same email already exists
	_, err := s.repo.GetByEmail(ctx, createUserDTO.Email)
	if err == nil {
		return nil, appError.NewBadRequestError("email already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, appError.NewServerError("error checking existing user", err)
	}

	// Create new user
	user := &models.User{
		Email:     createUserDTO.Email,
		Password:  createUserDTO.Password,
		FirstName: createUserDTO.FirstName,
		LastName:  createUserDTO.LastName,
		// Role:      "user", // Default role
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, appError.NewServerError("error hashing password", err)
	}

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, appError.NewServerError("error creating user", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("user not found")
		}
		return nil, appError.NewServerError("error retrieving user", err)
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("user not found")
		}
		return nil, appError.NewServerError("error retrieving user", err)
	}
	return user, nil
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, id uint, updateUserDTO dtos.UpdateUserDTO) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewNotFoundError("user not found")
		}
		return nil, appError.NewServerError("error retrieving user", err)
	}

	// Check if email is being changed and is already in use
	if updateUserDTO.Email != "" && updateUserDTO.Email != user.Email {
		existingUser, err := s.repo.GetByEmail(ctx, updateUserDTO.Email)
		if err == nil && existingUser.ID != id {
			return nil, appError.NewBadRequestError("email already in use")
		} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, appError.NewServerError("error checking existing email", err)
		}
		user.Email = updateUserDTO.Email
	}

	// Update fields
	if updateUserDTO.FirstName != "" {
		user.FirstName = updateUserDTO.FirstName
	}
	if updateUserDTO.LastName != "" {
		user.LastName = updateUserDTO.LastName
	}
	// user.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, appError.NewServerError("error updating user", err)
	}

	return user, nil
}

// Delete removes a user
func (s *UserService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return appError.NewNotFoundError("user not found")
		}
		return appError.NewServerError("error deleting user", err)
	}
	return nil
}

// List retrieves users with pagination
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get users
	users, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, appError.NewServerError("error retrieving users", err)
	}

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, appError.NewServerError("error counting users", err)
	}

	return users, total, nil
}
