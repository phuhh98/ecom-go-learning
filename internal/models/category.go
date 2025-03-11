package models

// Category for products
type Category struct {
	BaseModel
	Name          string      `json:"name" gorm:"uniqueIndex;not null"`
	Description   string      `json:"description"`
	SubCategories []*Category `json:"sub_categories" gorm:"many2many:subcategories;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Products      []*Product  `json:"products" gorm:"many2many:product_category;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
