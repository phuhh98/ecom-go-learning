package dtos

// DOCCategoryResponse represents a category response for swagger
// @Description Category data response
type DOCCategoryResponse struct {
	ID            uint                 `json:"id" example:"1"`
	Name          string               `json:"name" example:"Electronics"`
	Description   string               `json:"description" example:"Electronic devices and gadgets"`
	Products      []DOCProductSummary  `json:"products,omitempty"`
	Subcategories []DOCCategorySummary `json:"subcategories,omitempty"`
	ParentID      *uint                `json:"parent_id,omitempty" example:"0"`
}

// DOCCategorySummary represents a simplified category for swagger
// @Description Summary of category information
type DOCCategorySummary struct {
	ID          uint   `json:"id" example:"2"`
	Name        string `json:"name" example:"Smartphones"`
	Description string `json:"description" example:"Mobile phones and accessories"`
}

// DOCCategoryListResponse represents a list of categories for swagger
// @Description List of categories
type DOCCategoryListResponse struct {
	Data []DOCCategoryResponse `json:"data"`
}
