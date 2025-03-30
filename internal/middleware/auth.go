package middleware

import (
	"ecom-go/internal/config"
	"ecom-go/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware holds the configuration for authentication
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

// getToken extracts token from Authorization header or cookie
func (m *AuthMiddleware) getToken(c *gin.Context) (string, error) {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}
	
	// Try cookie
	token, err := c.Cookie("access_token")
	if err == nil && token != "" {
		return token, nil
	}
	
	return "", errors.NewUnauthorizedError("no valid authentication token found")
}

// Authenticate verifies the JWT token in Authorization header or cookie
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header or cookie
		tokenString, err := m.getToken(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		
		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewUnauthorizedError("invalid token signing method")
			}
			return []byte(m.config.GetJWTSecret()), nil
		})
		
		// Handle token validation errors
		if err != nil {
			// Check if the error is due to an expired token
			if strings.Contains(err.Error(), "token is expired") {
				c.Error(errors.NewUnauthorizedError("token has expired"))
			} else {
				c.Error(errors.NewUnauthorizedError("invalid token: " + err.Error()))
			}
			c.Abort()
			return
		}
		
		if !token.Valid {
			c.Error(errors.NewUnauthorizedError("invalid token"))
			c.Abort()
			return
		}
		
		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(errors.NewUnauthorizedError("invalid token claims"))
			c.Abort()
			return
		}
		
		// Set user info in context for handlers
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.Error(errors.NewUnauthorizedError("invalid token: missing user_id"))
			c.Abort()
			return
		}
		
		c.Set("userID", uint(userID))
		
		// Set email if present in token
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		
		// Set role if present in token
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}
		
		c.Next()
	}
}

// RequireRole middleware to check if user has required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by Authenticate middleware)
		role, exists := c.Get("role")
		if !exists {
			c.Error(errors.NewUnauthorizedError("role not found in context"))
			c.Abort()
			return
		}
		
		userRole, ok := role.(string)
		if !ok {
			c.Error(errors.NewUnauthorizedError("invalid role type"))
			c.Abort()
			return
		}
		
		// Check if user has any of the required roles
		hasRole := false
		for _, r := range roles {
			if r == userRole {
				hasRole = true
				break
			}
		}
		
		if !hasRole {
			c.Error(errors.NewForbiddenError("insufficient permissions"))
			c.Abort()
			return
		}
		
		c.Next()
	}
}
