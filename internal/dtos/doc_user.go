package dtos

// DOCUserDetailResponse represents a detailed user response for swagger
// @Description Detailed user data
type DOCUserDetailResponse struct {
	ID        uint     `json:"id" example:"1"`
	Name      string   `json:"name" example:"John Doe"`
	Email     string   `json:"email" example:"john.doe@example.com"`
	Roles     []string `json:"roles" example:"[\"user\"]"`
	Verified  bool     `json:"verified" example:"true"`
	CreatedAt string   `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt string   `json:"updated_at" example:"2023-01-01T12:30:00Z"`
}

// DOCUserListItem represents a user summary for listings in swagger
// @Description Summary of user information for listings
type DOCUserListItem struct {
	ID       uint     `json:"id" example:"1"`
	Name     string   `json:"name" example:"John Doe"`
	Email    string   `json:"email" example:"john.doe@example.com"`
	Roles    []string `json:"roles" example:"[\"user\"]"`
	Verified bool     `json:"verified" example:"true"`
}
