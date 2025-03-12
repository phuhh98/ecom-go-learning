package models

import (
	"strings"

	"gorm.io/gorm"
)

type Product struct {
	BaseModel
	Code        string      `json:"code" gorm:"uniqueIndex;not null"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Categories  []*Category `json:"categories" gorm:"many2many:product_category;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	p.Code = strings.ToLower(p.Code)
	return nil
}
