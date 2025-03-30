package handler

import (
	"ecom-go/internal/dtos"
	"net/http"
	"strconv"

	"ecom-go/internal/service"
	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, tokenService *service.TokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// Register sets up routes for the user handler
func (h *UserHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc ) {
	users := router.Group("/users")
	{
		users.POST("", h.Create)
		users.POST("/login", h.Login)
		users.POST("/refresh", h.RefreshToken)
		users.POST("/logout", h.Logout)
	}

	protected := users.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("", h.List)
		protected.GET("/:id", h.GetByID)
		protected.PUT("/:id", h.Update)
		protected.DELETE("/:id", h.Delete)
	}
}

// Login handles user login
// @Summary Login a user
// @Description Authenticates a user and returns access and refresh tokens
// @Tags User
// @Accept json
// @Produce json
// @Param   email     body    string     true  "Email"    binding:"required,email"
// @Param   password  body    string     true  "Password" binding:"required"
// @Param   use_cookies query  boolean    false "Set tokens as cookies"
// @Success 200 {object} service.TokenResponse "Token pair"
// @Failure 400 {object} errors.Error "Invalid input"
// @Failure 401 {object} errors.Error "Invalid credentials"
// @Failure 500 {object} errors.Error "Server error"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginDTO dtos.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}
		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	// Authenticate user
	user, err := h.userService.GetByEmail(c.Request.Context(), loginDTO.Email)
	if err != nil {
		c.Error(errors.NewUnauthorizedError("invalid credentials"))
		return
	}

	// Check password
	if err := user.ComparePassword(loginDTO.Password); err != nil {
		c.Error(errors.NewUnauthorizedError("invalid credentials"))
		return
	}

	// Generate token pair
	tokenPair, err := h.tokenService.GenerateTokenPair(c.Request.Context(), user)
	if err != nil {
		c.Error(err)
		return
	}

	// Check if using cookie-based auth
	useCookies := c.Query("use_cookies") == "true"
	if useCookies {
		// Set access token cookie
		c.SetCookie(
			"access_token",
			tokenPair.AccessToken,
			tokenPair.ExpiresIn,
			"/",
			"",
			c.Request.TLS != nil, // Secure (HTTPS only)
			true,                 // HTTP only
		)
		
		// Set refresh token cookie with longer expiration
		refreshExpiration := 60 * 60 * 24 * 7 // 7 days in seconds
		c.SetCookie(
			"refresh_token",
			tokenPair.RefreshToken,
			refreshExpiration,
			"/",
			"",
			c.Request.TLS != nil, // Secure
			true,                 // HTTP only
		)
	}

	response.Success(c, http.StatusOK, tokenPair)
}

// RefreshToken handles token refresh
// @Summary Refresh tokens
// @Description Exchanges a valid refresh token for a new token pair
// @Tags User
// @Accept json
// @Produce json
// @Param   refresh_token body string true "Refresh token"
// @Param   use_cookies query boolean false "Set tokens as cookies"
// @Success 200 {object} service.TokenResponse "New token pair"
// @Failure 400 {object} errors.Error "Invalid input"
// @Failure 401 {object} errors.Error "Invalid refresh token"
// @Failure 500 {object} errors.Error "Server error"
// @Router /users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	// Try to get refresh token from cookie first
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		req.RefreshToken = refreshToken
	} else {
		// If not in cookie, get from request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.NewBadRequestError("invalid input", err))
			return
		}
	}
	
	// Refresh tokens
	tokenPair, err := h.tokenService.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}
	
	// Check if using cookie-based auth
	useCookies := c.Query("use_cookies") == "true" || refreshToken != ""
	if useCookies {
		// Set access token cookie
		c.SetCookie(
			"access_token",
			tokenPair.AccessToken,
			tokenPair.ExpiresIn,
			"/",
			"",
			c.Request.TLS != nil, // Secure (HTTPS only)
			true,                 // HTTP only
		)
		
		// Set refresh token cookie with longer expiration
		refreshExpiration := 60 * 60 * 24 * 7 // 7 days in seconds
		c.SetCookie(
			"refresh_token",
			tokenPair.RefreshToken,
			refreshExpiration,
			"/",
			"",
			c.Request.TLS != nil, // Secure
			true,                 // HTTP only
		)
	}
	
	response.Success(c, http.StatusOK, tokenPair)
}

// Logout revokes the refresh token
// @Summary Logout
// @Description Invalidates the refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param   refresh_token body string false "Refresh token (not required when using cookies)"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} errors.Error "Invalid input"
// @Failure 401 {object} errors.Error "Invalid refresh token"
// @Failure 500 {object} errors.Error "Server error"
// @Router /users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	
	// Try to get refresh token from cookie first
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		req.RefreshToken = refreshToken
	} else {
		// If not in cookie, get from request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.NewBadRequestError("invalid input", err))
			return
		}
	}
	
	// If no refresh token provided, just clear cookies and return
	if req.RefreshToken == "" {
		if refreshToken != "" {
			// Clear cookies
			c.SetCookie("access_token", "", -1, "/", "", c.Request.TLS != nil, true)
			c.SetCookie("refresh_token", "", -1, "/", "", c.Request.TLS != nil, true)
		}
		response.Success(c, http.StatusOK, gin.H{"message": "Successfully logged out"})
		return
	}
	
	// Revoke refresh token
	err = h.tokenService.RevokeRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}
	
	// Clear cookies if they were used
	if refreshToken != "" {
		c.SetCookie("access_token", "", -1, "/", "", c.Request.TLS != nil, true)
		c.SetCookie("refresh_token", "", -1, "/", "", c.Request.TLS != nil, true)
	}
	
	response.Success(c, http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var createUserDTO dtos.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}
		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), createUserDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// GetByID handles retrieving a user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Update handles updating a user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	var updateUserDTO dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(errors.NewValidationError(&ve))
			return
		}

		c.Error(errors.NewBadRequestError("invalid input", err))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), updateUserDTO)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, user)
}

// Delete handles deleting a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.NewBadRequestError("invalid user ID"))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// List handles retrieving users with pagination
func (h *UserHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.Error(errors.NewBadRequestError("invalid page parameter"))
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil || pageSize <= 0 {
		c.Error(errors.NewBadRequestError("invalid per_page parameter"))
		return
	}

	users, total, err := h.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, users, page, pageSize, total)
}
