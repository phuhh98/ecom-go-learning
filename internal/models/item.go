package models

// Product items
type Item struct {
	BaseModel
	// an Item belongs to a Product
	ProductID uint    `json:"product_id"`
	Product   Product `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Order has many Items, nullable
	OrderID *uint `json:"order_id"`

	// Possible value for status: available, sold
	Status string `json:"status" gorm:"default:available"`
}
