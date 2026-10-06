package domain

import "time"

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Merchant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID          string    `json:"id"`
	Reference   string    `json:"reference"`
	CustomerID  string    `json:"customer_id"`
	MerchantID  string    `json:"merchant_id"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	PaymentPlan *Plan     `json:"payment_plan,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Plan struct {
	ID               string        `json:"id"`
	OrderID          string        `json:"order_id"`
	InstallmentCount int           `json:"installment_count"`
	Status           string        `json:"status"`
	Installments     []Installment `json:"installments,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
}

type Installment struct {
	ID        string     `json:"id"`
	PlanID    string     `json:"payment_plan_id"`
	Sequence  int        `json:"sequence"`
	Amount    int64      `json:"amount"`
	Status    string     `json:"status"`
	DueAt     time.Time  `json:"due_at"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Payment struct {
	ID            string    `json:"id"`
	InstallmentID string    `json:"installment_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	FailureReason string    `json:"failure_reason,omitempty"`
	Refunded      int64     `json:"refunded_amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type Refund struct {
	ID        string    `json:"id"`
	PaymentID string    `json:"payment_id"`
	Amount    int64     `json:"amount"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionEvent struct {
	ID          string         `json:"id"`
	CustomerID  string         `json:"customer_id"`
	OrderID     string         `json:"order_id"`
	Type        string         `json:"type"`
	Status      string         `json:"status"`
	Amount      int64          `json:"amount"`
	Currency    string         `json:"currency"`
	ReferenceID string         `json:"reference_id"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Page[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
