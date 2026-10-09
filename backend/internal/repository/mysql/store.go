package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/DevPooks/payflow/backend/internal/domain"
	"github.com/DevPooks/payflow/backend/internal/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Store struct{ db *sql.DB }

func Open(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &Store{db: db}, nil
}

func (s *Store) Close() error                   { return s.db.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

func (s *Store) CreateCustomer(ctx context.Context, c domain.Customer) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO customers (id,name,email) VALUES (?,?,?)`, c.ID, c.Name, c.Email)
	return translate(err)
}

func (s *Store) GetCustomer(ctx context.Context, id string) (domain.Customer, error) {
	var c domain.Customer
	err := s.db.QueryRowContext(ctx, `SELECT id,name,email,created_at FROM customers WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	return c, translate(err)
}

func (s *Store) CreateMerchant(ctx context.Context, m domain.Merchant) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO merchants (id,name) VALUES (?,?)`, m.ID, m.Name)
	return translate(err)
}

func (s *Store) GetMerchant(ctx context.Context, id string) (domain.Merchant, error) {
	var m domain.Merchant
	err := s.db.QueryRowContext(ctx, `SELECT id,name,created_at FROM merchants WHERE id=?`, id).Scan(&m.ID, &m.Name, &m.CreatedAt)
	return m, translate(err)
}

func (s *Store) CreateOrder(ctx context.Context, o domain.Order) error {
	return s.withTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO orders (id,reference,customer_id,merchant_id,amount,currency,status) VALUES (?,?,?,?,?,?,?)`, o.ID, o.Reference, o.CustomerID, o.MerchantID, o.Amount, o.Currency, o.Status)
		if err != nil {
			return translate(err)
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO transaction_events (id,customer_id,order_id,event_type,status,amount,currency,reference_id) VALUES (?,?,?,?,?,?,?,?)`, uuid.NewString(), o.CustomerID, o.ID, "order_created", o.Status, o.Amount, o.Currency, o.ID)
		return translate(err)
	})
}

func (s *Store) GetOrder(ctx context.Context, id string) (domain.Order, error) {
	var o domain.Order
	err := s.db.QueryRowContext(ctx, `SELECT id,reference,customer_id,merchant_id,amount,currency,status,created_at,updated_at FROM orders WHERE id=?`, id).Scan(&o.ID, &o.Reference, &o.CustomerID, &o.MerchantID, &o.Amount, &o.Currency, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return o, translate(err)
	}
	var planID string
	err = s.db.QueryRowContext(ctx, `SELECT id FROM payment_plans WHERE order_id=?`, id).Scan(&planID)
	if err == nil {
		p, getErr := s.GetPlan(ctx, planID)
		if getErr != nil {
			return o, getErr
		}
		o.PaymentPlan = &p
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return o, err
	}
	return o, nil
}

func (s *Store) ListOrders(ctx context.Context) ([]domain.Order, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,reference,customer_id,merchant_id,amount,currency,status,created_at,updated_at FROM orders ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.Reference, &o.CustomerID, &o.MerchantID, &o.Amount, &o.Currency, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (s *Store) CreatePlan(ctx context.Context, p domain.Plan, installments []domain.Installment) error {
	return s.withTx(ctx, func(tx *sql.Tx) error {
		var customerID, currency string
		var amount int64
		if err := tx.QueryRowContext(ctx, `SELECT customer_id,currency,amount FROM orders WHERE id=? FOR UPDATE`, p.OrderID).Scan(&customerID, &currency, &amount); err != nil {
			return translate(err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO payment_plans (id,order_id,installment_count,status) VALUES (?,?,?,?)`, p.ID, p.OrderID, p.InstallmentCount, p.Status); err != nil {
			return translate(err)
		}
		for _, item := range installments {
			if _, err := tx.ExecContext(ctx, `INSERT INTO installments (id,payment_plan_id,sequence_number,amount,status,due_at) VALUES (?,?,?,?,?,?)`, item.ID, item.PlanID, item.Sequence, item.Amount, item.Status, item.DueAt); err != nil {
				return translate(err)
			}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE orders SET status='active' WHERE id=?`, p.OrderID); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, `INSERT INTO transaction_events (id,customer_id,order_id,event_type,status,amount,currency,reference_id) VALUES (?,?,?,?,?,?,?,?)`, uuid.NewString(), customerID, p.OrderID, "plan_created", p.Status, amount, currency, p.ID)
		return translate(err)
	})
}

func (s *Store) GetPlan(ctx context.Context, id string) (domain.Plan, error) {
	var p domain.Plan
	err := s.db.QueryRowContext(ctx, `SELECT id,order_id,installment_count,status,created_at FROM payment_plans WHERE id=?`, id).Scan(&p.ID, &p.OrderID, &p.InstallmentCount, &p.Status, &p.CreatedAt)
	if err != nil {
		return p, translate(err)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,payment_plan_id,sequence_number,amount,status,due_at,paid_at,created_at FROM installments WHERE payment_plan_id=? ORDER BY sequence_number`, id)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var i domain.Installment
		if err := rows.Scan(&i.ID, &i.PlanID, &i.Sequence, &i.Amount, &i.Status, &i.DueAt, &i.PaidAt, &i.CreatedAt); err != nil {
			return p, err
		}
		p.Installments = append(p.Installments, i)
	}
	return p, rows.Err()
}

