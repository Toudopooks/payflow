package domain

import "errors"

var (
	ErrNotFound                = errors.New("resource not found")
	ErrInvalidInput            = errors.New("invalid input")
	ErrInvalidCurrency         = errors.New("unsupported currency")
	ErrInvalidInstallmentCount = errors.New("installment count must be 2, 3, or 4")
	ErrInvalidPaymentAmount    = errors.New("payment amount must equal the installment amount")
	ErrInstallmentAlreadyPaid  = errors.New("installment is already paid")
	ErrPaymentFailed           = errors.New("payment was intentionally declined")
	ErrPaymentNotRefundable    = errors.New("payment cannot be refunded")
	ErrInvalidRefundAmount     = errors.New("refund exceeds the refundable amount")
	ErrIdempotencyConflict     = errors.New("idempotency key was reused with a different request")
	ErrIdempotencyInProgress   = errors.New("idempotent request is still processing")
)
