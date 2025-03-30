package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig holds all the server-related configuration
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig holds all the database-related configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

// RedisConfig holds all the Redis-related configuration
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// RabbitMQConfig holds all the RabbitMQ-related configuration
type RabbitMQConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// JWTConfig holds all JWT-related configuration
type JWTConfig struct {
	Secret                string `mapstructure:"secret"`
	AccessExpirationMins  int    `mapstructure:"access_expiration_mins"`
	RefreshExpirationDays int    `mapstructure:"refresh_expiration_days"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	// Set default configuration paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Explicitly bind environment variables
	viper.BindEnv("database.host", "APP_DATABASE_HOST")
	viper.BindEnv("database.port", "APP_DATABASE_PORT")
	viper.BindEnv("database.name", "APP_DATABASE_NAME")
	viper.BindEnv("database.user", "APP_DATABASE_USER")
	viper.BindEnv("database.password", "APP_DATABASE_PASSWORD")
	viper.BindEnv("database.sslmode", "APP_DATABASE_SSLMODE")

	// Bind server, Redis, and RabbitMQ configs similarly
	viper.BindEnv("server.port", "APP_SERVER_PORT")
	viper.BindEnv("redis.host", "APP_REDIS_HOST")
	viper.BindEnv("redis.port", "APP_REDIS_PORT")
	viper.BindEnv("rabbitmq.host", "APP_RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "APP_RABBITMQ_PORT")
	
	// Bind JWT config
	viper.BindEnv("jwt.secret", "APP_JWT_SECRET")
	viper.BindEnv("jwt.access_expiration_mins", "APP_JWT_ACCESS_EXPIRATION_MINS")
	viper.BindEnv("jwt.refresh_expiration_days", "APP_JWT_REFRESH_EXPIRATION_DAYS")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			fmt.Printf("Warning: Config file not found, using environment variables only\n")
		} else {
			// Config file was found but another error occurred
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode)
}

// GetJWTSecret returns the JWT secret with a fallback to a default
func (c *Config) GetJWTSecret() string {
	if c.JWT.Secret == "" {
		return "default_secret_change_in_production" // Fallback default
	}
	return c.JWT.Secret
}

// GetJWTAccessExpirationMinutes returns the JWT access token expiration time in minutes
func (c *Config) GetJWTAccessExpirationMinutes() int {
	if c.JWT.AccessExpirationMins <= 0 {
		return 30 // Default to 30 minutes
	}
	return c.JWT.AccessExpirationMins
}

// GetJWTRefreshExpirationDays returns the JWT refresh token expiration time in days
func (c *Config) GetJWTRefreshExpirationDays() int {
	if c.JWT.RefreshExpirationDays <= 0 {
		return 7 // Default to 7 days
	}
	return c.JWT.RefreshExpirationDays
}
