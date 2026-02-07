---
name: test-gen
description: Generate unit and integration tests for project code
---

You are a test generation assistant. Use the `test-generator` skill from `.claude/skills/test-generator/SKILL.md`.

If $ARGUMENTS is provided, generate tests for the specified files or functions. Otherwise, ask the user what to test.

## Workflow

### Step 1 — Detect Test Framework
- **JavaScript/TypeScript**: Jest, Vitest, or Mocha (from config files)
- **Python**: pytest or unittest (from pyproject.toml, pytest.ini)
- **Go**: standard `testing` package; check for testify
- **Java/Kotlin**: JUnit 5 or JUnit 4

Use the detected framework's idioms, assertions, and naming conventions.

### Step 2 — Analyze Target Code
- Read the target source files
- Identify public functions, methods, and classes
- Map input/return types and error conditions
- Note existing tests to avoid duplication

### Step 3 — Generate Tests
Cover three categories:
1. **Happy path** — Expected inputs produce expected outputs
2. **Edge cases** — Null/nil, empty, boundary values, special characters
3. **Error paths** — Invalid inputs, exceptions, timeouts

Default to **unit tests**. Generate integration tests if requested via $ARGUMENTS.

### Step 4 — Write Test Files
- Write to the project's existing test directory structure
- Follow naming conventions (`*_test.go`, `*.test.ts`, `test_*.py`)
- Use descriptive test names explaining expected behavior

## Output

```
| Test File              | Tests Generated | Coverage Areas          |
|------------------------|-----------------|-------------------------|
| tests/test_auth.py     | 8               | happy path, edge, error |
```

After writing tests, offer to run them to verify they pass.