func (s *Store) GetOrderInstallments(ctx context.Context, orderID string) ([]domain.Installment, error) {
	var planID string
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM payment_plans WHERE order_id=?`, orderID).Scan(&planID); err != nil {
		return nil, translate(err)
	}
	p, err := s.GetPlan(ctx, planID)
	return p.Installments, err
}

func (s *Store) CreatePayment(ctx context.Context, cmd repository.PaymentCommand) (domain.Payment, bool, error) {
	var result domain.Payment
	replayed := false
	err := s.withTx(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `INSERT IGNORE INTO idempotency_keys (idempotency_key,request_hash,resource_type,expires_at) VALUES (?,?, 'payment', ?)`, cmd.IdempotencyKey, cmd.RequestHash, cmd.Now.Add(24*time.Hour))
		if err != nil {
			return err
		}
		inserted, _ := res.RowsAffected()
		if inserted == 0 {
			var hash, state string
			var resourceID sql.NullString
			if err := tx.QueryRowContext(ctx, `SELECT request_hash,state,resource_id FROM idempotency_keys WHERE idempotency_key=? FOR UPDATE`, cmd.IdempotencyKey).Scan(&hash, &state, &resourceID); err != nil {
				return err
			}
			if hash != cmd.RequestHash {
				return domain.ErrIdempotencyConflict
			}
			if state != "completed" || !resourceID.Valid {
				return domain.ErrIdempotencyInProgress
			}
			replayed = true
			return scanPayment(tx.QueryRowContext(ctx, `SELECT id,installment_id,amount,currency,status,failure_reason,refunded_amount,created_at FROM payments WHERE id=?`, resourceID.String), &result)
		}
		var installment domain.Installment
		var orderID, customerID, currency string
		err = tx.QueryRowContext(ctx, `SELECT i.id,i.payment_plan_id,i.sequence_number,i.amount,i.status,i.due_at,i.paid_at,i.created_at,p.order_id,o.customer_id,o.currency FROM installments i JOIN payment_plans p ON p.id=i.payment_plan_id JOIN orders o ON o.id=p.order_id WHERE i.id=? FOR UPDATE`, cmd.InstallmentID).Scan(&installment.ID, &installment.PlanID, &installment.Sequence, &installment.Amount, &installment.Status, &installment.DueAt, &installment.PaidAt, &installment.CreatedAt, &orderID, &customerID, &currency)
		if err != nil {
			return translate(err)
		}
		if installment.Status == "paid" || installment.Status == "refunded" {
			return domain.ErrInstallmentAlreadyPaid
		}
		if cmd.Amount != installment.Amount {
			return domain.ErrInvalidPaymentAmount
		}
		status, reason, eventType := "succeeded", "", "payment_succeeded"
		if cmd.SimulateFailure {
			status, reason, eventType = "failed", "demo decline requested", "payment_failed"
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO payments (id,installment_id,amount,currency,status,failure_reason) VALUES (?,?,?,?,?,?)`, cmd.ID, cmd.InstallmentID, cmd.Amount, currency, status, reason); err != nil {
			return err
		}
		if status == "succeeded" {
			if _, err := tx.ExecContext(ctx, `UPDATE installments SET status='paid',paid_at=? WHERE id=?`, cmd.Now, cmd.InstallmentID); err != nil {
				return err
			}
		}
		if status == "failed" {
			if _, err := tx.ExecContext(ctx, `UPDATE installments SET status='failed' WHERE id=?`, cmd.InstallmentID); err != nil {
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO transaction_events (id,customer_id,order_id,event_type,status,amount,currency,reference_id) VALUES (?,?,?,?,?,?,?,?)`, uuid.NewString(), customerID, orderID, eventType, status, cmd.Amount, currency, cmd.ID); err != nil {
			return err
		}
		if status == "succeeded" {
			var pending int
			if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM installments WHERE payment_plan_id=? AND status!='paid'`, installment.PlanID).Scan(&pending); err != nil {
				return err
			}
			if pending == 0 {
				if _, err := tx.ExecContext(ctx, `UPDATE payment_plans SET status='completed' WHERE id=?`, installment.PlanID); err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, `UPDATE orders SET status='completed' WHERE id=?`, orderID); err != nil {
					return err
				}
			}
		}
		if _, err := tx.ExecContext(ctx, `UPDATE idempotency_keys SET state='completed',resource_id=?,response_code=201 WHERE idempotency_key=?`, cmd.ID, cmd.IdempotencyKey); err != nil {
			return err
		}
		result = domain.Payment{ID: cmd.ID, InstallmentID: cmd.InstallmentID, Amount: cmd.Amount, Currency: currency, Status: status, FailureReason: reason, CreatedAt: cmd.Now}
		return nil
	})
	return result, replayed, err
}

