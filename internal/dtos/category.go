package dtos

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,max=2048"`
	Description string `json:"description" binding:"max=50000"`
}

type UpdateCategoryDTO struct {
	CreateCategoryDTO
}
