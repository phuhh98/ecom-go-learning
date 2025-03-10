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
