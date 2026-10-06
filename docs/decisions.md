# Engineering Decisions

## Go and Chi

Go makes concurrency and error paths explicit while producing a small deployment artifact. Chi adds routing and middleware composition without imposing a larger application framework.

## MySQL

The domain is relational and needs foreign keys, row locks, unique constraints, and transactions. MySQL 8 provides those properties and is common in web platforms.

## Integer Minor Units

Binary floating point cannot exactly represent many decimal values. Integer minor units make equality and sum checks deterministic. Currency rules remain explicit at formatting boundaries.

## Idempotency Keys

Clients retry when connections fail, even when the server may have committed. A persisted key lets the retry recover the original payment instead of charging twice. The request hash prevents accidental reuse for a different command.

## Transaction Boundaries

A payment is not only a row insert. Payment, installment, order progress, event history, and idempotency state must agree. Keeping these writes in one transaction provides an all-or-nothing outcome.

## REST

The resources and commands map cleanly to HTTP and JSON. REST keeps the project approachable and makes retries, status codes, and idempotency behavior visible.

## Repository and Service Separation

The service expresses application rules. The MySQL repository expresses locking and persistence. There is one concrete repository interface, not a hierarchy of generic data-access layers.

## Prometheus and Grafana

Counters and latency histograms make failure and retry behavior visible. Prometheus fits pull-based HTTP metrics; Grafana keeps the useful operational queries version controlled.

## GitLab CI/CD

The pipeline demonstrates a service-backed integration job and independent frontend and backend builds. GitHub can remain the public home while a GitLab mirror runs CI.

## No Card Data

Real payment credentials would add security and compliance requirements unrelated to this project's learning goals. All outcomes are simulated and demo failure is explicitly gated by configuration.

