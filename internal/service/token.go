package service

import (
	"context"
	"ecom-go/internal/config"
	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenResponse represents the response for token operations
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Access token expiration in seconds
}

// TokenService handles token operations
type TokenService struct {
	config   *config.Config
	redis    *repository.RedisClient
	userRepo repository.UserRepository
}

// NewTokenService creates a new token service
func NewTokenService(config *config.Config, redis *repository.RedisClient, userRepo repository.UserRepository) *TokenService {
	return &TokenService{
		config:   config,
		redis:    redis,
		userRepo: userRepo,
	}
}

// GenerateTokenPair generates both access and refresh tokens for a user
func (s *TokenService) GenerateTokenPair(ctx context.Context, user *models.User) (*TokenResponse, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, appError.NewServerError("failed to generate access token", err)
	}
	
	// Generate refresh token
	refreshToken, tokenID, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, appError.NewServerError("failed to generate refresh token", err)
	}
	
	// Calculate expiration time
	expirationDays := s.config.GetJWTRefreshExpirationDays()
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(expirationDays))
	
	// Store in Redis
	if err := s.storeRefreshToken(ctx, user.ID, tokenID, expiresAt); err != nil {
		return nil, appError.NewServerError("failed to store refresh token", err)
	}
	
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.GetJWTAccessExpirationMinutes() * 60,
	}, nil
}

// generateAccessToken creates a short-lived access token
func (s *TokenService) generateAccessToken(user *models.User) (string, error) {
	return user.GenerateJWT(
		s.config.GetJWTSecret(),
		s.config.GetJWTAccessExpirationMinutes(),
	)
}

// generateRefreshToken creates a long-lived refresh token
func (s *TokenService) generateRefreshToken(user *models.User) (string, string, error) {
	// Generate unique token ID
	tokenID := uuid.New().String()
	
	// Generate token
	token, err := user.GenerateRefreshToken(
		s.config.GetJWTSecret(),
		tokenID,
		s.config.GetJWTRefreshExpirationDays(),
	)
	
	if err != nil {
		return "", "", err
	}
	
	return token, tokenID, nil
}

// storeRefreshToken stores a refresh token in Redis
func (s *TokenService) storeRefreshToken(ctx context.Context, userID uint, tokenID string, expiresAt time.Time) error {
	// Calculate TTL
	ttl := time.Until(expiresAt)
	
	// Create Redis key
	key := fmt.Sprintf("refresh_token:%s", tokenID)
	
	// Store in Redis with expiration
	return s.redis.Set(ctx, key, userID, ttl)
}

// ValidateRefreshToken validates a refresh token and returns the associated user
func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*models.User, string, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appError.NewUnauthorizedError("invalid token signing method")
		}
		return []byte(s.config.GetJWTSecret()), nil
	})
	
	if err != nil || !token.Valid {
		return nil, "", appError.NewUnauthorizedError("invalid refresh token")
	}
	
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", appError.NewUnauthorizedError("invalid token claims")
	}
	
	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, "", appError.NewUnauthorizedError("invalid token type")
	}
	
	// Extract user ID and token ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, "", appError.NewUnauthorizedError("invalid user ID in token")
	}
	userID := uint(userIDFloat)
	
	tokenID, ok := claims["token_id"].(string)
	if !ok {
		return nil, "", appError.NewUnauthorizedError("invalid token ID")
	}
	
	// Check if token exists in Redis
	key := fmt.Sprintf("refresh_token:%s", tokenID)
	storedValue, err := s.redis.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return nil, "", appError.NewUnauthorizedError("refresh token not found or expired")
		}
		return nil, "", appError.NewServerError("failed to validate refresh token", err)
	}
	
	// Parse stored user ID
	storedUserID, err := strconv.ParseUint(storedValue, 10, 64)
	if err != nil {
		return nil, "", appError.NewServerError("invalid stored user ID", err)
	}
	
	if uint(storedUserID) != userID {
		return nil, "", appError.NewUnauthorizedError("token mismatch")
	}
	
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", appError.NewServerError("failed to get user", err)
	}
	
	return user, tokenID, nil
}

// RevokeRefreshToken invalidates a refresh token
func (s *TokenService) RevokeRefreshToken(ctx context.Context, tokenString string) error {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.GetJWTSecret()), nil
	})
	
	if err != nil || !token.Valid {
		return appError.NewUnauthorizedError("invalid refresh token")
	}
	
	// Extract token ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return appError.NewUnauthorizedError("invalid token claims")
	}
	
	tokenID, ok := claims["token_id"].(string)
	if !ok {
		return appError.NewUnauthorizedError("invalid token ID")
	}
	
	// Delete from Redis
	key := fmt.Sprintf("refresh_token:%s", tokenID)
	err = s.redis.Del(ctx, key)
	if err != nil {
		return appError.NewServerError("failed to revoke refresh token", err)
	}
	
	return nil
}

// RefreshTokens validates a refresh token and generates a new token pair
func (s *TokenService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	user, oldTokenID, err := s.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	
	// Generate new token pair
	tokenPair, err := s.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, appError.NewServerError("failed to generate new tokens", err)
	}
	
	// Delete old refresh token
	key := fmt.Sprintf("refresh_token:%s", oldTokenID)
	if err := s.redis.Del(ctx, key); err != nil {
		return nil, appError.NewServerError("failed to revoke old refresh token", err)
	}
	
	return tokenPair, nil
}
