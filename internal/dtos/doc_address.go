package dtos

// DOCErrorResponse represents an error response for swagger
// @Description Error response with message and details
type DOCErrorResponse struct {
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"Bad request"`
	Details interface{} `json:"details,omitempty"`
}

// DOCAddressResponse represents an address response for swagger
// @Description Address data response
type DOCAddressResponse struct {
	ID         uint   `json:"id" example:"1"`
	UserID     uint   `json:"user_id" example:"1"`
	Name       string `json:"name" example:"Home"`
	Phone      string `json:"phone" example:"+1234567890"`
	Street     string `json:"street" example:"123 Main St"`
	City       string `json:"city" example:"New York"`
	State      string `json:"state" example:"NY"`
	PostalCode string `json:"postal_code" example:"10001"`
	Country    string `json:"country" example:"USA"`
	IsDefault  bool   `json:"is_default" example:"true"`
}

// DOCAddressListResponse represents a paginated list of addresses for swagger
// @Description Paginated list of addresses
type DOCAddressListResponse struct {
	Data []DOCAddressResponse `json:"data"`
	Meta DOCPaginationMeta    `json:"meta"`
}

// DOCMessageResponse represents a simple message response for swagger
// @Description Simple message response
type DOCMessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
