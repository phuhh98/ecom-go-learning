package repository

import (
	"fmt"

	"ecom-go/internal/config"

	"gorm.io/gorm"
)

// Factory provides access to all repositories
type Factory struct {
	db    *gorm.DB
	Redis *RedisClient
	User  UserRepository
	// Add other repositories here as you implement them
	Product  ProductRepository
	Address  AddressRepository
	Category CategoryRepository
	Order    OrderRepository
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

	// Initialize Redis client
	redisClient, err := NewRedisClient(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	// Create repository instances
	return &Factory{
		db:    db,
		Redis: redisClient,
		User:  NewUserRepo(db),
		// Initialize other repositories here as you implement them
		Product:  NewProductRepo(db),
		Address:  NewAddressRepo(db),
		Category: NewCategoryRepo(db),
		Order:    NewOrderRepo(db),
	}, nil
}

// Close closes the database and Redis connections
func (f *Factory) Close() error {
	// Close database connection
	sqlDB, err := f.db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	// Close Redis connection
	if f.Redis != nil {
		if err := f.Redis.Close(); err != nil {
			return err
		}
	}

	return nil
}
