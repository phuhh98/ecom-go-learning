package dtos

// CreateAddressDTO represents the input for creating a new user
type CreateAddressDTO struct {
	UserID      uint   `json:"user_id"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	AddressLine string `json:"address_line" binding:"required"`
}

// UpdateAddressDTO represents the input for updating a user
type UpdateAddressDTO struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AddressLine string `json:"address_line"`
}
