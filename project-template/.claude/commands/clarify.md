---
name: clarify
description: Structured ambiguity scan — identify and resolve gaps in the problem definition through targeted questions
---

Act as a **Requirements Analyst** performing a structured ambiguity scan on the problem definition.

Your goal is to find gaps, ambiguities, and underspecified areas in `problem.yaml`, then resolve them through a focused sequence of questions.

## Prerequisites

Read `.claude/output/problem.yaml`.

- **If it exists**: use it as the basis for the ambiguity scan.
- **If it does not exist AND `$ARGUMENTS` contains a project description**: treat the argument as an inline problem statement. Run the ambiguity scan against it, and write the clarified result to `.claude/output/problem.yaml` at the end (effectively combining break + clarify in one pass).
- **If it does not exist AND no `$ARGUMENTS`**: tell the user to run `/bmad-break` first and stop.

Also read `.claude/output/principles.md` if it exists — use it to calibrate expectations (e.g., if principles mandate 95% test coverage, verify that testability fields are sufficient).

## Stage 1: Ambiguity Scan

Systematically scan `problem.yaml` across this taxonomy:

1. **Functional Scope & Business Value** — Are feature boundaries clear? Are there features that overlap or have undefined edges? Are all user stories complete (actor, action, value)? **Critically: does every story have a clear "so_that" (WHY)?** A story without business value justification is a story that shouldn't exist. Challenge any feature where the "why" is vague, circular, or missing.
2. **Data Model** — What data entities are implied but not defined? What are the relationships? Are there CRUD operations without a clear schema?
3. **UX Flow** — Are user journeys specified end-to-end? Are error states, empty states, and edge cases covered? Are there dead-end flows?
4. **Non-Functional Requirements** — Are performance targets measurable? Are availability SLAs defined? Is observability scoped?
5. **Integrations** — Are API contracts specified? Are authentication flows for third-party services clear? What happens when an integration is down?
6. **Edge Cases** — What happens with empty inputs, concurrent access, rate limits, large payloads, timezone differences, unicode?
7. **Constraints** — Are budget, timeline, compliance, and deployment constraints explicit? Are there implicit constraints (e.g., must work offline)?
8. **Terminology** — Are domain terms used consistently? Could any term mean different things to different stakeholders?

For each category, rate it:
- **Clear** — no ambiguity found
- **Minor gap** — small clarification needed, won't block implementation
- **Significant gap** — could lead to wrong implementation if not clarified

## Stage 2: Targeted Questions

Ask **max 5 questions**, one at a time, sequentially. Prioritize by impact — ask about significant gaps first.

Each question must include:
1. **Context** — why this matters (what could go wrong if not clarified)
2. **Your recommendation** — what you think the answer should be, based on the existing context
3. **Options** — 2-4 multiple-choice options (including your recommendation marked as such) plus "other"

Example format:

```
Question 1/5: [Category — Functional Scope]

The "user authentication" feature lists JWT but doesn't specify token refresh strategy.
Without this, teammates may implement incompatible approaches.

My recommendation: Short-lived access tokens (15min) + long-lived refresh tokens (7d)

Options:
  a) Short-lived access + refresh tokens (recommended)
  b) Long-lived access tokens only (simpler but less secure)
  c) Session-based auth instead of JWT
  d) Other — please specify
```

Wait for the user's answer before asking the next question. If the user's answer resolves the gap clearly, move on. If it raises a new question, it counts toward the 5-question budget.

If all significant gaps are resolved before reaching 5 questions, stop early — don't ask questions for the sake of filling the quota.

## Stage 3: Update Problem Definition

Write clarifications back to `problem.yaml`:

1. **Inline updates** — where the clarification fills an existing gap (e.g., adding a missing `testability` field to a user story, specifying a performance target), update the value in place.
2. **Clarifications section** — add a `clarifications:` block at the end of `problem.yaml` capturing the Q&A trail:

```yaml
clarifications:
  - question: <the question asked>
    category: <taxonomy category>
    answer: <user's answer>
    impact: <what was updated in the spec as a result>
```

This creates an audit trail — future commands can see what was clarified and why.

## Stage 4: Report

Present a summary to the user:

```
Clarification Summary
─────────────────────
Scanned: 8 categories
Clear: 5  |  Minor gaps: 2  |  Significant gaps: 1
Questions asked: 3/5
Updates applied: 4 inline changes + 3 clarification records

Remaining minor gaps (won't block implementation):
  - [gap description]
  - [gap description]
```

If any significant gaps remain unresolved (user chose to skip), flag them clearly — they may surface as issues during `/analyze`.

If $ARGUMENTS is provided, use it as focus area or additional context: $ARGUMENTS
