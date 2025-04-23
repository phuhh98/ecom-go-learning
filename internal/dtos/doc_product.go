package dtos

// DOCProductResponse represents a product response for swagger
// @Description Product data response
type DOCProductResponse struct {
	ID          uint    `json:"id" example:"1"`
	Name        string  `json:"name" example:"iPhone 13"`
	Description string  `json:"description" example:"Latest iPhone model with A15 chip"`
	Price       float64 `json:"price" example:"999.99"`
	Stock       int     `json:"stock" example:"100"`
	SKU         string  `json:"sku" example:"IPHN-13-128GB"`
	ImageURL    string  `json:"image_url" example:"https://example.com/images/iphone13.jpg"`
	Categories  []DOCCategorySummary `json:"categories,omitempty"`
	CreatedAt   string  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2023-01-01T12:30:00Z"`
}

// DOCProductSummary represents a simplified product for listings
// @Description Summary of product information
type DOCProductSummary struct {
	ID          uint    `json:"id" example:"1"`
	Name        string  `json:"name" example:"iPhone 13"`
	Description string  `json:"description" example:"Latest iPhone model with A15 chip"`
	Price       float64 `json:"price" example:"999.99"`
	ImageURL    string  `json:"image_url" example:"https://example.com/images/iphone13.jpg"`
}

// DOCProductListResponse represents a paginated list of products
// @Description Paginated list of products
type DOCProductListResponse struct {
	Data       []DOCProductSummary `json:"data"`
	Meta    DOCPaginationMeta `json:"meta"`}
