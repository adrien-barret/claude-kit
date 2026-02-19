---
name: principles
description: Define project principles — structured debate between Product Owner and Tech Lead to establish governance standards
---

You will role-play two distinct personas — a **Product Owner** (PO) and a **Tech Lead** (TL) — who **challenge each other** to produce a robust set of project principles. These two roles have natural tension: the PO optimizes for user value, speed-to-market, and scope; the TL optimizes for code quality, maintainability, and technical sustainability. Good principles emerge from this friction.

The output is a human-readable governance document (`.claude/output/principles.md`) that guides all future development decisions. It is OPTIONAL — other commands work without it, but when it exists, commands like `/bmad-break`, `/analyze`, and `/ralph` will reference it.

## Stage 1: Gather Existing Context

1. Check if `.claude/output/principles.md` already exists. If so, read it and ask the user if they want to update it or start fresh.
2. Scan the codebase for signals that hint at existing standards:
   - Linting config (`.eslintrc`, `.prettierrc`, `golangci-lint`, `ruff.toml`, `.flake8`, etc.)
   - Test frameworks and config (`jest.config`, `pytest.ini`, `vitest.config`, `go test`, etc.)
   - CI/CD pipelines (`.github/workflows/`, `.gitlab-ci.yml`, `Jenkinsfile`, etc.)
   - Security tools (`snyk`, `trivy`, `dependabot.yml`, `.npmrc`, etc.)
   - Type checking (`tsconfig.json`, `mypy.ini`, `pyright`, etc.)
   - Code formatting (`prettier`, `black`, `gofmt`, `rustfmt`, etc.)
3. Read `.claude/output/problem.yaml` if it exists — extract project context (tech stack, constraints).
4. Read `CLAUDE.md` for any existing code principles.

Summarize what you found to the user before proceeding.

## Stage 2: Structured Debate (PO vs TL)

Run **2-3 rounds** of structured debate across the relevant categories below. In each round:

1. **One role proposes** principles for a category (with rationale)
2. **The other role plays devil's advocate** — challenges assumptions, identifies blind spots, argues trade-offs
3. **They switch roles** for the next category
4. **Convergence**: after the exchange, state where they agree and where they disagree

### Debate format (show this to the user)

For each category, output the debate visibly:

```
📋 [Category Name]

PO proposes: ...
TL challenges: ...
PO responds: ...
→ Agreed: [principle]
→ Conflict: [PO position] vs [TL position]
```

Or when TL proposes first:

```
📋 [Category Name]

TL proposes: ...
PO challenges: ...
TL responds: ...
→ Agreed: [principle]
→ Conflict: [TL position] vs [PO position]
```

### Categories

1. **Code Quality** — style, naming, error handling, review requirements
2. **Testing Standards** — required types, coverage, TDD/BDD, what must always have tests
3. **Security** — auth, secrets, validation, scanning, compliance
4. **UX Principles** (if applicable) — accessibility, performance budgets, responsive, i18n
5. **Performance** (if applicable) — latency targets, throughput, resource budgets, caching
6. **Architecture** — patterns, API standards, data layer, module boundaries

Skip categories that don't apply (e.g., skip UX for a CLI tool, skip Performance for an internal script).

### Depth control

- Default: 2 rounds of debate (propose → challenge → respond → converge)
- If `$ARGUMENTS` contains "advocate": add a 3rd round where unresolved conflicts get a deeper back-and-forth before escalating to the user
- Keep it focused — this is not a continuous loop. 2-3 exchanges per category, then move on.

## Stage 3: User Arbitration

After the debate rounds, present the results to the user:

1. **Agreed principles** — both roles converged, no conflict. List them clearly.
2. **Unresolved conflicts** — present each conflict with both positions and their arguments. Ask the user to pick a side or propose a middle ground.

If there are no conflicts, skip straight to confirmation.

Wait for the user to resolve all conflicts before proceeding.

## Stage 4: Write Principles Document

Create `.claude/output/principles.md`:

```markdown
# Project Principles

> Governing principles for {project_name}. All code, architecture decisions, and reviews
> must align with these standards. Generated on {date}.

## Code Quality

{principles as bullet points, each with a brief rationale}

## Testing Standards

{principles as bullet points}

## Security

{principles as bullet points}

## UX Principles

{principles as bullet points, or omit section if N/A}

## Performance

{principles as bullet points, or omit section if N/A}

## Architecture

{principles as bullet points}

## Exceptions

{any explicitly agreed exceptions or deviations, with rationale}
```

## Stage 5: Final Confirmation

Present the full document to the user. Highlight any areas where a default was inferred rather than explicitly discussed.

Once confirmed, save to `.claude/output/principles.md` and report completion.

If $ARGUMENTS is provided, use it as context or instructions (e.g., "advocate" for deeper debate, or specific areas to focus on): $ARGUMENTS
