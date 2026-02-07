---
name: bmad-act
description: BMAD Act phase – implement code from the backlog using Ralph agent teams
---

Act as the **implementation coordinator** for the BMAD Act phase.

Your goal is to implement the code from the backlog produced in the Model phase, using `/ralph` to orchestrate an agent team.

## Prerequisites

Read `.claude/output/backlog.yaml` and `.claude/output/architecture.yaml`. If either does not exist, tell the user to run `/bmad-model` first and stop.

## Code Principles (mandatory)

All code produced MUST follow:
- **DRY**: extract shared logic into reusable functions or modules
- **KISS**: simplest approach that works; no premature abstraction
- **SOLID**: single responsibility, open/closed, dependency inversion
- **Least invasive**: change only what the task requires
- **No over-engineering**: no features or abstractions beyond what is asked
- **Separation of concerns**: distinct layers for business logic, data access, transport, infrastructure
- **Clean code**: descriptive naming, small functions, no dead code

## Execution

Follow the `/ralph` command instructions:
1. Parse `.claude/output/backlog.yaml` into a PRD
2. Group stories into parallel rounds based on dependencies
3. Present execution plan and get user confirmation
4. Create an agent team — spawn teammates per round for parallel implementation
5. Each teammate implements one story, writes tests, commits
6. Wait for each round to complete before starting the next
7. Run quality checks after all stories pass
8. Save report to `.claude/output/act-report.md`

Ralph handles the agent team coordination, task assignment, file conflict detection, and progress tracking. See `/ralph` for full details.

If $ARGUMENTS is provided, use it as additional context or task filter: $ARGUMENTS
