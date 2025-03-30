package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"userId"`
	Status    string    `json:"status"` // Use enum: created, processed, intransit, complete
	Items     []*Item   `json:"items" gorm:"-"` // Virtual field, not stored in database
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
