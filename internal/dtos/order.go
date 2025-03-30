package dtos

type OrderDTO struct {
	UserID uint   `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=created processed intransit complete"`
	Items  []uint `json:"items" binding:"required"`
}
