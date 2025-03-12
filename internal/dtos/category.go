package dtos

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required,max=2048"`
	Description string `json:"description" binding:"max=50000"`
	// SubCategories []uint `json:"sub_categories"`
	// Products      []uint `json:"products"`
}

type UpdateCategoryDTO struct {
	CreateCategoryDTO
}
