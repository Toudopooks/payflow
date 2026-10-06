package repository

import (
	"context"
	"time"

	"github.com/santi/payflow/backend/internal/domain"
)

type Store interface {
	Ping(context.Context) error
	CreateCustomer(context.Context, domain.Customer) error
	GetCustomer(context.Context, string) (domain.Customer, error)
	CreateMerchant(context.Context, domain.Merchant) error
	GetMerchant(context.Context, string) (domain.Merchant, error)
	CreateOrder(context.Context, domain.Order) error
	GetOrder(context.Context, string) (domain.Order, error)
	ListOrders(context.Context) ([]domain.Order, error)
	CreatePlan(context.Context, domain.Plan, []domain.Installment) error
	GetPlan(context.Context, string) (domain.Plan, error)
	GetOrderInstallments(context.Context, string) ([]domain.Installment, error)
	CreatePayment(context.Context, PaymentCommand) (domain.Payment, bool, error)
	GetPayment(context.Context, string) (domain.Payment, error)
	CreateRefund(context.Context, RefundCommand) (domain.Refund, error)
	CustomerTransactions(context.Context, string, int, int) (domain.Page[domain.TransactionEvent], error)
}

type PaymentCommand struct {
	ID              string
	InstallmentID   string
	Amount          int64
	IdempotencyKey  string
	RequestHash     string
	SimulateFailure bool
	Now             time.Time
}

type RefundCommand struct {
	ID        string
	PaymentID string
	Amount    int64
	Reason    string
	Now       time.Time
}
