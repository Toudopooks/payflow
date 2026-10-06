package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Toudopooks/payflow/backend/internal/domain"
	appmiddleware "github.com/Toudopooks/payflow/backend/internal/middleware"
	"github.com/Toudopooks/payflow/backend/internal/observability"
	"github.com/Toudopooks/payflow/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	service  *service.Service
	metrics  *observability.Metrics
	demoMode bool
}

func New(service *service.Service, metrics *observability.Metrics, demoMode bool, allowedOrigins []string, observe func(http.Handler) http.Handler) http.Handler {
	a := &API{service: service, metrics: metrics, demoMode: demoMode}
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, appmiddleware.RequestIDMiddleware, appmiddleware.SecurityHeaders, appmiddleware.CORS(allowedOrigins), middleware.RequestSize(1<<20), observe)
	r.Get("/health", a.health)
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/customers", a.createCustomer)
		r.Get("/customers/{id}", a.getCustomer)
		r.Get("/customers/{id}/transactions", a.transactions)
		r.Post("/merchants", a.createMerchant)
		r.Get("/merchants/{id}", a.getMerchant)
		r.Post("/orders", a.createOrder)
		r.Get("/orders", a.listOrders)
		r.Get("/orders/{id}", a.getOrder)
		r.Get("/orders/{id}/installments", a.orderInstallments)
		r.Post("/payment-plans", a.createPlan)
		r.Get("/payment-plans/{id}", a.getPlan)
		r.Post("/payments", a.createPayment)
		r.Post("/installments/{id}/pay", a.payInstallment)
		r.Get("/payments/{id}", a.getPayment)
		r.Post("/payments/{id}/refunds", a.refund)
	})
	return r
}

type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func decode(r *http.Request, target any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("%w: request body must be valid JSON", domain.ErrInvalidInput)
	}
	return nil
}
func writeError(w http.ResponseWriter, err error) {
	status, code, msg := http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed."
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status, code, msg = http.StatusNotFound, "NOT_FOUND", "The requested resource was not found."
	case errors.Is(err, domain.ErrIdempotencyConflict):
		status, code, msg = http.StatusConflict, "IDEMPOTENCY_CONFLICT", err.Error()
	case errors.Is(err, domain.ErrIdempotencyInProgress):
		status, code, msg = http.StatusConflict, "IDEMPOTENCY_IN_PROGRESS", err.Error()
	case errors.Is(err, domain.ErrInstallmentAlreadyPaid):
		status, code, msg = http.StatusConflict, "INSTALLMENT_ALREADY_PAID", err.Error()
	case errors.Is(err, domain.ErrInvalidPaymentAmount):
		status, code, msg = http.StatusUnprocessableEntity, "INVALID_PAYMENT_AMOUNT", err.Error()
	case errors.Is(err, domain.ErrInvalidRefundAmount):
		status, code, msg = http.StatusUnprocessableEntity, "INVALID_REFUND_AMOUNT", err.Error()
	case errors.Is(err, domain.ErrPaymentNotRefundable):
		status, code, msg = http.StatusConflict, "PAYMENT_NOT_REFUNDABLE", err.Error()
	case errors.Is(err, domain.ErrInvalidCurrency), errors.Is(err, domain.ErrInvalidInstallmentCount), errors.Is(err, domain.ErrInvalidInput):
		status, code, msg = http.StatusBadRequest, "INVALID_REQUEST", err.Error()
	}
	body := errorBody{}
	body.Error.Code = code
	body.Error.Message = msg
	writeJSON(w, status, body)
}

