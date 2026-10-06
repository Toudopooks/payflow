package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/Toudopooks/payflow/backend/internal/domain"
	"github.com/Toudopooks/payflow/backend/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	store repository.Store
	now   func() time.Time
}

func New(store repository.Store) *Service {
	return &Service{store: store, now: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) Ping(ctx context.Context) error { return s.store.Ping(ctx) }

func (s *Service) CreateCustomer(ctx context.Context, name, email string) (domain.Customer, error) {
	name, email = strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email))
	if name == "" {
		return domain.Customer{}, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.Customer{}, fmt.Errorf("%w: valid email is required", domain.ErrInvalidInput)
	}
	c := domain.Customer{ID: uuid.NewString(), Name: name, Email: email, CreatedAt: s.now()}
	return c, s.store.CreateCustomer(ctx, c)
}
func (s *Service) GetCustomer(ctx context.Context, id string) (domain.Customer, error) {
	return s.store.GetCustomer(ctx, id)
}
func (s *Service) CreateMerchant(ctx context.Context, name string) (domain.Merchant, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Merchant{}, fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	m := domain.Merchant{ID: uuid.NewString(), Name: name, CreatedAt: s.now()}
	return m, s.store.CreateMerchant(ctx, m)
}
func (s *Service) GetMerchant(ctx context.Context, id string) (domain.Merchant, error) {
	return s.store.GetMerchant(ctx, id)
}
func (s *Service) CreateOrder(ctx context.Context, customerID, merchantID string, amount int64, currency string) (domain.Order, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if amount <= 0 {
		return domain.Order{}, fmt.Errorf("%w: amount must be positive", domain.ErrInvalidInput)
	}
	if currency != "USD" && currency != "CLP" {
		return domain.Order{}, domain.ErrInvalidCurrency
	}
	id := uuid.NewString()
	o := domain.Order{ID: id, Reference: "PF-" + strings.ToUpper(strings.ReplaceAll(id[:8], "-", "")), CustomerID: customerID, MerchantID: merchantID, Amount: amount, Currency: currency, Status: "pending", CreatedAt: s.now(), UpdatedAt: s.now()}
	return o, s.store.CreateOrder(ctx, o)
}
func (s *Service) GetOrder(ctx context.Context, id string) (domain.Order, error) {
	return s.store.GetOrder(ctx, id)
}
func (s *Service) ListOrders(ctx context.Context) ([]domain.Order, error) {
	return s.store.ListOrders(ctx)
}
func (s *Service) CreatePlan(ctx context.Context, orderID string, count int) (domain.Plan, error) {
	o, err := s.store.GetOrder(ctx, orderID)
	if err != nil {
		return domain.Plan{}, err
	}
	id := uuid.NewString()
	items, err := BuildInstallments(id, o.Amount, count, s.now())
	if err != nil {
		return domain.Plan{}, err
	}
	p := domain.Plan{ID: id, OrderID: orderID, InstallmentCount: count, Status: "active", Installments: items, CreatedAt: s.now()}
	return p, s.store.CreatePlan(ctx, p, items)
}
func (s *Service) GetPlan(ctx context.Context, id string) (domain.Plan, error) {
	return s.store.GetPlan(ctx, id)
}
func (s *Service) GetOrderInstallments(ctx context.Context, id string) ([]domain.Installment, error) {
	return s.store.GetOrderInstallments(ctx, id)
}
func (s *Service) Pay(ctx context.Context, installmentID string, amount int64, key string, fail bool) (domain.Payment, bool, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return domain.Payment{}, false, fmt.Errorf("%w: Idempotency-Key header is required", domain.ErrInvalidInput)
	}
	if len(key) > 128 {
		return domain.Payment{}, false, fmt.Errorf("%w: idempotency key is too long", domain.ErrInvalidInput)
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%t", installmentID, amount, fail)))
	return s.store.CreatePayment(ctx, repository.PaymentCommand{ID: uuid.NewString(), InstallmentID: installmentID, Amount: amount, IdempotencyKey: key, RequestHash: hex.EncodeToString(sum[:]), SimulateFailure: fail, Now: s.now()})
}
func (s *Service) GetPayment(ctx context.Context, id string) (domain.Payment, error) {
	return s.store.GetPayment(ctx, id)
}
func (s *Service) Refund(ctx context.Context, paymentID string, amount int64, reason string) (domain.Refund, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return domain.Refund{}, fmt.Errorf("%w: reason is required", domain.ErrInvalidInput)
	}
	return s.store.CreateRefund(ctx, repository.RefundCommand{ID: uuid.NewString(), PaymentID: paymentID, Amount: amount, Reason: reason, Now: s.now()})
}
func (s *Service) Transactions(ctx context.Context, customerID string, page, pageSize int) (domain.Page[domain.TransactionEvent], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.store.CustomerTransactions(ctx, customerID, page, pageSize)
}
