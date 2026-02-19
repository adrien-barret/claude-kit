---
name: checklist
description: Pre-implementation readiness check — quality gate before coding begins
---

Act as a **QA Lead** performing a pre-implementation readiness check. Verify that all artifacts are complete and aligned before any code is written.

## Prerequisites

All three are required:
- `.claude/output/problem.yaml` — if missing, tell the user to run `/bmad-break` first and stop
- `.claude/output/architecture.yaml` — if missing, tell the user to run `/bmad-model` first and stop
- `.claude/output/backlog.yaml` — if missing, tell the user to run `/bmad-model` first and stop

Optional (used if present):
- `.claude/output/principles.md` — project principles for validation

Read all artifacts before starting the checklist.

## Checklist Categories

### 1. Requirements Completeness

| Check | Evidence |
|-------|----------|
| Every feature has at least one user story | Count stories per feature |
| Every user story has a `so_that` (business value / WHY) | List any missing |
| Every user story has at least one Given/When/Then scenario | List any missing |
| Every user story has a `testability` field | List any missing |
| All P1 features have complete user stories | List incomplete P1s |
| No orphan stories (every story belongs to a feature) | Cross-reference |
| Clarifications section exists if `/clarify` was run | Check for `clarifications:` block |

### 2. Architecture Alignment

| Check | Evidence |
|-------|----------|
| Every feature maps to at least one architecture component | Coverage matrix |
| Tech stack in `problem.yaml` matches `architecture.yaml` | Compare fields |
| All integrations have defined contracts or interface descriptions | List undefined |
| Data model covers all entities implied by user stories | Cross-reference |
| No architecture component is unused (not referenced by any task) | List orphans |

### 3. Security

| Check | Evidence |
|-------|----------|
| Authentication/authorization approach is defined | Reference architecture section |
| Secrets management strategy is specified | Reference constraints or architecture |
| Input validation is addressed in relevant stories | Check acceptance criteria |
| Dependency scanning is planned (if applicable) | Check CI/CD or constraints |
| Compliance requirements are reflected in stories (if any) | Cross-reference constraints |

### 4. Testing Strategy

| Check | Evidence |
|-------|----------|
| Testing approach is defined (unit, integration, e2e) | Reference architecture or principles |
| Every P1 story has testable acceptance criteria | Review Given/When/Then |
| Test infrastructure is specified (framework, CI integration) | Reference tech stack |
| Coverage expectations are defined | Reference principles or constraints |

### 5. Performance

| Check | Evidence |
|-------|----------|
| Performance targets are measurable (latency, throughput) | Reference constraints |
| Performance-critical paths are identified | Cross-reference architecture |
| Caching strategy is defined (if applicable) | Reference architecture |
| Resource budgets exist (if applicable) | Reference constraints |

## Evaluation

For each check item, assign:
- **PASS** — evidence exists and is sufficient
- **FAIL** — missing, incomplete, or contradictory
- **N/A** — not applicable to this project (with brief justification)

## Output

Write `.claude/output/checklist.md`:

```markdown
# Pre-Implementation Checklist

> Generated on {date} for {project_name}

## Summary

| Category | Pass | Fail | N/A | Status |
|----------|------|------|-----|--------|
| Requirements Completeness | X | X | X | PASS/FAIL |
| Architecture Alignment | X | X | X | PASS/FAIL |
| Security | X | X | X | PASS/FAIL |
| Testing Strategy | X | X | X | PASS/FAIL |
| Performance | X | X | X | PASS/FAIL |

**Overall: {PASS / FAIL — N items need attention}**

## Detailed Results

### 1. Requirements Completeness

- [PASS] Every feature has at least one user story (8/8 features covered)
- [FAIL] User story US-004 missing `so_that` — no business value justification
- [PASS] All P1 features have complete stories
...

### 2. Architecture Alignment
...

### 3. Security
...

### 4. Testing Strategy
...

### 5. Performance
...

## Action Items

{Numbered list of FAIL items with suggested fixes, ordered by impact}

1. [FAIL] US-004 missing business value → add `so_that` field in problem.yaml
2. [FAIL] No caching strategy defined → add to architecture.yaml
...
```

## Gate Behavior

When run as part of `/bmad-run`:
- **Any FAIL item blocks progression** — the workflow stops until all FAILs are resolved or explicitly waived by the user
- Present the action items and ask the user to fix or waive each one

When run standalone (`/checklist`):
- All results are reported; the user decides what to act on

If $ARGUMENTS is provided, use it as focus area or additional context: $ARGUMENTS
