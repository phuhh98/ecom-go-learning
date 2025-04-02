package dtos

// CreateUserDTO represents the input for creating a new user
type CreateUserDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserDTO represents the input for updating a user
type UpdateUserDTO struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role string `json:"role" binding:"oneof=admin user"` // Role can be either "admin" or "user"
}

// LoginDTO represents the input for user login
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
