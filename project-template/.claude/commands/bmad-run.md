---
name: bmad-run
description: Run the full BMAD v6 workflow (Break → Model → Act → Deliver) with dev skills, security, and FinOps
---

Orchestrate the complete **BMAD v6 workflow**. Execute each phase in sequence, carrying context forward between phases.

## Workflow

### Phase 1: Break -- Define the Problem

Follow the instructions in `/bmad-break`:
- Gather the project brief from the user
- Clarify requirements through targeted questions
- Produce `.claude/output/problem.yaml`
- Get user confirmation before proceeding

**Gate**: Do not proceed to Phase 2 until the user confirms the problem definition.

### Phase 2: Model -- Design Architecture & Backlog

Follow the instructions in `/bmad-model`:
- Design system architecture based on the confirmed problem definition
- Produce architecture decision records
- Generate a prioritized implementation backlog
- Produce `.claude/output/architecture.yaml` and `.claude/output/backlog.yaml`
- Get user confirmation before proceeding

**Gate**: Do not proceed to Phase 3 until the user confirms the architecture and backlog.

### Phase 3: Act -- Implement Code (Ralph Agent Team)

Follow the instructions in `/ralph`:
- Parse the backlog into a PRD with parallel implementation rounds
- Create an agent team — spawn teammates per story in each round
- Teammates implement stories in parallel, write tests, commit
- Wait for each round to complete before starting the next
- Run quality checks (code review, tests, security scan)
- Produce `.claude/output/act-report.md`

**Gate**: Do not proceed to Phase 4 until all stories pass and quality checks are complete.

### Phase 4: Dev Skills

Run these checks on the implemented code:
- **Code review**: Review all produced code for quality, principle adherence, and bugs
- **Test coverage**: Verify test coverage is adequate; generate additional tests if needed
- **API documentation**: Generate or update API docs if the project has an API
- **Dependency audit**: Check for vulnerable or outdated dependencies

### Phase 5: Security & FinOps

Run security and cost checks:
- **Security audit**: Check for code vulnerabilities, infra misconfigurations, auth weaknesses, exposed secrets
- **FinOps review** (if infrastructure code exists): Check tagging, rightsizing, waste, cost optimization

Report findings and apply fixes for critical issues. Present non-critical findings to the user.

### Phase 6: Deliver -- Prepare Release

Follow the instructions in `/bmad-deliver`:
- Create deployment configuration
- Update documentation
- Run final security review
- Produce `.claude/output/release-notes.md`
- Present release checklist to the user

## Principles

Throughout the entire workflow, enforce:
- **DRY, KISS, SOLID**: In all code and infrastructure
- **Least invasive**: Minimal changes, no unnecessary refactoring
- **No over-engineering**: Only what the requirements call for
- **User confirmation at each gate**: Never skip a gate

## Output Artifacts

At the end of the workflow, the following files will exist in `.claude/output/`:
- `problem.yaml` -- Problem definition (Break)
- `architecture.yaml` -- Architecture design (Model)
- `backlog.yaml` -- Implementation backlog (Model)
- `act-report.md` -- Implementation report (Act)
- `release-notes.md` -- Release notes (Deliver)

Plus the actual project source code, tests, infrastructure, and documentation in the project tree.

If $ARGUMENTS is provided, use it as the initial project brief: $ARGUMENTS
