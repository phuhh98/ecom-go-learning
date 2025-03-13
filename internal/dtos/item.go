package dtos

type ItemStatusDTO struct {
	Status string `json:"status" binding:"oneof=available sold"`
}

type CreateItemsDTO struct {
	ProductId uint `json:"-"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
	ItemStatusDTO
}

type UpdateItemDto struct {
	ItemId uint `json:"-"`
	ItemStatusDTO
}
