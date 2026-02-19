---
name: analyze
description: Cross-artifact consistency analysis — read-only check for gaps, conflicts, and drift across all BMAD outputs
---

Act as a **QA Analyst** performing a cross-artifact consistency check. This command is **read-only** — it modifies nothing, only reports findings.

## Prerequisites

All three are required:
- `.claude/output/problem.yaml` — if missing, tell the user to run `/bmad-break` first and stop
- `.claude/output/architecture.yaml` — if missing, tell the user to run `/bmad-model` first and stop
- `.claude/output/backlog.yaml` — if missing, tell the user to run `/bmad-model` first and stop

Optional (used if present):
- `.claude/output/principles.md` — project principles for violation checks

Read all artifacts before starting the analysis.

## Analysis Categories

### 1. Duplications
Detect requirements, features, or tasks that appear in multiple places with different wording but the same intent. Flag cases where:
- Two backlog tasks cover the same functionality
- A feature in `problem.yaml` maps to multiple overlapping tasks
- Architecture components duplicate responsibilities

### 2. Ambiguities
Find requirements or tasks that are vague enough to produce different implementations depending on interpretation:
- User stories missing acceptance scenarios
- Tasks without clear acceptance criteria
- Architecture decisions that reference undefined components
- "TBD", "TODO", or placeholder values

### 3. Underspecification
Identify areas where critical detail is missing:
- Features without user stories
- User stories without `so_that` (missing business value / WHY)
- Tasks without `depends_on` that clearly need prior work
- Architecture components without defined interfaces
- Integrations without error handling strategy

### 4. Principles Violations (if `principles.md` exists)
Cross-reference artifacts against the project principles:
- Testing standards not reflected in acceptance criteria
- Security principles not addressed in relevant features
- Architecture principles contradicted by design decisions
- Performance targets without corresponding non-functional requirements

### 5. Coverage Gaps
Check that the full chain is complete: every requirement maps to architecture, every architecture component maps to backlog tasks:
- Requirements in `problem.yaml` not covered by any backlog task
- Architecture components not exercised by any task
- Backlog tasks that don't trace back to any requirement (orphan tasks)

### 6. Inconsistencies
Detect contradictions between artifacts:
- `problem.yaml` says REST but `architecture.yaml` defines GraphQL
- Priority P1 in problem but task is marked low priority in backlog
- Tech stack mismatch between problem definition and architecture
- Conflicting non-functional requirements (e.g., "real-time" + "batch processing" for same data)

### 7. Terminology Drift
Find cases where the same concept uses different names across artifacts:
- "user" vs "customer" vs "account" for the same entity
- "order" vs "purchase" vs "transaction" for the same flow
- Component names that don't match between architecture and backlog

## Severity Levels

- **CRITICAL** — Will cause implementation failure or major rework. Blocks `/bmad-run` progression.
- **HIGH** — Likely to cause bugs or incorrect implementation. Should be fixed before coding.
- **MEDIUM** — May cause confusion or suboptimal implementation. Fix recommended.
- **LOW** — Minor inconsistency or style issue. Fix at convenience.

## Output Format

Report to stdout (do NOT write files):

```
Cross-Artifact Analysis Report
═══════════════════════════════

Artifacts analyzed:
  - problem.yaml (version X, N features, N user stories)
  - architecture.yaml (N components, N decisions)
  - backlog.yaml (N tasks across N rounds)
  - principles.md (present/absent)

Findings: N total (N critical, N high, N medium, N low)

─── CRITICAL ───────────────────────────────────────

[C-001] Category: Coverage Gap
  Requirement "payment processing" (problem.yaml, feature #3) has no
  corresponding backlog task. This feature will not be implemented.
  → Fix: Add tasks for payment processing to the backlog.

─── HIGH ───────────────────────────────────────────

[H-001] Category: Inconsistency
  problem.yaml specifies PostgreSQL but architecture.yaml references
  MongoDB in the data layer.
  → Fix: Align database choice across artifacts.

─── MEDIUM ─────────────────────────────────────────

[M-001] Category: Underspecification
  User story US-003 has no acceptance scenarios (Given/When/Then).
  → Fix: Add acceptance scenarios to problem.yaml.

─── LOW ────────────────────────────────────────────

[L-001] Category: Terminology Drift
  "user" in problem.yaml, "account" in architecture.yaml,
  "customer" in backlog.yaml — all refer to the same entity.
  → Fix: Standardize on one term.

─── Coverage Summary ───────────────────────────────

| Requirement          | Architecture Component | Backlog Tasks | Status   |
|----------------------|------------------------|---------------|----------|
| User auth            | auth-service           | T-001, T-002  | Covered  |
| Payment processing   | —                      | —             | MISSING  |
| Notifications        | notification-service   | T-005         | Covered  |

─── Verdict ────────────────────────────────────────

{PASS — no critical issues, safe to proceed}
{BLOCK — N critical issues must be resolved before implementation}
```

## Gate Behavior

When run as part of `/bmad-run`:
- **CRITICAL issues block progression** — the workflow stops until they are resolved
- HIGH/MEDIUM/LOW issues are reported but don't block

When run standalone (`/analyze`):
- All findings are reported; the user decides what to act on

If $ARGUMENTS is provided, use it as focus area or additional context: $ARGUMENTS
