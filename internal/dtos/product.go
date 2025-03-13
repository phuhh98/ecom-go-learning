package dtos

type CreateProductDTO struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"max=2048"`
	Description string  `json:"description" binding:"max=50000"`
	Price       float64 `json:"price" binding:"gt=0"`
}

type UpdateProductDTO struct {
	CreateProductDTO
}
