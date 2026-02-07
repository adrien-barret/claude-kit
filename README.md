# Claude Kit (ck)

A Go CLI for managing Claude Code project templates — interactive TUI setup, component management, stack-aware docs generation, and template synchronization.

Built with [Charm](https://charm.sh): Bubble Tea + huh + lipgloss.

---

## Quick Start

### Install from source

```bash
cd claude-cli
make install
```

This builds the `claude-kit` binary (aliased as `ck`), copies it to `/usr/local/bin`, and installs templates to `~/.bmad/templates/`.

### Initialize a project

```bash
cd my-project

# Interactive TUI — pick components from a categorized multi-select
ck init

# AI-guided setup — Claude recommends components based on your project
ck init --plan

# Install to global ~/.claude
ck init --global
```

### Add agents interactively

```bash
# Interactive agent picker — auto-installs skills + rules
ck add

# Add specific agents by name (with their deps)
ck add backend devops

# Add a specific component type
ck add skill code-reviewer
ck add command review
ck add rule testing
```

### Other commands

```bash
ck list                      # See available vs installed components
ck remove                    # Interactive removal picker
ck remove backend            # Remove an agent
ck sync                      # Update installed components from templates
ck docs                      # Generate stack-aware docs-index.md
```

---

## CLI Reference

| Command | Description |
|---------|-------------|
| `ck init` | Interactive setup — categorized multi-select of components |
| `ck init --plan` | AI-guided setup via Claude session |
| `ck init --global` | Install to `~/.claude` |
| `ck add` | Interactive agent picker (auto-installs skills + rules) |
| `ck add <name> [name...]` | Add agents by name with their dependencies |
| `ck add <type> <name>` | Add a specific component (skill, command, rule) |
| `ck remove` | Interactive removal picker |
| `ck remove <name>` | Remove an agent |
| `ck remove <type> <name>` | Remove a specific component |
| `ck list` | Available vs installed side-by-side table |
| `ck list --available` | Available components only |
| `ck list --installed` | Installed components only |
| `ck sync` | Update installed components + refresh docs-index |
| `ck docs` | Generate docs-index.md via stack detection |
| `ck docs --refresh` | Force regenerate even if fresh |
| `ck version` | Print version |

### How `add` works

When you add an agent, it automatically installs:
- The agent definition itself
- All skills listed in the agent's frontmatter
- Related rules based on the agent's role (e.g. backend → code-style, testing, security, api)

```bash
ck add security
# ✓ Added agent: security
#   Installing 4 skill dependencies:
#     ✓ Added skill: security/code-security-audit
#     ✓ Added skill: security/infra-security-audit
#     ✓ Added skill: security/auth-review
#     ✓ Added skill: security/secret-rotation
#   Installing 1 related rules:
#     ✓ Added rule: security
```

### Component types

For explicit type prefixes (`ck add <type> <name>`):

- `agent` / `agents`
- `skill` / `skills`
- `command` / `commands`
- `rule` / `rules`

---

## What's Included

| Category | Components |
|----------|-----------|
| **BMAD Workflow** | Break → Model → Act → Deliver |
| **Agents** | Backend, Tech Lead, DevOps, Security, Pentester, FinOps |
| **Dev Skills** | Code review, test generation, API docs, commit helper, README updater, dependency audit |
| **Security Skills** | Code audit, infra audit, auth review, secret rotation, pentest simulation, threat modeling |
| **FinOps Skills** | Cost optimization, tagging audit, waste detection, budget forecasting |
| **New Skills** | Performance audit, accessibility audit, database review, Terraform review, skill creator |
| **Rules** | Code style, testing, security, API, frontend, infrastructure, documentation, FinOps |

### Slash Commands

**BMAD Workflow:**
`/bmad-run`, `/bmad-break`, `/bmad-model`, `/bmad-act`, `/bmad-deliver`

**Dev Skills:**
`/review`, `/pr-review`, `/test-gen`, `/docs-gen`, `/commit-msg`, `/code-only`

**Security & FinOps:**
`/security-check`, `/pentest`, `/cost-review`

**Roles:**
`/role-backend`, `/role-tech-lead`, `/role-devops`, `/role-security`, `/role-pentester`, `/role-finops`

**Utilities:**
`/ck-sync`

---

## Docs Index

The docs-index system generates compressed, stack-specific notes that stay in Claude's context.

### How it works

1. `ck docs` scans your project root for dependency files (package.json, go.mod, requirements.txt, etc.)
2. Detects your tech stack (languages, frameworks, tools)
3. Generates `.claude/docs-index.md` with framework-specific directives
4. Stores metadata in `.claude/.docs-meta.json` for staleness tracking

### Auto-sync

The docs-index is considered stale when:
- Dependency files have changed (hash mismatch)
- More than 14 days since last generation

`ck sync` automatically refreshes the docs-index after updating components.

### Supported stacks

Languages: JavaScript, TypeScript, Python, Go, Ruby, Rust, Java, Kotlin, PHP
Frameworks: Next.js, React, Vue, Nuxt, Svelte, Angular, Express, Fastify, NestJS, Hono, Django, Flask, FastAPI, Rails, Sinatra, Laravel, Symfony
Tools: Docker, Terraform, Kubernetes, Helm, GitHub Actions, Prisma, Drizzle, Tailwind

---

## Build & Development

### Prerequisites

- Go 1.21+
- Make

### Build

```bash
cd claude-cli
make build              # Compile binary to ./claude-kit
make install            # Build + copy to /usr/local/bin (+ ck alias) + install templates
make install-templates  # Copy templates to ~/.bmad/templates/ only
make clean              # Remove build artifacts
make uninstall          # Remove binary, alias, and templates
```

### Template directory resolution

The binary resolves the template directory in this order:
1. `$BMAD_TEMPLATE_DIR` environment variable
2. `~/.bmad/templates/` (installed via `make install-templates`)
3. Adjacent `project-template/.claude/` (for development from source)

### Go dependencies

- [cobra](https://github.com/spf13/cobra) — subcommand structure
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [huh](https://github.com/charmbracelet/huh) — forms, multi-select, confirm dialogs
- [lipgloss](https://github.com/charmbracelet/lipgloss) — styling, tables, colors

---

## Legacy Installer

The `install.sh` script still works as a fallback. If `ck` / `claude-kit` is available, it delegates automatically:

```bash
# These are equivalent:
bmad-setup --plan           →  ck init --plan
bmad-setup --global         →  ck init --global
bmad-setup                  →  ck init
```

If `ck` is not installed, `install.sh` falls back to the original bash-based installer.

---

## Project Structure

```
claude-cli/
├── cmd/claude-kit/         # Go CLI source
│   ├── main.go             # Cobra root command + version
│   ├── init.go             # ck init — huh multi-select + --plan mode
│   ├── add.go              # ck add — interactive agent picker + auto-deps
│   ├── remove.go           # ck remove — interactive removal + warnings
│   ├── list.go             # ck list — lipgloss table
│   ├── sync.go             # ck sync — update + docs refresh
│   └── docs.go             # ck docs — stack detection + generation
├── internal/
│   ├── catalog/            # Template scanning + component operations
│   ├── stack/              # Stack detection from dependency files
│   ├── docsindex/          # Docs-index generation + staleness
│   └── config/             # Path resolution + defaults
├── project-template/.claude/  # Template files
│   ├── CLAUDE.md           # Project memory
│   ├── settings.json       # Permissions
│   ├── agents/             # 6 agent role definitions
│   ├── skills/             # 18+ skill directories
│   ├── commands/           # 18 slash commands
│   └── rules/              # 8 project rules
├── go.mod / go.sum
├── Makefile                # build, install, install-templates, clean
├── install.sh              # Legacy wrapper → delegates to ck
├── prompts.sh              # AI-guided setup prompts (used by --plan)
└── README.md
```

---

## Skills Reference

### Dev Skills

| Skill | Description |
|-------|-------------|
| `code-reviewer` | Code review with severity levels (critical/warning/info), auto-fix suggestions |
| `test-generator` | Test generation with framework detection and coverage gap analysis |
| `api-documenter` | OpenAPI/Swagger documentation generation |
| `git-commit-helper` | Conventional commit message generation |
| `readme-updater` | Keep README in sync with code |
| `dependency-auditor` | Vulnerability scanning, license compatibility matrix, supply-chain risk scoring |

### Security Skills

| Skill | Description |
|-------|-------------|
| `security` | Orchestrator — runs all security sub-skills |
| `security/code-security-audit` | OWASP Top 10, injection, XSS, hardcoded secrets |
| `security/infra-security-audit` | Cloud config, permissions, encryption |
| `security/auth-review` | OAuth/JWT, RBAC, token policies |
| `security/secret-rotation` | Secret storage and rotation policies |
| `security/pentest-web` | Auth bypass, IDOR, SSRF, rate-limit bypass, JWT attacks |
| `security/threat-model` | STRIDE threat modeling |

### FinOps Skills

| Skill | Description |
|-------|-------------|
| `finops` | Orchestrator — runs all FinOps sub-skills |
| `finops/cost-optimization` | Rightsizing, auto-scaling, reserved instances |
| `finops/tagging-audit` | Cost allocation tag compliance |
| `finops/waste-detection` | Idle resources, orphaned volumes |
| `finops/budget-forecast` | Cost estimation from IaC |

### New Skills

| Skill | Description |
|-------|-------------|
| `performance-audit` | N+1 queries, bundle size, caching, lazy loading, connection pooling |
| `accessibility-audit` | WCAG 2.1 AA, ARIA, keyboard nav, contrast, screen reader |
| `database-review` | Schema, indexing, query optimization, migration safety |
| `terraform-review` | Module structure, state management, provider versioning |
| `skill-creator` | Meta-skill to generate new SKILL.md files |

---

## How to Describe Your Application

To use the BMAD workflow, provide a **project brief**:

```
1. PROJECT NAME
2. PROBLEM STATEMENT — what it solves, who it's for
3. TECH STACK — language, framework, database, cloud
4. CORE FEATURES — prioritized list
5. CONSTRAINTS — performance, compliance, multi-tenancy
6. INTEGRATIONS — external APIs, payment, notifications
7. INFRASTRUCTURE — deployment, CI/CD, containers
```

Then run `/bmad-run` for the full workflow, or phase by phase:

```bash
/bmad-break       # Define the problem → problem.yaml
/bmad-model       # Design architecture → architecture.yaml, backlog.yaml
/bmad-act         # Implement code from backlog
/bmad-deliver     # Prepare release → release-notes.md
```

---

## Rules

Rules are modular project instructions loaded based on file patterns:

| Rule | Globs | What it enforces |
|------|-------|-----------------|
| `code-style.md` | `src/**`, `lib/**`, `app/**` | DRY, KISS, SOLID, clean code |
| `testing.md` | `tests/**`, `**/*.test.*`, `**/*.spec.*` | Test-first, edge cases, independent tests |
| `security.md` | _(all files)_ | No secrets, input validation, least privilege |
| `api.md` | `src/routes/**`, `src/api/**`, `src/controllers/**` | REST conventions, pagination, error format |
| `frontend.md` | `src/components/**`, `**/*.tsx`, `**/*.jsx` | Small components, accessibility, state handling |
| `infrastructure.md` | `infra/**`, `*.tf`, `Dockerfile*`, `k8s/**` | IaC, least-privilege IAM, non-root containers |
| `documentation.md` | `docs/**`, `**/*.md` | Close to code, examples, keep updated |
| `finops.md` | `infra/**`, `*.tf`, `k8s/**`, `helm/**` | Tagging, rightsizing, lifecycle, scheduling |
