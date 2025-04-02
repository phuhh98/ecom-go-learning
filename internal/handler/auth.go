package handler

import (
	"ecom-go/internal/dtos"
	"net/http"

	"ecom-go/internal/service"
	"ecom-go/pkg/errors"
	"ecom-go/pkg/http/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles HTTP requests related to users
type AuthHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewUserHandler creates a new user handler
func NewAuthHandler(userService *service.UserService, tokenService *service.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// Register sets up routes for the user handler
func (h *AuthHandler) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
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
func (h *AuthHandler) Login(c *gin.Context) {
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
func (h *AuthHandler) RefreshToken(c *gin.Context) {
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
func (h *AuthHandler) Logout(c *gin.Context) {
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
