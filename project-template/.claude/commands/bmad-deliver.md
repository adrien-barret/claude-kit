---
name: bmad-deliver
description: BMAD Deliver phase – prepare release with deployment scripts, docs, and release notes
---

Act as a **DevOps Engineer** and **Tech Lead** working together, using the agents defined in `.claude/agents/devops.md` and `.claude/agents/tech-lead.md`.

Your goal is to prepare the project for release: deployment configuration, documentation, and release notes.

## Prerequisites

1. Read `.claude/output/architecture.yaml` for infrastructure and deployment design.
2. Read `.claude/output/act-report.md` for what was implemented.
3. Read `.claude/output/principles.md` if it exists — use security and deployment principles to guide the release process.
4. Scan the project source tree to understand the current state of the code.

If no code exists in the project, tell the user to run `/bmad-act` first and stop.

## Stage 1: Deployment Configuration

Based on the architecture, create or verify:

1. **Containerization** (if applicable):
   - Dockerfile with multi-stage build, non-root user, pinned base image
   - .dockerignore excluding unnecessary files

2. **Infrastructure as Code** (if applicable):
   - Terraform, CloudFormation, or Pulumi configs for cloud resources
   - Proper variable extraction (no hardcoded values)
   - Least-privilege IAM roles and policies
   - Resource tagging (environment, team, service, cost-center)

3. **CI/CD Pipeline**:
   - Build, test, lint, security-scan stages
   - Environment-specific deployment (dev, staging, production)
   - Secrets injected from environment/secret manager (never hardcoded)

4. **Configuration Management**:
   - Environment-based configuration (env vars or config files)
   - Separation between secrets and non-sensitive config
   - Health check endpoints

## Stage 2: Documentation

Create or update:

1. **README.md**: Project overview, setup instructions, development guide, deployment guide.
2. **API documentation**: If an API exists, generate or update OpenAPI/Swagger docs.
3. **Architecture diagram description**: Text-based description of the system architecture.
4. **Runbook**: Basic operational procedures (deploy, rollback, monitor, troubleshoot).

## Stage 3: Security Review

Run a final security check:

1. No hardcoded secrets in source code or config files.
2. Dependencies are pinned and free of known critical vulnerabilities.
3. Authentication and authorization are properly configured.
4. HTTPS/TLS is enforced for external communication.
5. Logging does not expose sensitive data.

## Stage 4: Release Notes

Create `.claude/output/release-notes.md`:

```markdown
# Release Notes - <project name> v<version>

## Summary
<one-paragraph description of what this release includes>

## Features
- <feature 1>: <description>
- <feature 2>: <description>

## Architecture
- <key architectural decisions>

## Infrastructure
- **Compute**: <what's used>
- **Database**: <what's used>
- **CI/CD**: <pipeline summary>

## Security
- <authentication method>
- <key security measures>

## Known Limitations
- <limitation 1>

## Deployment
1. <step 1>
2. <step 2>
3. <step 3>

## Configuration
| Variable | Description | Required |
|----------|-------------|----------|
| <VAR> | <description> | yes/no |
```

## Stage 5: Validate

Present the release checklist to the user:

- [ ] Deployment config is complete and tested
- [ ] Documentation is up to date
- [ ] Security review passed
- [ ] Release notes are accurate
- [ ] All acceptance criteria from the backlog are met

Ask the user to confirm the release is ready.

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
