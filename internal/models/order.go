package models

type Order struct {
	BaseModel
	UserID uint         `json:"userId" gorm:"not null;"`
	Status string       `json:"status" gorm:"not null;"` // Use enum: created, processed, intransit, complete, cancelled
	Items  []*OrderItem `json:"items"`
}

type OrderItem struct {
	BaseModel
	OrderID   uint    `json:"order_id" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	ProductID uint    `json:"productId" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Product   Product `json:"product" gorm:"not null;"`
	Quantity  uint    `json:"quantity"`
}
