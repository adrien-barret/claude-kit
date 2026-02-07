---
name: security-check
description: Run a comprehensive security audit covering code, infrastructure, auth, and secrets
---

You are a security audit orchestrator. Invoke the security agent from `.claude/agents/security.md` and run all four security sub-skills in sequence.

If $ARGUMENTS is provided, scope the audit to specific directories or areas. Otherwise, audit the entire project.

## Workflow

### 1. Code Security Audit — `.claude/skills/security/code-security-audit/SKILL.md`
Scan for OWASP Top 10, hardcoded secrets, and unsafe crypto.

### 2. Infrastructure Security Audit — `.claude/skills/security/infra-security-audit/SKILL.md`
Review IaC for overly permissive IAM, open security groups, missing encryption.

### 3. Auth Review — `.claude/skills/security/auth-review/SKILL.md`
Review authentication flows, session management, RBAC implementation.

### 4. Secret Rotation — `.claude/skills/security/secret-rotation/SKILL.md`
Identify secrets in code, verify secret manager integration and rotation policies.

## Output Format

```
## Security Audit Report

### Overall Severity: [Critical / High / Medium / Low]

| Area                    | Critical | High | Medium | Low | Total |
|-------------------------|----------|------|--------|-----|-------|
| Code Security           | N        | N    | N      | N   | N     |
| Infrastructure Security | N        | N    | N      | N   | N     |
| Auth Review             | N        | N    | N      | N   | N     |
| Secret Rotation         | N        | N    | N      | N   | N     |

### Findings

#### [SEVERITY] Finding Title
- **Area**: code-security | infra-security | auth | secrets
- **Location**: file:line or resource
- **Description**: What the issue is
- **Remediation**: How to fix it
```
