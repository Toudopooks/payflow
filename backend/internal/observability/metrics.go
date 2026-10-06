package observability

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HTTPRequests       *prometheus.CounterVec
	HTTPDuration       *prometheus.HistogramVec
	Payments           *prometheus.CounterVec
	Refunds            prometheus.Counter
	DatabaseErrors     prometheus.Counter
	IdempotencyReplays prometheus.Counter
}

func New(registerer prometheus.Registerer) *Metrics {
	m := &Metrics{
		HTTPRequests:       prometheus.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total", Help: "HTTP requests by method, route, and status."}, []string{"method", "route", "status"}),
		HTTPDuration:       prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration.", Buckets: prometheus.DefBuckets}, []string{"method", "route"}),
		Payments:           prometheus.NewCounterVec(prometheus.CounterOpts{Name: "payments_total", Help: "Payment attempts by result."}, []string{"status"}),
		Refunds:            prometheus.NewCounter(prometheus.CounterOpts{Name: "refunds_created_total", Help: "Successful refunds."}),
		DatabaseErrors:     prometheus.NewCounter(prometheus.CounterOpts{Name: "database_errors_total", Help: "Database operation failures."}),
		IdempotencyReplays: prometheus.NewCounter(prometheus.CounterOpts{Name: "idempotency_replays_total", Help: "Completed payment requests replayed."}),
	}
	registerer.MustRegister(m.HTTPRequests, m.HTTPDuration, m.Payments, m.Refunds, m.DatabaseErrors, m.IdempotencyReplays)
	return m
}
