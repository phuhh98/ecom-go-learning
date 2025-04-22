package models

type Order struct {
	BaseModel
	UserID    uint         `json:"user_id" gorm:"not null;"`
	Status    string       `json:"status" gorm:"not null;"` // Use enum: created, processed, intransit, complete, cancelled
	Items     []*OrderItem `json:"items"`
	AddressID uint         `json:"-" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Address   Address      `json:"address" gorm:"foreignKey:AddressID"`
}

type OrderItem struct {
	BaseModel
	OrderID   uint    `json:"-" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	ProductID uint    `json:"-" gorm:"not null;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Product   Product `json:"product" gorm:"not null;"`
	Quantity  uint    `json:"quantity"`
}
