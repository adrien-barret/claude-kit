---
name: gsd-prep
description: GSD Prep — codebase mapping, gap analysis, task atomisation, and bounded context packs for Ralph teammates
---

You are a **context engineer**. Your job is to prepare bounded, focused context packs so that each Ralph teammate receives only the information relevant to their story — eliminating context rot as projects grow.

## Prerequisites

Read `.claude/output/backlog.yaml` and `.claude/output/architecture.yaml`. If either does not exist, tell the user to run `/bmad-model` first and stop.

## Stage 1: Codebase Scan

Scan the existing project codebase and produce mapping files in `.claude/output/gsd/codebase/`.

**For existing projects (code already exists):**

1. **STACK.md** — Languages, frameworks, package managers, key dependencies detected. Scan `package.json`, `go.mod`, `requirements.txt`, `Gemfile`, `Cargo.toml`, `pubspec.yaml`, or equivalent. List versions where available.

2. **ARCHITECTURE.md** — Patterns and layers discovered from the actual code (not from the spec). Identify:
   - Architectural pattern (MVC, hexagonal, layered, microservices, monolith, etc.)
   - Key layers and their responsibilities
   - Data flow between components
   - Entry points (API routes, CLI commands, event handlers)
   - External integrations (databases, APIs, message queues)

3. **STRUCTURE.md** — Directory layout with annotations. List key directories and files with a one-line description of their purpose. Focus on files that stories will interact with.

4. **CONVENTIONS.md** — Coding standards and patterns found in the existing code:
   - Naming conventions (files, functions, variables, types)
   - Error handling patterns
   - Testing patterns and framework
   - Import ordering
   - Configuration approach
   - Logging patterns

5. **CONCERNS.md** — Tech debt, known issues, and risks discovered during the scan:
   - Outdated dependencies
   - Missing tests or low coverage areas
   - Inconsistent patterns
   - Security concerns
   - Performance bottlenecks visible from code structure

**For greenfield projects (no code yet):**

Produce minimal stubs for each file noting this is a new project. STACK.md should reflect the tech stack from `architecture.yaml`. ARCHITECTURE.md should note the planned architecture. STRUCTURE.md, CONVENTIONS.md, and CONCERNS.md should be near-empty with a note that they will be populated as code is written.

## Stage 2: Gap Analysis

Read `.claude/output/backlog.yaml` and `.claude/output/architecture.yaml`. Compare each backlog task against the codebase mapping from Stage 1.

Produce `.claude/output/gsd/gap-analysis.md` with:

### Already Implemented
Tasks where code already exists and acceptance criteria are met or nearly met. These can be dropped or reduced to verification-only stories.

### Partially Implemented
Tasks where some implementation exists but acceptance criteria are not fully met. Note what exists and what remains.

### Purely New
Tasks that require full implementation — no existing code covers them.

### Hidden Dependencies
Dependencies not captured in the backlog's `depends_on` fields:
- Shared utilities or types that multiple stories need but no story creates
- Database migrations that must run before feature code
- Configuration or environment setup not listed as a task
- Third-party service setup or API key provisioning

Present findings to the user. Ask for confirmation before proceeding. If the user wants to adjust the backlog (drop, reduce, or add tasks), make those changes and update the gap analysis accordingly.

## Stage 3: Task Atomisation

For each story in the backlog, estimate the context size it will need when assigned to a teammate:

- Story acceptance criteria
- Relevant architecture subset
- Relevant codebase mapping (files, patterns, conventions)
- Contract definitions
- File ownership list

The total context pack for a single story must fit within ~100k tokens (50% of a 200k context window), leaving room for the teammate's own reasoning, code generation, and tool usage.

**If a story is too large** (estimated context pack > 100k tokens):
- Split it into sub-stories (max 3 sub-stories per original story)
- Each sub-story must be self-contained with its own acceptance criteria
- Update `depends_on` so sub-stories execute in the correct order
- Preserve the original story ID as a prefix (e.g., T-005 splits into T-005a, T-005b, T-005c)

**If a story is small enough**, leave it as-is.

Present any proposed splits to the user for confirmation. Write confirmed changes back to the backlog representation used by Ralph (these are suggestions — the user decides).

## Stage 4: Context Pack Generation

For each story, grouped by round (from the backlog's dependency ordering):

Create `.claude/output/gsd/context-packs/round-N/T-XXX.md` containing ONLY:

```markdown
# Context Pack: T-XXX — {story title}

## Round
{round number}

## Acceptance Criteria
{acceptance criteria from the backlog, as a bullet list}

## Architecture Context
{ONLY the components, layers, and data flows from architecture.yaml that this story touches — NOT the full architecture}

## Codebase Context
{ONLY the relevant subset from the codebase mapping:
- Files this story will read or modify
- Patterns and conventions relevant to this story's domain
- Related existing code that this story interacts with}

## Contracts
{Interfaces, types, API contracts, or schemas this story must use or produce.
If this story depends on contracts from earlier rounds, list them.
If this story produces contracts for later rounds, note that.}

## File Ownership
{Specific files this teammate will own — only these files should be created or modified}

## Dependencies
{Stories this depends on, with a one-line summary of what each provides}

## Estimated Token Budget
{Estimated tokens for this context pack. Must be < 100k.}
```

Create the directory structure:
```
.claude/output/gsd/context-packs/
├── round-1/
│   ├── T-001.md              # Context pack (created by GSD Prep)
│   ├── T-001-feedback.md     # Feedback file (created by Ralph on validation failure, deleted on pass)
│   ├── T-002.md
│   └── ...
├── round-2/
│   └── ...
└── ...
```

Note: `T-XXX-feedback.md` files are NOT created during GSD Prep. They are written by Ralph during Phase D (validation) when a story fails. They are included here for completeness — the context pack directory is the single source of truth for all per-story context, including feedback loops.

## Stage 5: Validation & Report

Verify:
1. **Coverage**: every story in the backlog has a corresponding context pack
2. **Budget**: no context pack exceeds the 100k token budget
3. **Dependencies**: all `depends_on` references point to stories that exist and are in earlier or same rounds
4. **File ownership**: no file is owned by two different stories in the same round
5. **Completeness**: every acceptance criterion from the backlog appears in exactly one context pack

Produce `.claude/output/gsd/prep-report.md`:

```markdown
# GSD Prep Report

## Summary
- Total stories: {count}
- Stories by round: Round 1: {count}, Round 2: {count}, ...
- Stories dropped (already implemented): {count and list}
- Stories reduced (partially implemented): {count and list}
- Stories split (too large): {count, original → sub-stories}
- Average context pack size: {tokens}
- Largest context pack: {story id, tokens}

## Codebase Mapping
- Stack: {languages and frameworks}
- Architecture pattern: {pattern}
- Key concerns: {summary of tech debt or risks}

## Gap Analysis Summary
- Already implemented: {count}
- Partially implemented: {count}
- Purely new: {count}
- Hidden dependencies found: {count}

## Context Packs
| Story | Round | Est. Tokens | Files Owned | Status |
|-------|-------|-------------|-------------|--------|
| T-001 | 1     | ~45k        | 5           | Ready  |
| T-002 | 1     | ~62k        | 8           | Ready  |
| ...   | ...   | ...         | ...         | ...    |

## Validation
- All stories have context packs: {PASS/FAIL}
- All packs within budget: {PASS/FAIL}
- All dependencies satisfiable: {PASS/FAIL}
- No file ownership conflicts: {PASS/FAIL}
```

Present the report to the user for confirmation before proceeding.

If $ARGUMENTS is provided, use it as additional context or instructions: $ARGUMENTS
