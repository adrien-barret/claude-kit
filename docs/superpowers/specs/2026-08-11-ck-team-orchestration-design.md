# `ck team` — multi-role pane orchestration (design)

**Status:** approved (brainstorming) — 2026-08-11
**Scope:** one new `ck` subcommand + one internal package + role/board templates.

## Problem

The "Kerios company" runs several Claude Code sessions (one tmux window per
department: `lead`, `eng`, `marketing`, `sales`, `board`) coordinating through a
shared `BOARD.md`, with global hooks that (a) **poke** a pane when an `@<role>:`
request is written and (b) **rescan** pending requests when a pane boots. That
machinery is hardcoded to Kerios: fixed board path, fixed pane names, standalone
Python scripts in `~/nix-config/kerios/`, wired into the global
`~/.claude/settings.json`.

We want the same mechanism as a **reusable, project-agnostic** `ck` command, so
any project can spin up its own role/board/pane team without touching Kerios.

## Goals / non-goals

- **Goal:** `ck team <ProjectName>` scaffolds a root project dir with one
  subfolder per role, a shared board, per-role `CLAUDE.md`, generic hook wiring,
  and a launched tmux session.
- **Goal:** roles are discovered by a **temporary headless Claude session** that
  proposes a role set + per-role objectives from a project description; the user
  confirms/edits; `ck` does all file writes.
- **Goal:** the hook machinery is **generic** and **coexists with Kerios** — it
  no-ops in any directory that isn't a `ck team` project.
- **Non-goal:** modifying Kerios's board, scripts, or wiring.
- **Non-goal:** per-role git repos (Kerios uses sibling repos; we use subfolders).
- **Non-goal:** replacing `ck`'s existing `internal/stack` tech-stack detector —
  unrelated; the cobra command is named `team` precisely to avoid that collision.

## User-facing flow

```
ck team MyProject
  1. prompt: one-paragraph project description
  2. headless `claude -p` → JSON proposal { board_focus, roles:[{name, folder, mission, objectives[]}] }
  3. huh confirm/edit form (toggle / rename / edit objectives)   ← ck is the only writer
  4. scaffold role subfolders + per-role .claude/CLAUDE.md (archetype template + objectives)
  5. write board/BOARD.md + board/RESOURCE-REQUESTS.md + .team/config.json
  6. idempotently wire two generic hooks into global ~/.claude/settings.json
  7. create + launch the tmux session (one window per role, cd'd into its folder)
```

**Fallback (no `claude` on PATH, or bad/empty JSON):** skip step 2 and drive
step 3 from `huh` alone — pick archetypes from the catalog, type one objective
line each. No hard dependency on Claude.

## On-disk layout produced

```
MyProject/
├── .team/config.json         # discovery anchor: session, board path, role→folder map
├── board/
│   ├── BOARD.md              # Current Focus, per-role sections, Cross-Role Requests (@role: + DoD)
│   └── RESOURCE-REQUESTS.md
├── lead/      .claude/CLAUDE.md
├── eng/       .claude/CLAUDE.md
└── …          one subfolder per chosen role
```

- Roles are plain subdirectories (not separate git repos).
- `ck team` does **not** `git init`.
- **Idempotent re-run:** adds missing roles, re-asserts `.team/config.json` +
  hook wiring + tmux windows; never clobbers an existing role `CLAUDE.md`
  (skips with a note).

## `.team/config.json` (the contract)

```json
{
  "project": "MyProject",
  "tmux_session": "myproject",
  "board": "board/BOARD.md",
  "roles": [
    { "name": "lead", "folder": "lead" },
    { "name": "eng",  "folder": "eng"  }
  ]
}
```

Everything the hooks need lives here. `board` is relative to the config's
directory (the project root). A pane's role = the `folder` its `cwd` is under.

## Hooks — implemented as hidden `ck` subcommands (not Python)

`ck team` merges exactly two entries into global `~/.claude/settings.json`
(only if its exact command string is not already present):

```json
{ "hooks": {
  "SessionStart": [ { "hooks": [ { "type": "command", "command": "ck team hook rescan" } ] } ],
  "PostToolUse":  [ { "matcher": "Write|Edit",
                      "hooks": [ { "type": "command", "command": "ck team hook notify" } ] } ]
} }
```

Both subcommands read the Claude hook JSON on **stdin** (which carries `cwd`),
then **walk up from `cwd`** to find the nearest `.team/config.json`.

- `ck team hook notify`: if the just-written file is the board, grep it for
  `@<role>:` where `<role>` ∈ config roles, and for each run
  `tmux send-keys -t <session>:<role> "[team] new @<role>: request — re-read <board>" C-m`.
