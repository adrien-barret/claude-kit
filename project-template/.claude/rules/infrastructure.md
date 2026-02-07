---
description: Infrastructure rules applied when working on IaC, Docker, CI/CD, and cloud configs
globs: ["infra/**", "terraform/**", "*.tf", "Dockerfile*", "docker-compose*", ".github/workflows/**", "k8s/**", "helm/**"]
---

> See `rules/security.md` for credential handling. See `rules/finops.md` for resource tagging.

## Infrastructure as Code

- Use variables and modules to avoid hardcoding values; every environment-specific value should be a variable
- Pin provider and module versions explicitly to avoid breaking changes
- Configure remote state backend with locking (S3+DynamoDB, GCS, Terraform Cloud)
- Run `terraform plan` in CI and require approval before `terraform apply`
- Use infrastructure testing tools (Terratest, Conftest, OPA) to validate policies before merge

## Network and IAM

- Use least-privilege IAM policies; avoid wildcard (`*`) permissions
- Enable encryption at rest and in transit for all data stores
- Segment networks using VPCs and subnets: separate public-facing, application, and data tiers
- Restrict security groups to specific CIDR ranges and ports; never allow `0.0.0.0/0` to management ports

## Containers

- Use multi-stage Docker builds to minimize image size
- Run containers as non-root users; set `USER` directive in Dockerfile
- Scan container images for vulnerabilities in CI (Trivy, Grype) before pushing to registry
- Pin base image versions to a specific digest or tag — do not use `latest`

## Drift and Compliance

- Enable drift detection to catch out-of-band changes
- Treat all infrastructure changes as code: manual console changes must be imported or reverted
