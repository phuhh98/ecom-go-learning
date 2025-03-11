package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreation"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoCreation"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:"autoCreation"`
}
