---
name: bmad-run
description: Run the full BMAD v6 workflow (Principles → Break → Clarify → Model → Analyze → Checklist → GSD Prep → Act → Deliver) with dev skills, security, and FinOps
---

Orchestrate the complete **BMAD v6 workflow**. Execute each phase in sequence, carrying context forward between phases.

## Workflow

### Phase 0: Principles -- Project Governance (optional)

Follow the instructions in `/principles`:
- PO and TL debate project principles across code quality, testing, security, UX, performance, architecture
- Produce `.claude/output/principles.md`
- Get user confirmation

**Gate**: Ask the user: "Define project principles now, or skip?" If the user says "skip", proceed to Phase 1 without principles. The rest of the workflow adapts — principles-dependent checks are skipped gracefully.

### Phase 1: Break -- Define the Problem

Follow the instructions in `/bmad-break`:
- Gather the project brief from the user
- Clarify requirements through targeted questions
- Break features into rich user stories with Given/When/Then acceptance scenarios and business value (WHY)
- Produce `.claude/output/problem.yaml`
- Get user confirmation before proceeding

**Gate**: Do not proceed to Phase 1.5 until the user confirms the problem definition.

### Phase 1.5: Clarify -- Structured Ambiguity Scan

Follow the instructions in `/clarify`:
- Scan `problem.yaml` across 8 categories (functional scope, data model, UX flow, non-functional, integrations, edge cases, constraints, terminology)
- Ask max 5 targeted questions with recommendations and options
- Update `problem.yaml` with clarifications (inline + audit trail)
- Report summary

**Gate**: Do not proceed to Phase 2 until clarification is complete. If significant gaps remain, warn the user.

### Phase 2: Model -- Design Architecture & Backlog

Follow the instructions in `/bmad-model`:
- Design system architecture based on the confirmed problem definition
- Produce architecture decision records
- Generate a prioritized implementation backlog
- Produce `.claude/output/architecture.yaml` and `.claude/output/backlog.yaml`
- Get user confirmation before proceeding

**Gate**: Do not proceed to Phase 2.25 until the user confirms the architecture and backlog.

### Phase 2.25: Analyze -- Cross-Artifact Consistency

Follow the instructions in `/analyze`:
- Read-only analysis across `problem.yaml`, `architecture.yaml`, `backlog.yaml`, and `principles.md` (if present)
- Detect duplications, ambiguities, underspecification, principles violations, coverage gaps, inconsistencies, terminology drift
- Report findings with severity (CRITICAL/HIGH/MEDIUM/LOW)
- Coverage summary: requirement → architecture → backlog mapping

**Gate**: CRITICAL issues block progression. The user must resolve them (by updating artifacts) before proceeding. After the user resolves CRITICAL issues, re-run `/analyze` to verify the fixes before proceeding to Phase 2.3. HIGH/MEDIUM/LOW issues are reported but don't block.

### Phase 2.3: Checklist -- Pre-Implementation Quality Gate

Follow the instructions in `/checklist`:
- Evaluate 5 categories: requirements completeness, architecture alignment, security, testing strategy, performance
- Each item scored PASS/FAIL/N/A with evidence
- Produce `.claude/output/checklist.md`

**Gate**: Any FAIL item blocks progression. The user must fix or explicitly waive each FAIL before proceeding. If the user fixes (rather than waives) FAIL items, re-run `/checklist` to verify before proceeding.

### Phase 2.5: GSD Prep -- Context Engineering

Follow the instructions in `/gsd-prep`:
- Scan the existing codebase and produce mapping files
- Run gap analysis: compare backlog tasks against existing code
- Atomise tasks: verify each story fits in a sub-agent context
- Generate bounded context packs for each story
- Produce `.claude/output/gsd/prep-report.md`
- Get user confirmation before proceeding

**Gate**: Do not proceed to Phase 3 until the user confirms the GSD prep results.

### Phase 3: Act -- Implement Code (Ralph Agent Team)

Follow the instructions in `/ralph`:
- Ralph uses the context packs from GSD Prep (`.claude/output/gsd/context-packs/`)
- Each teammate receives a bounded context instead of the full architecture dump
- Parse the backlog into a PRD with parallel implementation rounds
- Create numbered feature branch (`NNN-<feature-name>`)
- Create an agent team — spawn teammates per story in each round
- Teammates implement stories in parallel, write tests, commit
- Wait for each round to complete before starting the next
- Run quality checks (code review, tests, security scan)
- Produce `.claude/output/act-report.md`

**Gate**: Do not proceed to Phase 4 until all stories pass and quality checks are complete. Note: Ralph runs its own quality checks (Step 7). Phase 4 and 5 extend those — do not re-run the same checks, build on Ralph's output.

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
- **Business value first**: Every feature must justify its existence (WHY)

## Output Artifacts

At the end of the workflow, the following files will exist in `.claude/output/`:
- `principles.md` -- Project principles (Principles, optional)
- `problem.yaml` -- Problem definition with rich user stories (Break)
- `architecture.yaml` -- Architecture design (Model)
- `backlog.yaml` -- Implementation backlog (Model)
- `checklist.md` -- Pre-implementation readiness check (Checklist)
- `gsd/prep-report.md` -- GSD preparation report (GSD Prep)
- `gsd/gap-analysis.md` -- Gap analysis (GSD Prep)
- `gsd/codebase/` -- Codebase mapping (GSD Prep)
- `gsd/context-packs/` -- Bounded context packs per story (GSD Prep)
- `act-report.md` -- Implementation report (Act)
- `release-notes.md` -- Release notes (Deliver)

Plus the actual project source code, tests, infrastructure, and documentation in the project tree.

If $ARGUMENTS is provided, use it as the initial project brief: $ARGUMENTS
