---
name: bmad-model
description: BMAD Model phase – design architecture, produce ADRs, and generate a prioritized backlog
---

Act as an **Architect** and **Tech Lead** working together.

Your goal is to take the problem definition from the Break phase and produce an architecture design and implementation backlog.

## Prerequisites

Read `.claude/output/problem.yaml`. If it does not exist, tell the user to run `/bmad-break` first and stop.

Read `.claude/output/principles.md` if it exists — use it to inform architectural decisions, tech stack choices, and acceptance criteria in the backlog.

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

After the backlog is confirmed, the recommended next steps are:

1. Run `/analyze` to check for cross-artifact inconsistencies
2. Run `/checklist` to verify implementation readiness
3. Run `/gsd-prep` to generate codebase mapping and context packs
4. Run `/ralph` (or `/bmad-act`) to begin implementation

Or simply run `/bmad-run` which orchestrates all of this automatically.

If $ARGUMENTS is provided, use it as additional context: $ARGUMENTS
