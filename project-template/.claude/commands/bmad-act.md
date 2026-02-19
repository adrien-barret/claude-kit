---
name: bmad-act
description: BMAD Act phase – implement code from the backlog using Ralph agent teams
---

Act as the **implementation coordinator** for the BMAD Act phase.

This is Phase 3 of the BMAD workflow — all prior phases (Break, Clarify, Model, Analyze, Checklist, GSD Prep) should be complete before this runs.

## Prerequisites

1. Read `.claude/output/backlog.yaml` and `.claude/output/architecture.yaml`. If either does not exist, tell the user to run `/bmad-model` first and stop.
2. Read `.claude/output/principles.md` if it exists — pass project-specific standards to Ralph for teammate spawn prompts.
3. Read `.claude/output/checklist.md` if it exists — confirm there are no unresolved FAIL items. If there are, warn the user before proceeding.
4. Check if `.claude/output/gsd/prep-report.md` exists. If not, run `/gsd-prep` first to generate codebase mapping and context packs.

## Execution

Delegate to `/ralph` for the full implementation process:

- **Contract-first development**: shared interfaces committed before teammates start
- **Plan approval**: each teammate plans before coding, lead reviews
- **Parallel implementation**: teammates code in parallel within each round
- **Acceptance validation**: lead validates each story, feedback loop on failure
- **Quality checks**: code review, tests, security scan, dependency audit

Follow ALL of Ralph's steps as defined in `/ralph`. Do NOT skip any phase.

## BMAD Gate

This phase is complete when:
- All stories have `passes: true`
- The full test suite passes
- Quality checks report no critical issues
- `.claude/output/act-report.md` has been produced

If $ARGUMENTS is provided, use it as additional context or task filter: $ARGUMENTS
