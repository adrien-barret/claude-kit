---
name: cost-review
description: Run a full FinOps review covering cost optimization, tagging, waste detection, and budget forecasting
---

You are a FinOps review orchestrator. Invoke the finops agent from `.claude/agents/finops.md` and run all four sub-skills in sequence.

## Scope

If $ARGUMENTS is provided, scope the review to the specified directories, services, or cloud providers (e.g., `aws`, `gcp`, `terraform/`, `k8s/`). Otherwise, scan the entire project.

## Workflow

### 1. Cost Optimization
Skill: `.claude/skills/finops/cost-optimization/SKILL.md`
- Identify over-provisioned resources, missing reservations, and expensive configurations

### 2. Tagging Audit
Skill: `.claude/skills/finops/tagging-audit/SKILL.md`
- Verify all resources have required tags: `environment`, `team`, `service`, `cost-center`

### 3. Waste Detection
Skill: `.claude/skills/finops/waste-detection/SKILL.md`
- Identify unused resources: unattached volumes, idle load balancers, orphaned snapshots

### 4. Budget Forecast
Skill: `.claude/skills/finops/budget-forecast/SKILL.md`
- Estimate monthly cost based on current resource definitions

## Output Format

```
## FinOps Review Summary

| # | Severity | Skill Area        | Finding                          | Est. Monthly Impact | Remediation                   |
|---|----------|-------------------|----------------------------------|---------------------|-------------------------------|
| 1 | Critical | waste-detection   | 3 unattached EBS volumes (500GB) | ~$50                | Delete or snapshot and remove |
| 2 | High     | cost-optimization | RDS instance is db.r5.2xlarge    | ~$200 savings       | Downsize to db.r5.large       |

## Total Estimated Monthly Savings: $X
```

Sort by severity (Critical > High > Medium > Low), then by estimated impact descending.
