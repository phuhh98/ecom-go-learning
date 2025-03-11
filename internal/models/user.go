package models

import (
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

// CheckPassword compares the given password with the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
