# Architecture

## System Overview

```mermaid
flowchart TB
  subgraph Client
    SPA[React + TypeScript SPA]
  end
  subgraph Application
    HTTP[Chi HTTP transport]
    Service[Domain service]
    Store[MySQL repository]
  end
  subgraph Data
    MySQL[(MySQL 8)]
  end
  subgraph Operations
    Metrics[/Prometheus metrics/]
    Dashboard[Grafana]
  end
  SPA -->|REST| HTTP
  HTTP --> Service --> Store --> MySQL
  HTTP --> Metrics --> Dashboard
```

The HTTP layer parses input and maps domain errors to stable API responses. The service owns validation, installment construction, and command creation. The repository owns SQL, row locks, and transaction boundaries. This split keeps business rules testable without hiding SQL behind unnecessary generic abstractions.

## Payment Sequence

```mermaid
sequenceDiagram
  participant C as Client
  participant A as API
  participant D as MySQL
  C->>A: POST /payments + Idempotency-Key
  A->>D: BEGIN
  A->>D: insert or lock idempotency key
  alt completed key
    D-->>A: original payment ID
    A->>D: COMMIT
    A-->>C: 200 original payment
  else new key
    A->>D: lock installment
    A->>D: insert payment and event
    A->>D: update installment and plan/order if complete
    A->>D: complete idempotency record
    A->>D: COMMIT
    A-->>C: 201 payment
  end
```

The idempotency row and payment mutation share a transaction. `SELECT ... FOR UPDATE` serializes competing operations for the same installment. Any error rolls back the full mutation.

## Frontend

React Router maps the dashboard, order list, order creation, order details, and transaction lookup. TanStack Query owns server state and cache invalidation. React Hook Form with Zod owns input parsing and inline errors. Search state for transaction history lives in the URL.

## CI/CD

GitLab CI has lint, test, integration, and build stages. The integration job starts MySQL 8.4 and runs the same migration and payment flow used locally. Build artifacts are produced, but deployment is intentionally manual.

## Deployment

Docker Compose is the primary local environment. Kubernetes manifests demonstrate the mechanics of a cluster deployment. They are intentionally plain YAML to keep the operational model visible.

