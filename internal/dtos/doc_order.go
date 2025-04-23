package dtos

// DOCAddressListResponse represents a paginated list of addresses for swagger
// @Description Paginated list of orders
type DOCOrderListResponse struct {
	Data       []DOCOrderResponse `json:"data"`
	Meta       DOCPaginationMeta `json:"meta"`
}


// DOCOrderResponse represents an order response for swagger
// @Description Order data response
type DOCOrderResponse struct {
	ID            uint              `json:"id" example:"1"`
	UserID        uint              `json:"user_id" example:"1"`
	Status        string            `json:"status" example:"pending"`
	Total         float64           `json:"total" example:"99.99"`
	Items         []DOCOrderItem    `json:"items"`
	ShippingInfo  DOCOrderShipping  `json:"shipping_info"`
	PaymentInfo   DOCOrderPayment   `json:"payment_info"`
	CreatedAt     string            `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt     string            `json:"updated_at" example:"2023-01-01T12:30:00Z"`
}

// DOCOrderItem represents an order item for swagger
// @Description Order item details
type DOCOrderItem struct {
	ID        uint    `json:"id" example:"1"`
	ProductID uint    `json:"product_id" example:"5"`
	Name      string  `json:"name" example:"iPhone 13"`
	Price     float64 `json:"price" example:"999.99"`
	Quantity  int     `json:"quantity" example:"1"`
	Subtotal  float64 `json:"subtotal" example:"999.99"`
}

// DOCOrderShipping represents shipping information for swagger
// @Description Order shipping information
type DOCOrderShipping struct {
	AddressID   uint   `json:"address_id" example:"1"`
	FullName    string `json:"full_name" example:"John Doe"`
	AddressLine string `json:"address_line" example:"123 Main St, New York, NY 10001, USA"`
	PhoneNumber string `json:"phone_number" example:"+1234567890"`
	Method      string `json:"method" example:"standard"`
	TrackingID  string `json:"tracking_id,omitempty" example:"TRACK123456789"`
}

// DOCOrderPayment represents payment information for swagger
// @Description Order payment information
type DOCOrderPayment struct {
	Method      string  `json:"method" example:"credit_card"`
	Status      string  `json:"status" example:"completed"`
	Amount      float64 `json:"amount" example:"99.99"`
	Currency    string  `json:"currency" example:"USD"`
	PaymentID   string  `json:"payment_id,omitempty" example:"PAY123456789"`
	PaymentDate string  `json:"payment_date,omitempty" example:"2023-01-01T12:15:00Z"`
}

// DOCOrderSummary represents a simplified order for listings for swagger
// @Description Summary of order information
type DOCOrderSummary struct {
	ID         uint    `json:"id" example:"1"`
	Status     string  `json:"status" example:"pending"`
	Total      float64 `json:"total" example:"99.99"`
	ItemCount  int     `json:"item_count" example:"3"`
	CreatedAt  string  `json:"created_at" example:"2023-01-01T12:00:00Z"`
}
