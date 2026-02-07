---
description: Documentation rules applied when writing or updating docs and README files
globs: ["docs/**", "**/*.md", "README*"]
---

## Structure and Placement

- Keep documentation close to the code it describes — prefer co-located README files over a monolithic docs folder
- Use clear headings, short paragraphs, and bullet points; documentation should be scannable
- Document prerequisites, setup steps, and required environment variables in the project README

## Content Standards

- Include code examples for non-obvious usage; examples should be copy-pasteable and tested
- Update docs when the related code changes — stale docs are worse than no docs
- For API endpoint documentation, follow the standards in `rules/api.md`

## Inline Code Comments

- Comment "why" (intent, constraints, trade-offs), not "what" (the code already says what)
- Use `TODO(author):` for planned work and `HACK:` for intentional workarounds that need cleanup
- Do not leave commented-out code — use version control instead

## Changelog and Decision Records

- Maintain a `CHANGELOG.md` using Keep a Changelog format; update it with every user-facing change
- Record significant architecture decisions in `docs/adr/` using: Title, Status, Context, Decision, Consequences

## Diagrams

- Use text-based diagram tools (Mermaid, PlantUML) so diagrams live in version control
- Place architecture diagrams in `docs/` and reference them from the README
