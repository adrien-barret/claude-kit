---
name: commit-msg
description: Generate a conventional commit message from staged git changes
---

You are a commit message generator. Use the `git-commit-helper` skill from `.claude/skills/git-commit-helper/SKILL.md`.

## Instructions

1. Run `git diff --cached --stat` to see which files are staged. If nothing is staged, tell the user to stage changes first and stop.
2. Run `git diff --cached` to read the full staged diff.
3. Analyze the changes and generate a conventional commit message.

## Commit Format

Use the conventional commits specification:

```
<type>(<scope>): <short summary>

<body — optional, for complex changes>

<footer — optional, for breaking changes or issue refs>
```

### Types
- `feat` — new feature or capability
- `fix` — bug fix
- `chore` — maintenance, deps, config
- `docs` — documentation only
- `refactor` — code restructuring without behavior change
- `test` — adding or updating tests
- `ci` — CI/CD pipeline changes

### Scope
Include a scope when changes are confined to a specific module or area:
- `feat(api): add pagination to /users endpoint`
- `fix(auth): handle expired refresh tokens`
- `chore(deps): bump axios to 1.7.0`

### Body
For complex changes, add a blank line after the summary and explain **why** the change was made, not what (the diff shows that).

### Breaking Changes
If the change breaks backward compatibility, add a `BREAKING CHANGE:` footer:
```
feat(api)!: change response format for /users

BREAKING CHANGE: response is now paginated; clients must handle `next_cursor` field
```

## Output

Present the generated commit message in a code block, ready to copy. If the diff spans multiple areas, offer a multi-line version with body. Ask the user if they want to proceed with `git commit` or adjust the message.

If $ARGUMENTS is provided, use it as additional context for the commit message.
