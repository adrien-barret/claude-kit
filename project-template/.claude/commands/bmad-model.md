---
name: bmad-model
description: BMAD Model phase – design architecture, produce ADRs, and generate a prioritized backlog
---

Act as an **Architect** and **Tech Lead** working together.

Your goal is to take the problem definition from the Break phase and produce an architecture design and implementation backlog.

## Prerequisites

Read `.claude/output/problem.yaml`. If it does not exist, tell the user to run `/bmad-break` first and stop.

## Stage 1: Architecture Design

Based on the problem definition, design the system architecture:

1. **Component breakdown**: Identify the major components/services and their responsibilities.
2. **Data model**: Define core entities, their relationships, and storage strategy.
3. **API surface**: Define the main endpoints or interfaces between components.
4. **Infrastructure**: Define the deployment topology, networking, and cloud services.
5. **Security model**: Authentication, authorization, data protection approach.
6. **Cross-cutting concerns**: Logging, monitoring, error handling, configuration.

Follow code principles: KISS (simplest architecture that meets requirements), no over-engineering (do not add components or layers that are not justified by the requirements).

## Stage 2: Architecture Decision Records

For each significant decision, document the reasoning:

- What was decided and why
- Alternatives considered
- Trade-offs accepted

## Stage 3: Generate Backlog

Break the architecture into an ordered backlog of implementation tasks. Each task must be:

- **Small enough** to implement in a single session
- **Self-contained** with clear inputs and outputs
- **Ordered** by dependency (foundations first, then features, then polish)
- **Labeled** by component and priority

## Stage 4: Produce Output

Create `.claude/output/architecture.yaml`:

```yaml
project_name: <from problem.yaml>
version: "1.0"
phase: model

components:
  - name: <component name>
    type: <service | library | infrastructure | config>
    responsibility: <what it does>
    tech: <specific technology>
    depends_on:
      - <other component>

data_model:
  entities:
    - name: <entity>
      fields:
        - name: <field>
          type: <type>
          constraints: <nullable, unique, indexed, etc.>
      relationships:
        - type: <has_many | belongs_to | has_one>
          target: <entity>

api_surface:
  - method: <GET|POST|PUT|DELETE>
    path: <endpoint path>
    component: <which component owns this>
    description: <what it does>
    auth: <public | authenticated | admin>

infrastructure:
  compute: <ECS, Lambda, K8s, etc.>
  database: <RDS, DynamoDB, etc.>
  cache: <ElastiCache, Redis, etc.>
  storage: <S3, etc.>
  networking: <VPC layout, load balancer>
  ci_cd: <pipeline description>

security:
  authentication: <strategy>
  authorization: <strategy>
  encryption: <at rest, in transit>
  secrets_management: <approach>

adrs:
  - id: <ADR-001>
    title: <decision title>
    decision: <what was decided>
    rationale: <why>
    alternatives:
      - <option considered>
    trade_offs:
      - <accepted trade-off>
```

Create `.claude/output/backlog.yaml`:

```yaml
project_name: <from problem.yaml>
version: "1.0"
phase: model

tasks:
  - id: <T-001>
    title: <short title>
    component: <which component>
    priority: <P0|P1|P2>
    type: <setup | feature | integration | test | infra | docs>
    description: <what to implement>
    depends_on:
      - <task id>
    acceptance_criteria:
      - <criterion>
    files_to_create:
      - <path>
    files_to_modify:
      - <path>
```

## Stage 5: Validate

Present an architecture summary and the backlog to the user. Ask for confirmation before saving. Highlight:

- Key architectural decisions and their trade-offs
- Total number of tasks by priority
- Any assumptions made

Once confirmed, save both files and report completion.

## Next Step

After the backlog is confirmed, tell the user to run `/ralph` to start autonomous implementation. Ralph will:
- Parse the backlog into a PRD with parallel rounds
- Create an agent team with teammates per story
- Coordinate implementation across rounds, respecting dependencies
- Run quality checks and produce `.claude/output/act-report.md`

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
