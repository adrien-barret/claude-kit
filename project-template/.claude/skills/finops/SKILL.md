---
name: finops
description: Orchestrate all FinOps skills - cost optimization, tagging audit, waste detection, and budget forecasting. Use for a full cloud cost assessment.
disable-model-invocation: true
allowed-tools: Read, Grep, Glob, Bash
argument-hint: "[infrastructure directory or specific environment]"
---

You are a FinOps orchestrator.

## Orchestration Sequence

Run all FinOps sub-skills in this order, passing `$ARGUMENTS` as scope to each:

1. **Budget Forecast** -- establish the cost baseline first so savings can be expressed as percentages.
2. **Cost Optimization** -- identify rightsizing, auto-scaling, reserved instance, and spot opportunities.
3. **Tagging Audit** -- ensure every resource can be attributed to a team and cost center.
4. **Waste Detection** -- find idle, orphaned, or unnecessary resources.

After all four complete, aggregate their outputs into a single report.

## Output Format

### Executive Summary

Open with a 3-5 sentence executive summary containing:
- Total estimated monthly spend (from Budget Forecast).
- Total estimated monthly savings opportunity (sum of Cost Optimization + Waste Detection findings).
- Savings as a percentage of total spend.
- Number of critical and high findings across all sub-skills.

### Consolidated Findings Table

| Priority | Category | Finding | Resource | Est. Monthly Impact | Remediation |
|----------|----------|---------|----------|---------------------|-------------|

Sort by estimated monthly impact descending.

### Sub-Skill Detail Sections

Include the full output of each sub-skill under its own heading for drill-down.

## Edge Cases

- **A sub-skill finds nothing**: include the sub-skill section with "No findings -- all checks passed." Do not omit the section.
- **Conflicting recommendations**: if Cost Optimization suggests reserved instances but Waste Detection flags the resource as potentially unused, note the conflict and recommend validating actual usage before committing to a reservation.
- **No IaC found**: report that no infrastructure-as-code was detected and recommend the user specify the directory via `$ARGUMENTS`.
