package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	BaseModel
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"` // Password is not exposed in JSON
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" gorm:"default:user"`
	// Cascade delete of addresses on user delete, nullable address
	Addresses []*Address `json:"addresses"`
}

// HashPassword encrypts the password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword compares the given password with the stored hash.
func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}

// GenerateJWT generates a JWT access token for the user.
// The token contains the user ID, email, role, and expires based on the provided expiration minutes.
func (u *User) GenerateJWT(jwtSecret string, expirationMinutes int) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"role":    u.Role,
		"exp":     time.Now().Add(time.Minute * time.Duration(expirationMinutes)).Unix(),
		"type":    "access",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a JWT refresh token for the user.
// The token contains the user ID, token ID, and expires based on the provided expiration days.
func (u *User) GenerateRefreshToken(jwtSecret string, tokenID string, expirationDays int) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"user_id":  u.ID,
		"token_id": tokenID,
		"exp":      time.Now().Add(time.Hour * 24 * time.Duration(expirationDays)).Unix(),
		"type":     "refresh",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
