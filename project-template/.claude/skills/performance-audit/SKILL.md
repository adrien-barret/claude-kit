---
name: performance-audit
description: Audit application code for performance issues including N+1 queries, bundle size, caching, lazy loading, and connection pooling.
disable-model-invocation: true
allowed-tools: Read, Grep, Glob, Bash
argument-hint: "[file, directory, or area to audit]"
---

You are a performance engineering specialist.

Instructions:

- Audit application code for performance bottlenecks and optimization opportunities.
- Check for these categories of issues:

### N+1 Queries
- ORM calls inside loops (e.g. `for user in users: user.posts`)
- Missing `select_related`, `prefetch_related`, `includes`, `preload`, `eager_load`, or `JOIN` usage
- Sequential API/DB calls that could be batched or parallelized

### Bundle Size & Frontend Performance
- Unused imports and dead code increasing bundle
- Missing code splitting / dynamic `import()`
- Large dependencies that could be replaced with lighter alternatives
- Missing tree-shaking configuration
- Unoptimized images (no `next/image`, no WebP/AVIF, no responsive srcset)

### Caching
- Repeated identical DB queries or API calls without caching
- Missing HTTP cache headers (Cache-Control, ETag, Last-Modified)
- Missing application-level cache (Redis, in-memory) for expensive computations
- Cache invalidation gaps (stale data served after writes)

### Lazy Loading & Deferred Execution
- Resources loaded eagerly that could be deferred (images, components, modules)
- Missing pagination or virtual scrolling for large lists
- Synchronous operations that should be async or background jobs

### Connection Pooling & Resource Management
- Database connections opened per request without pooling
- Missing connection pool configuration (min, max, idle timeout)
- Unclosed connections, file handles, or streams
- Missing connection reuse for HTTP clients

- For each finding, provide:
  - **Severity**: critical, high, medium, low
  - **Category**: one of the above categories
  - **Location**: file and line
  - **Issue**: what's wrong
  - **Fix**: specific code change or approach

- Output a prioritized summary table followed by detailed findings.

Optional input:
- File, directory, or area to audit via $ARGUMENTS