func (a *API) fail(w http.ResponseWriter, err error) {
	if !errors.Is(err, domain.ErrNotFound) && !errors.Is(err, domain.ErrInvalidInput) && !errors.Is(err, domain.ErrInvalidCurrency) && !errors.Is(err, domain.ErrInvalidInstallmentCount) && !errors.Is(err, domain.ErrInvalidPaymentAmount) && !errors.Is(err, domain.ErrInstallmentAlreadyPaid) && !errors.Is(err, domain.ErrPaymentNotRefundable) && !errors.Is(err, domain.ErrInvalidRefundAmount) && !errors.Is(err, domain.ErrIdempotencyConflict) && !errors.Is(err, domain.ErrIdempotencyInProgress) {
		a.metrics.DatabaseErrors.Inc()
	}
	writeError(w, err)
}
func (a *API) health(w http.ResponseWriter, r *http.Request) {
	db := "ok"
	status := http.StatusOK
	if err := a.service.Ping(r.Context()); err != nil {
		db = "unavailable"
		status = http.StatusServiceUnavailable
		a.metrics.DatabaseErrors.Inc()
	}
	writeJSON(w, status, map[string]string{"status": map[bool]string{true: "ok", false: "degraded"}[status == http.StatusOK], "database": db})
}
func (a *API) createCustomer(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	v, err := a.service.CreateCustomer(r.Context(), in.Name, in.Email)
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}
func (a *API) getCustomer(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetCustomer(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (a *API) createMerchant(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	v, err := a.service.CreateMerchant(r.Context(), in.Name)
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}
func (a *API) getMerchant(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetMerchant(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var in struct {
		CustomerID string `json:"customer_id"`
		MerchantID string `json:"merchant_id"`
		Amount     int64  `json:"amount"`
		Currency   string `json:"currency"`
	}
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	v, err := a.service.CreateOrder(r.Context(), strings.TrimSpace(in.CustomerID), strings.TrimSpace(in.MerchantID), in.Amount, in.Currency)
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}
func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetOrder(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (a *API) listOrders(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.ListOrders(r.Context())
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": v})
}
func (a *API) createPlan(w http.ResponseWriter, r *http.Request) {
	var in struct {
		OrderID          string `json:"order_id"`
		InstallmentCount int    `json:"installment_count"`
	}
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	v, err := a.service.CreatePlan(r.Context(), strings.TrimSpace(in.OrderID), in.InstallmentCount)
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}
func (a *API) getPlan(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetPlan(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (a *API) orderInstallments(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetOrderInstallments(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": v})
}

type paymentInput struct {
	InstallmentID   string `json:"installment_id"`
	Amount          int64  `json:"amount"`
	SimulateFailure bool   `json:"simulate_failure"`
}

func (a *API) createPayment(w http.ResponseWriter, r *http.Request) {
	var in paymentInput
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	a.processPayment(w, r, in)
}
func (a *API) payInstallment(w http.ResponseWriter, r *http.Request) {
	var in paymentInput
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	in.InstallmentID = chi.URLParam(r, "id")
	a.processPayment(w, r, in)
}
func (a *API) processPayment(w http.ResponseWriter, r *http.Request, in paymentInput) {
	if in.SimulateFailure && !a.demoMode {
		a.fail(w, fmt.Errorf("%w: payment failure simulation is disabled", domain.ErrInvalidInput))
		return
	}
	v, replay, err := a.service.Pay(r.Context(), strings.TrimSpace(in.InstallmentID), in.Amount, r.Header.Get("Idempotency-Key"), in.SimulateFailure)
	if err != nil {
		a.metrics.Payments.WithLabelValues("error").Inc()
		a.fail(w, err)
		return
	}
	if replay {
		a.metrics.IdempotencyReplays.Inc()
		w.Header().Set("Idempotency-Replayed", "true")
	} else {
		a.metrics.Payments.WithLabelValues(v.Status).Inc()
	}
	status := http.StatusCreated
	if replay {
		status = http.StatusOK
	}
	writeJSON(w, status, v)
}
func (a *API) getPayment(w http.ResponseWriter, r *http.Request) {
	v, err := a.service.GetPayment(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
func (a *API) refund(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Amount int64  `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := decode(r, &in); err != nil {
		a.fail(w, err)
		return
	}
	v, err := a.service.Refund(r.Context(), chi.URLParam(r, "id"), in.Amount, in.Reason)
	if err != nil {
		a.fail(w, err)
		return
	}
	a.metrics.Refunds.Inc()
	writeJSON(w, http.StatusCreated, v)
}
func (a *API) transactions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	v, err := a.service.Transactions(r.Context(), chi.URLParam(r, "id"), page, size)
	if err != nil {
		a.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
