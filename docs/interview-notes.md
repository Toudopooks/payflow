# Interview Notes

## Installment Calculation

Divide the total minor units by the installment count. Assign the base amount to every installment, then add one minor unit to the first `remainder` installments. `10000 / 3` becomes `3334 + 3333 + 3333`.

## Why Integer Money

Binary floating point produces rounding surprises. Integers give exact addition, comparison, database constraints, and refund limits. Formatting into dollars or pesos happens only at the UI boundary.

## Transaction Protection

Payments and refunds use MySQL transactions plus row locks. If an event insert or state update fails, rollback removes every write from that attempt. No installment is left paid without its successful payment record.

## Idempotency

The caller sends a stable key for one logical payment. The server persists it with a request hash and payment ID. A repeated identical request returns the original result; different input under the same key is rejected.

## Why Duplicate Payments Matter

HTTP clients retry after timeouts and connection resets. The first request may have committed even when its response never arrived. Without server-side idempotency, normal retry behavior can duplicate a financial operation.

## Refunds

The payment row is locked, and `amount - refunded_amount` is calculated inside the transaction. Partial refunds update the cumulative amount. A fully refunded payment and installment move to `refunded`; the order follows once no paid installments remain.

## Tests

Unit tests target pure arithmetic. The MySQL suite targets locks, constraints, idempotency, persistence, and rollback-sensitive state. Component tests target rendering and recovery states. Playwright covers the user flow.

## GitLab CI

Lint and tests run before build jobs. The integration job starts a MySQL service, waits for its port, applies migrations from the test, then executes the payment flow. No deployment is claimed.

## Docker

Compose runs MySQL, a one-shot migration container, the Go API, Nginx-hosted SPA, Prometheus, and Grafana. Health and dependency conditions keep startup ordering explicit.

## Kubernetes

The manifests show Deployments for stateless services, a StatefulSet and claim for MySQL, Services, probes, resource requests, configuration, secrets, and optional Ingress routing. They are a fundamentals demonstration, not a production database design.

## Prometheus and Grafana

The API records labeled request counts and duration histograms. Domain counters track payment results, refunds, database errors, and idempotency replays. Grafana queries request rate, error ratio, p95 latency, and cumulative operation counts.

## Tradeoffs

- One service keeps consistency easier to reason about than distributed transactions.
- `READ COMMITTED` plus explicit row locks is targeted and understandable.
- One payment table records failed and successful attempts, preserving history.
- Authentication was deferred to keep attention on financial integrity.

## Three Useful Edge Cases

1. `10000 / 3` must not lose one cent.
2. The same idempotency key with changed input must not silently reuse an unrelated payment.
3. Two partial refunds must never exceed the successful payment amount.

## Real-World Improvements

Add merchant authentication and tenant scoping, processor webhooks, asynchronous settlement states, idempotency cleanup, audit retention policy, tracing, rate limits, fraud checks, reconciliation jobs, and a managed highly available database.

