---
description: Testing rules applied when generating or reviewing tests
globs: ["tests/**", "test/**", "**/*.test.*", "**/*.spec.*"]
---

## Test Design

- Write a failing test before fixing a bug to prove the bug exists and prevent regression
- Cover edge cases and error paths, not just happy paths — test nulls, empty inputs, boundary values
- Use descriptive test names explaining the scenario and expected outcome: `test_returns_404_when_user_not_found`
- Each test must be independent and not rely on execution order or shared mutable state

## Mocking Strategy

- Mock external boundaries (APIs, databases, file systems, clocks) to keep tests fast and deterministic
- Avoid mocking internal code — test real behavior; mocking internals makes tests brittle
- When a test requires complex setup, the code under test likely needs a simpler interface

## Test Data

- Use factories or builders to create test data with sensible defaults; override only what the test cares about
- Do not share mutable test data between tests — each test creates its own state
- Avoid hardcoded fixture files that grow stale; generate test data programmatically

## Test Types and Coverage

- **Unit tests**: test individual functions and classes in isolation; majority of tests
- **Integration tests**: test boundaries between your code and external systems
- **End-to-end tests**: cover critical user flows only; keep count low (slow)
- Focus coverage on critical paths and complex logic; do not chase arbitrary percentage targets

## CI and Reliability

- All tests must pass before merge; no exceptions
- Fix flaky tests immediately — do not skip, retry, or `@ignore`; flaky tests erode trust
- Keep test setup minimal and close to the assertion
