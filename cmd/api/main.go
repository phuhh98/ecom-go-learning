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
	"ecom-go/pkg/validators"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	tokenService := service.NewTokenService(cfg, repoFactory.Redis, repoFactory.User)
	userService := service.NewUserService(repoFactory.User, cfg)
	// TODO: Add other services here
	productService := service.NewProductService(repoFactory.Product)
	addressService := service.NewAddressService(repoFactory.Address)
	categoryService := service.NewCategoryService(repoFactory.Category, productService)
	orderService := service.NewOrderService(repoFactory.Order, productService, userService)

	// Set up HTTP server with Gin
	router := setupRouter(cfg)

	// Register handlers
	api := router.Group("/api/v1")

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, tokenService)
	authHandler.Register(api)

	userHandler := handler.NewUserHandler(userService, tokenService)
	userHandler.Register(api, authMiddleware.Authenticate())

	productHandler := handler.NewProductHandler(productService)
	productHandler.Register(api, authMiddleware.Authenticate(), authMiddleware.RequireRole("admin"))

	addressHandler := handler.NewAddressHandler(addressService)
	addressHandler.Register(api, authMiddleware.Authenticate())

	orderHandler := handler.NewOrderHandler(orderService)
	orderHandler.Register(api, authMiddleware.Authenticate())

	categoryHandler := handler.NewCategoryHandler(categoryService)
	categoryHandler.Register(api, authMiddleware.Authenticate(), authMiddleware.RequireRole("admin"))

	// // Admin-only routes
	// admin := api.Group("/admin")
	// admin.Use(authMiddleware.Authenticate(), authMiddleware.RequireRole("admin"))
	// {
	// 	// Admin-only endpoints can go here
	// }

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

func setupRouter(appConfig *config.Config) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("default", validators.Default)
	}

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Error())

	config := cors.DefaultConfig()

	// Get Cors allow origins from environment variable
	config.AllowOrigins = appConfig.Cors.AllowOrigins
	config.AllowCredentials = true

	router.Use(cors.New(config))

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return router
}
