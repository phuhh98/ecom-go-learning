package models

// User's addresses
type Address struct {
	BaseModel
	UserID      uint   `json:"user_id" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AddressLine string `json:"address_line"`
}
