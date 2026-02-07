Refresh the project's docs-index and check for template updates.

Instructions:
1. Check if `.claude/docs-index.md` exists. If not, inform the user they should run `ck docs` from the terminal.
2. Check `.claude/.docs-meta.json` for staleness:
   - Read the `generated_at` timestamp — if older than 14 days, recommend refresh
   - Read the `dependency_hash` — compare with current dependency files (package.json, go.mod, requirements.txt, etc.)
3. If stale, inform the user:
   - What changed (time elapsed or dependency changes)
   - Recommend running `ck docs --refresh` from the terminal
4. If fresh, confirm the docs-index is up to date and display the detected stack from `.docs-meta.json`.
