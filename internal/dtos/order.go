package dtos

type CreateOrderDTO struct {
	UserID uint           `json:"user_id" binding:"required"`
	Status string         `json:"status" binding:"required,oneof=created processed intransit complete cancelled"`
	Items  []OrderItemDTO `json:"items" binding:"required"`
}

type OrderItemDTO struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}

type UpdateOrderDTO struct {
	Status string `json:"status" binding:"required,oneof=created processed intransit complete cancelled"`
	Items  []uint `json:"items" binding:"required"`
}
