---
name: code-only
description: Generate only the project code, skipping tests, docs, and security checks
---

Act as a **development team** composed of the agents defined in `.claude/agents/` (backend, tech-lead).

Your goal is to implement code from the backlog as fast as possible, skipping tests, documentation, and security checks.

## Prerequisites

Read `.claude/output/backlog.yaml` and `.claude/output/architecture.yaml`. If either does not exist, tell the user to run `/bmad-model` first and stop.

## What to do

Follow the same implementation process as `/bmad-act` but with these differences:

- **Skip** test generation (Stage 2 step 4)
- **Skip** quality checks (Stage 3)
- **Skip** the implementation report (Stage 4)

Focus purely on writing the source code:

1. Read the backlog and sort by dependency order.
2. For each task, implement the code following the architecture and acceptance criteria.
3. Write code directly in the project source tree.
4. Follow all code principles (DRY, KISS, SOLID, least invasive, no over-engineering).

When done, list all files created and modified.

If $ARGUMENTS is provided, use it as a task filter or additional context: $ARGUMENTS
