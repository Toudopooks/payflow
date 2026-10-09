package integration

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DevPooks/payflow/backend/internal/domain"
	mysqlstore "github.com/DevPooks/payflow/backend/internal/repository/mysql"
	"github.com/DevPooks/payflow/backend/internal/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func TestPaymentIdempotencyAndRefund(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}
	ctx := context.Background()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	deadline := time.Now().Add(60 * time.Second)
	for {
		if err = db.PingContext(ctx); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("wait for MySQL: %v", err)
		}
		time.Sleep(time.Second)
	}
	migration, err := os.ReadFile(filepath.Join("..", "..", "migrations", "000001_initial.up.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.ExecContext(ctx, string(migration)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	store, err := mysqlstore.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	svc := service.New(store)
	customer, err := svc.CreateCustomer(ctx, "Maya Rivera", uuid.NewString()+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	merchant, err := svc.CreateMerchant(ctx, "North Harbor Goods")
	if err != nil {
		t.Fatal(err)
	}
	order, err := svc.CreateOrder(ctx, customer.ID, merchant.ID, 10000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := svc.CreatePlan(ctx, order.ID, 3)
	if err != nil {
		t.Fatal(err)
	}
	first := plan.Installments[0]
	if _, _, err = svc.Pay(ctx, first.ID, first.Amount+1, "rollback-check", false); err != domain.ErrInvalidPaymentAmount {
		t.Fatalf("invalid payment: got %v", err)
	}
	var rolledBackKeys, rolledBackPayments int
	if err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM idempotency_keys WHERE idempotency_key='rollback-check'`).Scan(&rolledBackKeys); err != nil {
		t.Fatal(err)
	}
	if err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments WHERE installment_id=?`, first.ID).Scan(&rolledBackPayments); err != nil {
		t.Fatal(err)
	}
	if rolledBackKeys != 0 || rolledBackPayments != 0 {
		t.Fatalf("failed transaction left %d keys and %d payments", rolledBackKeys, rolledBackPayments)
	}
	paid, replay, err := svc.Pay(ctx, first.ID, first.Amount, "integration-payment-1", false)
	if err != nil {
		t.Fatal(err)
	}
	if replay {
		t.Fatal("first request was marked replay")
	}
	replayed, isReplay, err := svc.Pay(ctx, first.ID, first.Amount, "integration-payment-1", false)
	if err != nil {
		t.Fatal(err)
	}
	if !isReplay || replayed.ID != paid.ID {
		t.Fatalf("idempotent replay returned a different payment")
	}
	if _, _, err = svc.Pay(ctx, first.ID, first.Amount, "integration-payment-2", false); err != domain.ErrInstallmentAlreadyPaid {
		t.Fatalf("second payment: got %v", err)
	}
	refund, err := svc.Refund(ctx, paid.ID, paid.Amount/2, "Customer requested a partial refund")
	if err != nil {
		t.Fatal(err)
	}
	if refund.Amount != paid.Amount/2 {
		t.Fatal("wrong refund amount")
	}
	if _, err = svc.Refund(ctx, paid.ID, paid.Amount-refund.Amount, "Complete the refund"); err != nil {
		t.Fatal(err)
	}
	refundedPayment, err := svc.GetPayment(ctx, paid.ID)
	if err != nil || refundedPayment.Status != "refunded" {
		t.Fatalf("payment status after full refund: %q, %v", refundedPayment.Status, err)
	}
	activeOrder, err := svc.GetOrder(ctx, order.ID)
	if err != nil || activeOrder.Status != "active" {
		t.Fatalf("order with pending installments should stay active: %q, %v", activeOrder.Status, err)
	}
	if _, err = svc.Refund(ctx, paid.ID, 1, "over refund"); err != domain.ErrInvalidRefundAmount {
		t.Fatalf("over refund: got %v", err)
	}
	second := plan.Installments[1]
	failed, _, err := svc.Pay(ctx, second.ID, second.Amount, "integration-decline-1", true)
	if err != nil || failed.Status != "failed" {
		t.Fatalf("failed payment result: %q, %v", failed.Status, err)
	}
	page, err := svc.Transactions(ctx, customer.ID, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total < 6 {
		t.Fatalf("got %d transaction events", page.Total)
	}
}
