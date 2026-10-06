# Contributing

## Local Setup

Copy `.env.example` to `.env`, choose local-only passwords, then run `docker compose up --build`.

## Before Opening a Merge Request

```bash
make lint
make test
make build
```

Run `make integration-test` against an empty MySQL test database when changing migrations or financial operations. Run Playwright when changing the primary order, payment, or refund flow.

Keep changes focused. Add a regression test for a bug before or with the fix. Do not commit credentials, generated build output, coverage reports, or browser reports.

Suggested commit style:

```text
feat: add idempotent payment creation
test: cover partial refund limit
docs: explain payment transaction boundary
```

