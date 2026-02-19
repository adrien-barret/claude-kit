---
name: ralph-loop
description: Resume Ralph — continue implementing remaining stories from ralph-prd.json
---

Resume **Ralph** — pick up where the last session left off.

## Setup Check

1. Read `.claude/ralph-prd.json`. If it doesn't exist, tell the user to run `/ralph` first.
2. Read `.claude/output/architecture.yaml` if it exists (for teammate context).
3. Read `.claude/output/gsd/prep-report.md` if it exists — use context packs from `.claude/output/gsd/context-packs/` when spawning teammates (each teammate gets their story's bounded context pack instead of the full architecture dump).
4. Show progress: X/Y stories done, list remaining stories grouped by round.

## Resume

1. Find all stories where `passes` is `false`
2. Recompute remaining rounds (stories whose `dependsOn` are all passing form the next round)
3. Check for existing feedback files in `.claude/output/gsd/context-packs/round-N/T-XXX-feedback.md` — stories with feedback files were mid-fix when the session was interrupted. These teammates should be re-spawned with both their context pack AND their feedback file, so they know exactly what failed and what was already attempted.
4. Follow the same process as `/ralph` Step 6:
   - Create an agent team, you are the lead in delegate mode
   - Spawn one teammate per story in each round, using the `/ralph` teammate spawn prompt
   - **For stories with feedback files**: include the feedback in the spawn prompt so the teammate picks up where the previous attempt left off — no wasted work rediscovering the same issues
   - Each teammate must:
     a. Read the project codebase to understand existing conventions
     b. If a feedback file exists: address ALL listed issues first
     c. Implement the story following ALL acceptance criteria
     d. Write tests (unit + integration where applicable)
     e. Run tests and verify they ALL pass
     f. Use skills when relevant (`/review`, `/test-gen`, `/security-check`, check `.claude/skills/`)
     g. Commit with: `feat(<story-id>): <title>`
     h. Message the lead with a completion report — do NOT update `.claude/ralph-prd.json` directly. The lead validates and marks stories as passed.
   - Wait for all teammates in a round to finish before starting the next round
   - Validate each story (acceptance-validator) and update `.claude/ralph-prd.json` as stories pass
5. After all stories pass, run quality checks:
   - Code review (use code-reviewer skill)
   - Full test suite execution
   - Security scan (use security skill)
   - Dependency audit (use dependency-auditor skill)
6. Update `.claude/output/act-report.md`

See `/ralph` for full details on teammate spawn prompt, file conflict prevention, and quality checks.

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
