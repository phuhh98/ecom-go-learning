package repository

import (
	"fmt"

	"ecom-go/internal/config"

	"gorm.io/gorm"
)

// Factory provides access to all repositories
type Factory struct {
	db   *gorm.DB
	User UserRepository
	// Add other repositories here as you implement them
	// Product ProductRepository
	// Order   OrderRepository
	Product ProductRepository
	Address AddressRepository
}

// NewFactory creates a new repository factory
func NewFactory(cfg *config.Config) (*Factory, error) {
	// Initialize database connection
	db, err := NewDatabase(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Auto-migrate database schema if needed
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	// Create repository instances
	return &Factory{
		db:   db,
		User: NewUserRepo(db),
		// Initialize other repositories here as you implement them
		Product: NewProductRepo(db),
		Address: NewAddressRepo(db),
	}, nil
}

// Close closes the database connection
func (f *Factory) Close() error {
	sqlDB, err := f.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
