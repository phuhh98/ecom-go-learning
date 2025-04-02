package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"ecom-go/internal/config"
	"ecom-go/internal/models"
	"ecom-go/pkg/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDatabase establishes a connection to the PostgreSQL database
// or creates it if it doesn't exist
func NewDatabase(config *config.DatabaseConfig) (*gorm.DB, error) {
	// First, ensure the database exists
	if err := createDatabaseIfNotExists(config); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Configure GORM logger
	gormLogger := gormLogger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to the database

	fmt.Println("--newDB: connecting")
	dsn := config.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("--newDB: connected")

	// Get the underlying SQL DB object
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping database to verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	logger.Info("Connected to database", "database", config.Name)
	return db, nil
}

// createDatabaseIfNotExists connects to PostgreSQL and creates
// the database if it doesn't already exist
func createDatabaseIfNotExists(config *config.DatabaseConfig) error {
	fmt.Printf("Host: %s\n", config.Host)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("User: %s\n", config.User)
	fmt.Printf("Password: %s\n", config.Password) // For security
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("SSLMode: %s\n", config.SSLMode)
	// Connect to the 'postgres' database to create our app database
	postgresConnStr := config.GetDSN()

	// Connect to postgres database

	fmt.Println("--createDB: connecting to DB")
	db, err := sql.Open("postgres", postgresConnStr)
	if err != nil {
		return fmt.Errorf("error connecting to postgres database: %w", err)
	}
	fmt.Println("--createDB: connected")
	defer db.Close()

	// Check if the database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", config.Name)

	fmt.Println("--createDB: checking")
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		return fmt.Errorf("error checking if database exists: %w", err)
	}
	fmt.Println("--createDB: checked: " + strconv.FormatBool(exists))

	// Create the database if it doesn't exist
	if !exists {
		logger.Info("Creating database", "database", config.Name)

		// Create the database - need to use fmt.Sprintf as the database name is an identifier
		// and shouldn't be passed as a parameter
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Name))
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
		logger.Info("Database created successfully", "database", config.Name)
	} else {
		logger.Info("Database already exists", "database", config.Name)
	}
	fmt.Println("--createDB: done")

	return nil
}

// AutoMigrate automatically migrates database schema
func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running auto migration")

	// Add models for migration here
	err := db.AutoMigrate(
		&models.User{},
		// Add other models for auto-migration here as they are created
		&models.Product{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
		&models.Address{},
	)

	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	logger.Info("Auto migration completed successfully")
	return nil
}
