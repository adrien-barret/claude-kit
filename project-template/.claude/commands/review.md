---
name: review
description: Run a structured code review on staged changes or specific files
---

You are a code review assistant. Use the `code-reviewer` skill from `.claude/skills/code-reviewer/SKILL.md`.

## Scope

If $ARGUMENTS is provided, review the specified files or directories. Otherwise, review staged changes via `git diff --cached`. If nothing is staged, review unstaged changes via `git diff`.

## Workflow

1. **Gather changes**: Run the appropriate git diff or read the specified files
2. **Analyze**: Review each changed file against the checklist below
3. **Report**: Output structured findings

## Review Checklist

- **Correctness** — Logic errors, off-by-one, null/undefined handling, race conditions
- **Security** — Injection, auth bypass, secrets in code, input validation
- **Performance** — N+1 queries, unnecessary allocations, blocking calls in async code
- **Readability** — Naming, function length, nesting depth, dead code, magic numbers
- **Test Coverage** — Are changed functions covered by tests?

Also check against project principles (from CLAUDE.md): DRY, KISS, SOLID, least invasive.

## Output Format

```
## Code Review Summary

| Severity | Count |
|----------|-------|
| Critical | N     |
| Warning  | N     |
| Info     | N     |

### [severity] Title — file:line
**Category**: correctness | security | performance | readability | test-coverage
**Issue**: Description of the problem.
**Suggestion**: How to fix it.
**Auto-fix**: [auto-fixable] (if unambiguous, provide corrected code)

## Verdict
[APPROVE / REQUEST CHANGES] — brief rationale
```

For `[auto-fixable]` findings, ask the user if they want to apply fixes automatically.
