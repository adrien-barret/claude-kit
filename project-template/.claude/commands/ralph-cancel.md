---
name: ralph-cancel
description: Cancel Ralph Loop — remove stop hook and show progress summary
---

Cancel the **Ralph Loop** and show a progress summary.

## Steps

1. Read `.claude/ralph-prd.json` and summarize progress:
   - Total stories
   - Completed stories (where `passes` is `true`)
   - Remaining stories (where `passes` is `false`)
   - List each story with its status

2. Remove the Ralph stop hook from `.claude/settings.json`:
   - Remove the `hooks.Stop` entry that references `.claude/hooks/ralph-stop.sh`
   - Remove the `Bash(.claude/hooks/ralph-stop.sh)` permission from `permissions.allow`
   - Preserve all other settings

3. Remove the iteration counter file `.claude/.ralph-iteration` if it exists.

4. Report what was done:
   - Show the progress summary
   - Confirm the stop hook has been removed
   - Note that `.claude/ralph-prd.json` and `.claude/ralph-prompt.md` are preserved for reference
   - Suggest the user can re-run `/ralph` to restart the loop if needed

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
