package domain

import "time"

type Order struct {
	ID            string    `json:"id"`
	ChainID       string    `json:"chain_id"`
	BranchID      string    `json:"branch_id"`
	TableID       *string   `json:"table_id,omitempty"`
	CustomerID    *string   `json:"customer_id,omitempty"`
	CustomerName  string    `json:"customer_name,omitempty"`
	GuestID       string    `json:"guest_id,omitempty"`
	OrderType     string    `json:"order_type"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod *string   `json:"payment_method,omitempty"`
	TotalAmount   float64   `json:"total_amount"`
	Notes         string    `json:"notes,omitempty"`
	Items         []OrderItem `json:"items,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	MenuItemID  string  `json:"menu_item_id"`
	VariantID   *string `json:"variant_id,omitempty"`
	VariantName string  `json:"variant_name,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

type OrderStore interface {
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	ListByBranchID(branchID string, page, limit int) ([]Order, int, error)
	ListByChainID(chainID string, page, limit int) ([]Order, int, error)
	UpdateStatus(id, status string) error
	UpdatePayment(id, status, method string) error
}
