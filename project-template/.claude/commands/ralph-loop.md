---
name: ralph-loop
description: Resume Ralph — continue implementing remaining stories from ralph-prd.json
---

Resume **Ralph** — pick up where the last session left off.

## Setup Check

1. Read `.claude/ralph-prd.json`. If it doesn't exist, tell the user to run `/ralph` first.
2. Read `.claude/output/architecture.yaml` if it exists (for teammate context).
3. Show progress: X/Y stories done, list remaining stories grouped by round.

## Resume

1. Find all stories where `passes` is `false`
2. Recompute remaining rounds (stories whose `dependsOn` are all passing form the next round)
3. Follow the same process as `/ralph` Step 6:
   - Create an agent team, you are the lead in delegate mode
   - Spawn one teammate per story in each round, using the `/ralph` teammate spawn prompt
   - Each teammate must:
     a. Read the project codebase to understand existing conventions
     b. Implement the story following ALL acceptance criteria
     c. Write tests (unit + integration where applicable)
     d. Run tests and verify they ALL pass
     e. Use skills when relevant (`/review`, `/test-gen`, `/security-check`, check `.claude/skills/`)
     f. Update `.claude/ralph-prd.json` — set `passes` to `true` for their story
     g. Commit with: `feat(<story-id>): <title>`
     h. Message the lead with a completion report
   - Wait for all teammates in a round to finish before starting the next round
   - Update `.claude/ralph-prd.json` as stories pass
4. After all stories pass, run quality checks:
   - Code review (use code-reviewer skill)
   - Full test suite execution
   - Security scan (use security skill)
   - Dependency audit (use dependency-auditor skill)
5. Update `.claude/output/act-report.md`

See `/ralph` for full details on teammate spawn prompt, file conflict prevention, and quality checks.

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
