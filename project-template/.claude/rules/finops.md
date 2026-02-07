---
description: FinOps rules applied when writing or reviewing infrastructure and cloud resource definitions
globs: ["infra/**", "terraform/**", "*.tf", "*.tfvars", "cloudformation/**", "*.yaml", "k8s/**", "helm/**", "docker-compose*"]
---

## Resource Tagging and Allocation

- Tag every cloud resource with at least: `environment`, `team`, `service`, `cost-center`
- Use cost allocation reports to attribute spend to teams and services; review monthly
- Set up budget alerts at 50%, 80%, and 100% thresholds; enable anomaly detection for unexpected spikes

## Compute Optimization

- Use auto-scaling with both min and max bounds instead of fixed instance counts for variable workloads
- Use spot/preemptible instances for fault-tolerant and batch workloads (60-90% savings)
- Right-size resources: do not default to large instance types without load-testing justification
- Schedule non-production environments to shut down outside business hours
- Review reserved instance and savings plan coverage before deploying new on-demand resources

## Storage and Data Transfer

- Set storage lifecycle policies: transition to cold/archive tiers after defined periods
- Set retention policies on logs, snapshots, and backups; delete what is no longer needed
- Avoid cross-region data transfer unless required by compliance or latency

## Service Selection

- Prefer managed services over self-hosted when total cost of ownership is lower
- Evaluate serverless options for sporadic workloads to avoid paying for idle compute
