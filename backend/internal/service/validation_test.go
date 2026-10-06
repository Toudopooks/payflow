package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Toudopooks/payflow/backend/internal/domain"
)

func TestValidationRejectsInvalidPublicInput(t *testing.T) {
	svc := New(nil)
	tests := []struct {
		name string
		run  func() error
		want error
	}{
		{"customer email", func() error { _, err := svc.CreateCustomer(context.Background(), "Maya", "not-an-email"); return err }, domain.ErrInvalidInput},
		{"currency", func() error {
			_, err := svc.CreateOrder(context.Background(), "customer", "merchant", 1000, "EUR")
			return err
		}, domain.ErrInvalidCurrency},
		{"negative order", func() error {
			_, err := svc.CreateOrder(context.Background(), "customer", "merchant", -1, "USD")
			return err
		}, domain.ErrInvalidInput},
		{"missing idempotency key", func() error { _, _, err := svc.Pay(context.Background(), "installment", 1000, "", false); return err }, domain.ErrInvalidInput},
		{"missing refund reason", func() error { _, err := svc.Refund(context.Background(), "payment", 1000, " "); return err }, domain.ErrInvalidInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); !errors.Is(err, test.want) {
				t.Fatalf("got %v, want %v", err, test.want)
			}
		})
	}
}
