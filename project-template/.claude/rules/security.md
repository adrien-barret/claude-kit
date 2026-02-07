---
description: Security rules applied across all code
---

## Secrets and Credentials

- Never hardcode secrets, API keys, tokens, or passwords in source code or config files
- Use environment variables or a secret manager (Vault, AWS Secrets Manager, GCP Secret Manager)
- Set token expiration and implement refresh strategies; short-lived tokens reduce blast radius

## Input Validation and Injection Prevention

- Validate and sanitize all user input at system boundaries — reject unexpected types, lengths, and characters
- Use parameterized queries to prevent SQL injection; never concatenate user input into queries
- Apply the principle of least privilege for all roles, permissions, and service accounts

## Transport and Headers

- Enable HTTPS/TLS for all external communication; redirect HTTP to HTTPS
- Set security response headers on all responses:
  - `Content-Security-Policy` — restrict resource origins
  - `Strict-Transport-Security` — enforce HTTPS (`max-age=31536000; includeSubDomains`)
  - `X-Content-Type-Options: nosniff` — prevent MIME sniffing
  - `X-Frame-Options: DENY` — prevent clickjacking
- Configure CORS to list explicit allowed origins; never use wildcard (`*`) in production

## Dependency and Container Security

- Scan dependencies for known vulnerabilities in CI (Dependabot, Snyk, `npm audit`, `pip-audit`)
- Keep dependencies updated; review and merge security patches within 48 hours
- Do not run containers as root; use minimal base images — see `rules/infrastructure.md`

## Audit Logging

- Log security-relevant events: authentication attempts, authorization failures, privilege changes
- Include: timestamp, user/principal, action, resource, result (success/failure)
- Never log secrets, tokens, passwords, or PII in plaintext
