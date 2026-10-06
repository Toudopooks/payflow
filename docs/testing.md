# Testing Strategy

Tests focus on financial and interaction risk instead of raw test count.

## Backend Unit Tests

`internal/service/installments_test.go` proves deterministic remainder distribution and rejection of unsupported plan sizes. Service-level validation covers invalid money, currencies, and public input through the same functions used by HTTP handlers.

## MySQL Integration Test

`backend/tests/integration/payment_flow_test.go` requires a real empty MySQL database. It applies the migration, creates the full aggregate, pays an installment, replays the same idempotency key, rejects a second payment, creates a partial refund, rejects an over-refund, and checks the event history.

It skips when `TEST_DATABASE_DSN` is absent. GitLab CI supplies the variable and a MySQL 8.4 service.

## Frontend Tests

Vitest and Testing Library cover order rendering, localized money, loading announcements, error recovery, and form schema validation. Playwright covers the cross-layer happy path against a running stack on desktop and mobile viewports.

## Commands

```bash
make test
make integration-test
cd frontend && npm run e2e
```

The browser test is separate from `make test` because it needs the API, MySQL, and browser binaries running.

