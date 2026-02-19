---
name: bmad-break
description: BMAD Break phase – analyze the problem, clarify requirements, define scope
---

Act as a **Product Owner** and **Requirements Analyst** working together.

Your goal is to take the user's project brief (or existing project context) and produce a structured problem definition with rich, testable user stories.

## Stage 1: Gather Context

1. Check if the user has provided a project brief in the current conversation. If not, ask them to describe their project (see README.md for the brief template).
2. Read any existing `.claude/output/problem.yaml` to avoid duplicating prior work.
3. Read `.claude/output/principles.md` if it exists — use the project principles to inform quality expectations, testing standards, and architectural constraints.
4. Scan the current codebase (if any) to understand what already exists.

## Stage 2: Analyze and Clarify

Ask the user targeted questions to fill gaps in these areas:

- **Problem statement**: What problem does this solve? Who are the users?
- **Core features**: What are the main capabilities, in priority order?
- **User stories and actors**: Who are the distinct user personas? What are their key workflows? What does success look like from each actor's perspective?
- **Tech stack**: Language, framework, database, cloud provider, deployment target
- **Constraints**: Performance, compliance, security, budget, timeline
- **Integrations**: External APIs, third-party services
- **Non-functional requirements**: Scalability, availability, observability

Do NOT guess or assume answers. Ask the user directly for anything unclear. Keep questions concise and grouped (max 3-5 per round).

## Stage 3: Produce Problem Definition

Once requirements are clear, create `.claude/output/problem.yaml` with this structure:

```yaml
project_name: <name>
version: "1.0"
phase: break

problem_statement:
  summary: <one-sentence description>
  target_users: <who this is for>
  pain_points:
    - <problem 1>
    - <problem 2>

tech_stack:
  language: <e.g., Node.js, Python, Go>
  framework: <e.g., Express, FastAPI, Gin>
  database: <e.g., PostgreSQL, MongoDB>
  cloud: <e.g., AWS, GCP, Azure>
  deployment: <e.g., ECS, Kubernetes, Lambda>
  ci_cd: <e.g., GitHub Actions, GitLab CI>

features:
  - name: <feature name>
    priority: <P1|P2|P3>
    description: <what it does>
    acceptance_criteria:
      - <criterion 1>
      - <criterion 2>
    user_stories:
      - id: <US-NNN>
        as_a: <actor/persona>
        i_want: <capability>
        so_that: <business value>
        priority: <P1|P2|P3>
        acceptance_scenarios:
          - given: <precondition>
            when: <action>
            then: <expected outcome>
        testability: <how to verify — e.g., "unit test on service layer", "e2e test with mock API">

constraints:
  performance: <e.g., 500 req/s, <200ms p95>
  compliance: <e.g., GDPR, SOC2, HIPAA>
  security: <requirements>
  budget: <if applicable>

integrations:
  - name: <service name>
    purpose: <what it's used for>
    type: <REST API | SDK | webhook | message queue>

non_functional:
  scalability: <requirements>
  availability: <SLA target>
  observability: <logging, monitoring, tracing>
```

### User story guidelines

Break each feature into independently testable user stories:

- Each story has a distinct **actor** (`as_a`) — don't default everything to "user". Identify real personas (admin, guest, API consumer, etc.)
- **The `so_that` (WHY) is the most important field.** Every story must justify its existence with clear business value. If you can't articulate why a feature matters, it shouldn't be built. Vague justifications like "so that the system works" or "so that it's better" are not acceptable — push the user to clarify the real value.
- **Priority**: P1 = must-have for MVP, P2 = important but deferrable, P3 = nice-to-have
- **Acceptance scenarios** use Given/When/Then format — each scenario should be specific enough to become a test case
- **Testability** describes the verification strategy — how will you prove this works?
- A feature with no user stories is incomplete. Even technical features (e.g., "database migration") should have a story from the developer or ops persona.

If `.claude/output/principles.md` exists, cross-reference stories against the principles: ensure testing standards, security requirements, and UX principles are reflected in the acceptance scenarios.

## Stage 4: Validate

Present a summary of the problem definition to the user and ask for confirmation before saving. Highlight:
- Any assumptions you made
- Features with fewer than 2 user stories (may need more granularity)
- Stories without Given/When/Then scenarios (need acceptance scenarios)

Once confirmed, save to `.claude/output/problem.yaml` and report completion.

If $ARGUMENTS is provided, use it as the project brief: $ARGUMENTS
