package models

type Order struct {
	BaseModel
	// an Order belongs to an User
	UserId uint `json:"user_id"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Possible value for status: created, processed, intransit, complete
	Status string `json:"status"  gorm:"default:created"`

	Items []Item `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
