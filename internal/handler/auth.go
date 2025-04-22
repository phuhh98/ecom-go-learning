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
		auth.POST("/verify", h.VerifyToken)
	}
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param use_cookies query bool false "Set to 'true' to use cookie-based authentication"
// @Param credentials body dtos.LoginDTO true "User credentials"
// @Success 200 {object} dtos.DOCResponseWrapper{data=dtos.DOCTokenPair}
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 401 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /auth/login [post]
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

// RefreshToken handles refreshing authentication tokens
// @Summary Refresh token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param use_cookies query bool false "Set to 'true' to use cookie-based authentication"
// @Param refresh_request body dtos.DOCRefreshTokenRequest false "Refresh token (not required if using cookies)"
// @Success 200 {object} dtos.DOCTokenPair
// @Failure 400 {object} dtos.DOCErrorResponse
// @Failure 401 {object} dtos.DOCErrorResponse
// @Failure 500 {object} dtos.DOCErrorResponse
// @Router /auth/refresh [post]
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
			c.Request.Host,
			c.Request.TLS != nil, // Secure (HTTPS only)
			false,                // HTTP only
		)

		// Set refresh token cookie with longer expiration
		refreshExpiration := 60 * 60 * 24 * 7 // 7 days in seconds
		c.SetCookie(
			"refresh_token",
			tokenPair.RefreshToken,
			refreshExpiration,
			"/",
			"localhost",
			c.Request.TLS != nil, // Secure
			false,                // HTTP only
		)

	}

	response.Success(c, http.StatusOK, tokenPair)
}

// Logout invalidates refresh tokens
// @Summary Logout user
// @Description Invalidate a refresh token to log out
// @Tags auth
// @Accept json
// @Produce json
// @Param logout_request body dtos.DOCLogoutRequest false "Refresh token (not required if using cookies)"
// @Success 200 {object} dtos.DOCSuccessMessageResponse
// @Failure 400 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Failure 500 {object} dtos.DOCErrorWrapper{error=dtos.DOCStandardError}
// @Router /auth/logout [post]
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

// VerifyToken validates an access token
// @Summary Verify token
// @Description Validate an access token
// @Tags auth
// @Accept json
// @Produce json
// @Param verify_request body dtos.DOCVerifyTokenRequest false "Access token (not required if using cookies)"
// @Success 200 {object} dtos.DOCMessageResponse
// @Failure 400 {object} dtos.DOCErrorResponse
// @Failure 401 {object} dtos.DOCErrorResponse
// @Failure 500 {object} dtos.DOCErrorResponse
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req struct {
		AccessToken string `json:"access_token" binding:"required"`
	}

	// Try to get refresh token from cookie first
	accessToken, err := c.Cookie("access_token")
	if err == nil && accessToken != "" {
		req.AccessToken = accessToken
	} else {
		// If not in cookie, get from request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.NewBadRequestError("invalid input", err))
			return
		}
	}

	// Verify the access token
	err = h.tokenService.VerifyAccessToken(c.Request.Context(), req.AccessToken)
	if err != nil {
		c.Error(errors.NewUnauthorizedError("invalid or expired token"))
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Token is valid"})
}
