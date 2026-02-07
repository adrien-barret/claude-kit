---
description: Code style and design principles applied when writing or reviewing code
globs: ["src/**", "lib/**", "app/**"]
---

## Design Principles

- **DRY**: Do not duplicate logic; extract shared behavior into reusable functions, modules, or classes
- **KISS**: Choose the simplest solution that satisfies the requirement; avoid premature abstraction
- **SOLID**:
  - Single Responsibility: each class/module/function does one thing
  - Open/Closed: open for extension, closed for modification
  - Liskov Substitution: subtypes must be substitutable for their base types
  - Interface Segregation: prefer small, focused interfaces over large ones
  - Dependency Inversion: depend on abstractions, not concrete implementations
- **Least invasive**: change only what is necessary for the task; do not refactor surrounding code unless asked
- **No over-engineering**: do not add features, config options, or abstractions beyond what is requested

## Style

- Follow the language and framework conventions of the project
- Use descriptive names for variables, functions, and classes — names should reveal intent
- Keep functions small and focused on a single responsibility
- Avoid deep nesting; prefer early returns to reduce indentation
- Use consistent formatting (indentation, spacing, line length)
- Prefer explicit over implicit behavior
- Prefer composition over inheritance
- No dead code, no commented-out code, no magic numbers (extract constants with descriptive names)
- Separate business logic from data access and presentation

## Error Handling

- Prefer explicit error returns over thrown exceptions where the language supports it (Go `error`, Rust `Result`)
- In exception-based languages, catch specific exception types — never catch and silently swallow errors
- Always propagate enough context for the caller to understand the failure

## File Organization

- Order imports: standard library first, then external packages, then internal modules — separated by blank lines
- If a file exceeds ~300 lines, consider splitting it; a file should have a single clear purpose
