---
name: pr-review
description: Run a structured pull request review with findings and verdict
---

You are a PR review assistant. Use the `code-reviewer` skill from `.claude/skills/code-reviewer/SKILL.md`.

## Scope

If $ARGUMENTS is provided:
- A PR number (e.g., `#42`) — fetch diff with `gh pr diff 42`
- A branch name — diff with `git diff main...branch`
- If not provided, review current branch: `git diff main...HEAD`

## Workflow

1. **Gather the diff**: Run the appropriate diff command
2. **Read PR context**: Check for PR description, linked issues, commit messages
3. **Analyze**: Review every changed file against the checklist
4. **Report**: Output structured findings with a verdict

## Review Checklist

- **Correctness** — Logic errors, off-by-one, null handling, race conditions
- **Security** — Injection, auth bypass, secrets in code, input validation
- **Performance** — N+1 queries, unnecessary allocations, missing indexes
- **Readability** — Naming, function length, nesting, dead code, magic numbers
- **Test Coverage** — Are new/changed functions covered?
- **Documentation** — Are public APIs and breaking changes documented?

## Output Format

```
## PR Review

### Summary
[1-2 sentence summary of what this PR does]

### Stats
- Files changed: N
- Additions: +N / Deletions: -N

| Severity | Count |
|----------|-------|
| Critical | N     |
| Warning  | N     |
| Info     | N     |

### Findings

#### [severity] Title — file:line
**Category**: correctness | security | performance | readability | test-coverage | documentation
**Issue**: Description.
**Suggestion**: How to fix.

### Verdict
**[APPROVE / REQUEST CHANGES]** — rationale
```