- `ck team hook rescan`: map `cwd`'s folder → role, print the board's
  un-`[DONE]` `@<role>:` lines so the booting session picks them up.
- **No `.team/config.json` found → print nothing, exit 0.**

**Coexistence with Kerios:** Kerios repos have no `.team/config.json`, so the
generic hooks no-op there; Kerios's own Python hooks are never touched. In a
`ck team` project, Kerios's hooks find no kerios pane/dept and no-op. Both hook
sets run globally without collision. Hidden subcommands (not standalone scripts)
keep everything in one testable Go binary with no Python dependency.

## Components (file map)

**New package `internal/team/`** (pure logic, unit-tested):
- `config.go` — `Config` struct, `Write`, `Discover(startDir)` (walk-up to
  `.team/config.json`), `RoleForDir(cfg, dir)`.
- `archetype.go` — load a role archetype template from the template dir, render
  with `{{.Mission}}`/`{{.Objectives}}`/`{{.Role}}`/`{{.Board}}`; generic
  fallback for unknown roles.
- `board.go` — render initial `BOARD.md` + `RESOURCE-REQUESTS.md` from the role
  set + `board_focus`.
- `propose.go` — build the claude prompt; `Propose(desc)` shells out to
  `claude -p … --output-format text`, extract the JSON object, validate; typed
  error so the caller can fall back.
- `settings.go` — idempotent merge of the two hook entries into
  `~/.claude/settings.json` (back up before write; safe-create if absent).
- `notify.go` / `rescan.go` — pure functions: `(boardContent, cfg, cwd) →
  []tmuxTarget` and `→ []pendingLine`. IO (`tmux`, stdout) behind a thin seam.

**New command `cmd/claude-kit/team.go`:**
- `teamCmd` (`ck team <ProjectName>`) — orchestrates the flow; huh confirm/edit
  using `ckTheme()`; styled output via `banner()`/`sectionHeader()`/`checkMark`.
- `team hook notify` / `team hook rescan` hidden subcommands (read stdin, call
  the pure `internal/team` logic, perform IO).
- tmux launch helper (`tmux new-session -d -s <s> -n <role> -c <folder>` + one
  `new-window` per remaining role); if `tmux` absent, print the command instead.
- Registered via `rootCmd.AddCommand(teamCmd)` in `main.go`.

**New templates under `project-template/team-roles/`:**
- `<archetype>/CLAUDE.md.tmpl` for `lead`, `board`, `eng`, `marketing`, `sales`,
  `design`, `research`, `ops`, and a `generic/CLAUDE.md.tmpl` fallback. Each
  carries the board/coordination protocol + "you are poked on `@<role>:`" text
  with placeholders for mission/objectives.
- `board/BOARD.md.tmpl`, `board/RESOURCE-REQUESTS.md.tmpl`.

## Data flow

```
description ──claude -p──▶ JSON proposal ──huh confirm──▶ RoleSet
RoleSet ──archetype.Render──▶ per-role CLAUDE.md
RoleSet ──board.Render──▶ BOARD.md
RoleSet ──▶ .team/config.json
config ──settings.Merge──▶ ~/.claude/settings.json (2 hooks, idempotent)
RoleSet ──tmux──▶ live session
(runtime) board Edit ──PostToolUse──▶ ck team hook notify ──walk-up config──▶ tmux poke
(runtime) pane boot ──SessionStart──▶ ck team hook rescan ──walk-up config──▶ pending lines
```

## Error handling

- No `claude` / non-JSON / empty roles → `huh` fallback; never abort.
- No `tmux` → finish scaffold + wiring, print the exact launch command.
- Global `settings.json` absent/malformed → create or merge safely; back up first.
- Re-run on an existing project → additive, non-destructive.
- `hook notify`/`hook rescan` with no config found → exit 0, no output.

## Testing

- **Unit (Go):** config walk-up discovery; archetype rendering (incl. generic
  fallback); claude-JSON extraction + validation; idempotent settings.json merge
  (absent / present / already-wired); `notify`/`rescan` pure logic
  (board+config+cwd → expected targets / pending lines).
- **E2E:** run `ck team` in a temp dir with a **fake `claude` on PATH** returning
  canned JSON → assert full tree + board + `.team/config.json` + settings merge;
  plus a `ck team hook notify` invocation with a synthetic stdin event asserting
  the computed tmux targets (tmux behind an interface so no real session needed).

## Conventions

Go 1.25 + cobra; module `github.com/AdeptMind/infra-tool/claude-cli`. Reuse
`config.TemplateDir()`, `styles.go` helpers, `ckTheme()`. Headless call mirrors
`smartadd.go` (`claude -p … --output-format text` + JSON extraction). Build/test
via `make build` / `make test`. Conventional commits; ponytail pass before PR.
