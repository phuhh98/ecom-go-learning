package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecom-go/internal/config"
	"ecom-go/internal/handler"
	"ecom-go/internal/middleware"
	"ecom-go/internal/repository"
	"ecom-go/internal/service"
	"ecom-go/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", "error", err)
	}

	// Set up repository
	repoFactory, err := repository.NewFactory(cfg)
	if err != nil {
		logger.Fatal("Failed to create repository factory", "error", err)
	}
	defer func() {
		if err := repoFactory.Close(); err != nil {
			logger.Error("Error closing database connection", "error", err)
		}
	}()

	// Set up services
	userService := service.NewUserService(repoFactory.User)
	// TODO: Add other services here
	itemService := service.NewItemService(repoFactory.Item)
	productService := service.NewProductService(repoFactory.Product, itemService)
	addressService := service.NewAddressService(repoFactory.Address)
	categoryService := service.NewCategoryService(repoFactory.Category, productService)

	// Set up HTTP server with Gin
	router := setupRouter()

	// Register handlers
	api := router.Group("/api/v1")
	userHandler := handler.NewUserHandler(userService)
	userHandler.Register(api)
	// TODO: Add other handlers here
	productHandler := handler.NewProductHandler(productService)
	productHandler.Register(api)

	addressHandler := handler.NewAddressHandler(addressService)
	addressHandler.Register(api)

	categoryHandler := handler.NewCategoryHandler(categoryService)
	categoryHandler.Register(api)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited properly")
}

func setupRouter() *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Error())

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return router
}
