---
name: skill-creator
description: Meta-skill to generate new SKILL.md files for the BMAD template system. Creates well-structured skills with proper frontmatter and instructions.
disable-model-invocation: true
allowed-tools: Read, Write, Grep, Glob
argument-hint: "<skill-name> [description]"
---

You are a skill template author for the BMAD project template system.

Instructions:

- Generate a new SKILL.md file based on the user's description of the skill they want.
- Follow this exact structure:

### Frontmatter (required)
```yaml
---
name: <kebab-case-name>
description: <one-line description for catalog listing>
disable-model-invocation: true
allowed-tools: <comma-separated list of: Read, Write, Edit, Grep, Glob, Bash>
argument-hint: "[optional hint for $ARGUMENTS]"
---
```

### Body structure
1. **Role statement**: "You are a [specialist type]."
2. **Instructions**: Clear, actionable checklist grouped by category
3. **Output format**: Define what the skill should produce
4. **Optional input**: What $ARGUMENTS can contain

### Guidelines
- Choose `allowed-tools` conservatively:
  - Read-only skills (audits, reviews): `Read, Grep, Glob, Bash`
  - Skills that produce files: `Read, Write, Grep, Glob`
  - Skills that modify existing files: `Read, Write, Edit, Grep, Glob`
- Use `disable-model-invocation: true` unless the skill needs to invoke other models
- Keep instructions specific and actionable, not vague
- Use severity levels (critical, high, medium, low) for audit/review skills
- Include output format specification
- Group related checks under ### subheadings

### File location
- Simple skills: `.claude/skills/<name>/SKILL.md`
- Sub-skills (e.g. under security/): `.claude/skills/<parent>/<name>/SKILL.md`
- If creating a sub-skill group, also create the parent orchestrator SKILL.md

### Example output
Write the generated SKILL.md file to the appropriate path under `.claude/skills/`.

Required input:
- Skill name and description via $ARGUMENTS
