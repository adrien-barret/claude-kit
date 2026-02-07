---
description: API design rules applied when building or modifying API endpoints
globs: ["src/routes/**", "src/api/**", "src/controllers/**", "src/handlers/**", "api/**", "routes/**"]
---

## RESTful Design

- Use proper HTTP methods (`GET` read, `POST` create, `PUT` replace, `PATCH` update, `DELETE` remove) and meaningful status codes (`201` for creation, `204` for no content, `404` for not found, `409` for conflict)
- Name resources as plural nouns (`/users`, `/orders`); nest sub-resources one level deep (`/users/{id}/orders`)
- Version APIs in the URL path (`/v1/users`) when introducing breaking changes

## Input and Output

- Validate all request inputs at the controller/handler level — reject invalid requests early with `400`
- Return a consistent error format: `{"error": {"code": "NOT_FOUND", "message": "...", "details": []}}`
- Use pagination for list endpoints (`limit`/`offset` or cursor-based); never return unbounded collections
- Support content negotiation via `Accept` headers; default to `application/json`

## Reliability and Security

- Make `POST` and `PUT` idempotent where possible using client-provided `Idempotency-Key` headers
- Apply rate limiting and authentication middleware on all public endpoints
- Configure CORS explicitly: list allowed origins, methods, and headers — never use wildcard (`*`) in production

## Documentation and Observability

- Document every endpoint with OpenAPI/Swagger annotations including request/response schemas and examples
- Expose a health check endpoint (`GET /healthz`) that returns service status
- Log request method, path, status code, and duration for every request; never log credentials or PII