func (s *Store) GetPayment(ctx context.Context, id string) (domain.Payment, error) {
	var p domain.Payment
	err := scanPayment(s.db.QueryRowContext(ctx, `SELECT id,installment_id,amount,currency,status,failure_reason,refunded_amount,created_at FROM payments WHERE id=?`, id), &p)
	return p, translate(err)
}

type rowScanner interface{ Scan(...any) error }

func scanPayment(row rowScanner, p *domain.Payment) error {
	return row.Scan(&p.ID, &p.InstallmentID, &p.Amount, &p.Currency, &p.Status, &p.FailureReason, &p.Refunded, &p.CreatedAt)
}

func (s *Store) CreateRefund(ctx context.Context, cmd repository.RefundCommand) (domain.Refund, error) {
	var refund domain.Refund
	err := s.withTx(ctx, func(tx *sql.Tx) error {
		var payment domain.Payment
		var orderID, customerID, planID string
		if err := tx.QueryRowContext(ctx, `SELECT py.id,py.installment_id,py.amount,py.currency,py.status,py.failure_reason,py.refunded_amount,py.created_at,p.order_id,o.customer_id,p.id FROM payments py JOIN installments i ON i.id=py.installment_id JOIN payment_plans p ON p.id=i.payment_plan_id JOIN orders o ON o.id=p.order_id WHERE py.id=? FOR UPDATE`, cmd.PaymentID).Scan(&payment.ID, &payment.InstallmentID, &payment.Amount, &payment.Currency, &payment.Status, &payment.FailureReason, &payment.Refunded, &payment.CreatedAt, &orderID, &customerID, &planID); err != nil {
			return translate(err)
		}
		if payment.Status != "succeeded" && payment.Status != "refunded" {
			return domain.ErrPaymentNotRefundable
		}
		if cmd.Amount <= 0 || cmd.Amount > payment.Amount-payment.Refunded {
			return domain.ErrInvalidRefundAmount
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO refunds (id,payment_id,amount,reason,status) VALUES (?,?,?,?,'succeeded')`, cmd.ID, cmd.PaymentID, cmd.Amount, cmd.Reason); err != nil {
			return err
		}
		newRefunded := payment.Refunded + cmd.Amount
		status := "succeeded"
		if newRefunded == payment.Amount {
			status = "refunded"
		}
		if _, err := tx.ExecContext(ctx, `UPDATE payments SET refunded_amount=?,status=? WHERE id=?`, newRefunded, status, cmd.PaymentID); err != nil {
			return err
		}
		if status == "refunded" {
			if _, err := tx.ExecContext(ctx, `UPDATE installments SET status='refunded' WHERE id=?`, payment.InstallmentID); err != nil {
				return err
			}
		}
		var notRefunded int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM installments WHERE payment_plan_id=? AND status!='refunded'`, planID).Scan(&notRefunded); err != nil {
			return err
		}
		if notRefunded == 0 {
			if _, err := tx.ExecContext(ctx, `UPDATE payment_plans SET status='refunded' WHERE id=?`, planID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE orders SET status='refunded' WHERE id=?`, orderID); err != nil {
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO transaction_events (id,customer_id,order_id,event_type,status,amount,currency,reference_id) VALUES (?,?,?,?,?,?,?,?)`, uuid.NewString(), customerID, orderID, "refund_succeeded", "succeeded", -cmd.Amount, payment.Currency, cmd.ID); err != nil {
			return err
		}
		refund = domain.Refund{ID: cmd.ID, PaymentID: cmd.PaymentID, Amount: cmd.Amount, Reason: cmd.Reason, Status: "succeeded", CreatedAt: cmd.Now}
		return nil
	})
	return refund, err
}

func (s *Store) CustomerTransactions(ctx context.Context, customerID string, page, pageSize int) (domain.Page[domain.TransactionEvent], error) {
	var total int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transaction_events WHERE customer_id=?`, customerID).Scan(&total); err != nil {
		return domain.Page[domain.TransactionEvent]{}, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,customer_id,order_id,event_type,status,amount,currency,reference_id,metadata,created_at FROM transaction_events WHERE customer_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?`, customerID, pageSize, (page-1)*pageSize)
	if err != nil {
		return domain.Page[domain.TransactionEvent]{}, err
	}
	defer rows.Close()
	items := []domain.TransactionEvent{}
	for rows.Next() {
		var e domain.TransactionEvent
		var raw []byte
		if err := rows.Scan(&e.ID, &e.CustomerID, &e.OrderID, &e.Type, &e.Status, &e.Amount, &e.Currency, &e.ReferenceID, &raw, &e.CreatedAt); err != nil {
			return domain.Page[domain.TransactionEvent]{}, err
		}
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &e.Metadata)
		}
		items = append(items, e)
	}
	return domain.Page[domain.TransactionEvent]{Data: items, Page: page, PageSize: pageSize, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, rows.Err()
}

func (s *Store) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
func translate(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
