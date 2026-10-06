# Security Policy

## Scope

PayFlow is an educational payment simulation. It has no authentication or tenant isolation and must not be deployed on a public network as-is. Never send real card numbers, bank credentials, identity documents, or personal customer data to it.

## Reporting

For a public fork, enable GitHub private vulnerability reporting and replace this section with the maintainer's security contact. Do not include exploit details in a public issue before a fix is available.

## Implemented Safeguards

- Parameterized SQL and relational constraints
- Input validation and request body size limits
- Stable client errors without database details or stack traces
- Configured CORS and basic browser security headers
- Environment-based secrets and an ignored `.env`
- Transactional payment and refund writes

These controls do not represent PCI DSS compliance.

