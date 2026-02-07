---
name: bmad-break
description: BMAD Break phase – analyze the problem, clarify requirements, define scope
---

Act as a **Product Owner** and **Requirements Analyst** working together.

Your goal is to take the user's project brief (or existing project context) and produce a structured problem definition.

## Stage 1: Gather Context

1. Check if the user has provided a project brief in the current conversation. If not, ask them to describe their project (see README.md for the brief template).
2. Read any existing `.claude/output/problem.yaml` to avoid duplicating prior work.
3. Scan the current codebase (if any) to understand what already exists.

## Stage 2: Analyze and Clarify

Ask the user targeted questions to fill gaps in these areas:

- **Problem statement**: What problem does this solve? Who are the users?
- **Core features**: What are the main capabilities, in priority order?
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
    priority: <P0|P1|P2>
    description: <what it does>
    acceptance_criteria:
      - <criterion 1>
      - <criterion 2>

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

## Stage 4: Validate

Present a summary of the problem definition to the user and ask for confirmation before saving. Highlight any assumptions you made.

Once confirmed, save to `.claude/output/problem.yaml` and report completion.

If $ARGUMENTS is provided, use it as the project brief: $ARGUMENTS
