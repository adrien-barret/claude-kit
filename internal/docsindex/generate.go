package docsindex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/stack"
)

const (
	DocsIndexFile = "docs-index.md"
	DocsMetaFile  = ".docs-meta.json"
	StaleDays     = 14
)

// Meta stores metadata about the last docs-index generation.
type Meta struct {
	GeneratedAt    string   `json:"generated_at"`
	DependencyHash string   `json:"dependency_hash"`
	DetectedStack  []string `json:"detected_stack"`
}

// frameworkDirectives contains compressed stack-specific notes for each technology.
var frameworkDirectives = map[string]string{
	"nextjs": `## Next.js
- App Router: use app/ dir, layouts, loading.tsx, error.tsx, not-found.tsx
- Server Components by default; add "use client" only for interactivity
- Use next/image for images, next/link for navigation, next/font for fonts
- Data fetching: fetch() in Server Components, not getServerSideProps
- Route handlers: app/api/route.ts (GET, POST exports), not pages/api
- Metadata: export metadata object or generateMetadata() for SEO
- Middleware: middleware.ts at project root for auth, redirects
- Environment: NEXT_PUBLIC_ prefix for client-side env vars only`,

	"react": `## React
- Functional components + hooks, never class components
- useState for local state, useReducer for complex state
- useEffect: always specify deps array, clean up subscriptions
- useMemo/useCallback only for measurable performance issues, not by default
- Custom hooks for shared stateful logic, prefix with "use"
- Lifting state: move state to lowest common ancestor
- Controlled components for forms; ref only for focus/scroll/measure`,

	"vue": `## Vue
- Composition API with <script setup>, not Options API for new code
- ref() for primitives, reactive() for objects
- computed() for derived state, watch() sparingly
- defineProps/defineEmits with TypeScript generics for type safety
- Pinia for global state, composables for shared logic
- Teleport for modals/tooltips, Suspense for async components`,

	"angular": `## Angular
- Standalone components, not NgModules for new code
- Signals for reactive state, not zone.js patterns
- inject() function over constructor injection
- OnPush change detection on all components
- RxJS: use async pipe in templates, avoid manual subscribe
- Strong typing: no "any", use strict template checking`,

	"svelte": `## Svelte
- Runes ($state, $derived, $effect) for reactivity in Svelte 5+
- $props() for component props, not export let
- Snippets over slots for content projection
- Use SvelteKit for full apps: +page.svelte, +layout.svelte, +server.ts
- Form actions for mutations, load functions for data fetching`,

	"express": `## Express
- Router-level middleware for route groups, app-level for global
- Error middleware: 4 args (err, req, res, next), register last
- Validate request body/params/query at handler entry
- Use helmet for security headers, cors for CORS
- Async handlers: wrap with try/catch or express-async-errors
- Never send stack traces in production error responses`,

	"fastify": `## Fastify
- Schema-based validation with JSON Schema or TypeBox
- Plugins for encapsulation: register() with prefix
- Decorators for shared utilities, not global state
- Use @fastify/autoload for route auto-discovery
- Serialization schemas for response type safety and speed`,

	"nestjs": `## NestJS
- Modules for feature boundaries, providers for services
- DTOs with class-validator for request validation
- Guards for auth, Interceptors for transform/logging, Pipes for validation
- Repository pattern for data access, never query in controllers
- ConfigModule with validation for environment variables`,

	"django": `## Django
- Class-based views for CRUD, function views for custom logic
- Models: use migrations, never modify DB directly
- Forms/serializers for validation, never trust request.data raw
- Templates: use template tags, avoid logic in templates
- Settings: split base/dev/prod, use django-environ for env vars
- ORM: select_related/prefetch_related to avoid N+1 queries`,

	"flask": `## Flask
- Application factory pattern with create_app()
- Blueprints for feature modules, not everything in app.py
- Flask-SQLAlchemy for ORM, Flask-Migrate for migrations
- Validate with marshmallow or pydantic, not manual checks
- Error handlers: @app.errorhandler for consistent error format`,

	"fastapi": `## FastAPI
- Pydantic models for request/response validation
- Dependency injection for auth, DB sessions, config
- Background tasks for non-blocking operations
- Router prefixes for API versioning
- async def for I/O-bound handlers, def for CPU-bound
- Settings with pydantic-settings for typed env config`,

	"rails": `## Rails
- Convention over configuration: follow Rails naming/structure
- Strong Parameters: require/permit in controllers
- ActiveRecord: scopes, validations, callbacks sparingly
- Service objects for complex business logic
- Concerns for shared model/controller behavior
- N+1: use includes/preload/eager_load`,

	"typescript": `## TypeScript
- Strict mode always: strict: true in tsconfig.json
- Prefer type over interface for unions and intersections
- Use satisfies for type-safe object literals
- Discriminated unions over type assertions
- Branded types for domain IDs (UserId, OrderId)
- Template literal types for string patterns
- No any — use unknown for truly unknown types, then narrow`,

	"go": `## Go
- Accept interfaces, return structs
- Error wrapping: fmt.Errorf("context: %w", err)
- Table-driven tests, testify for assertions
- Context propagation: first param, never store in struct
- Goroutines: always handle cleanup (defer, context cancellation)
- Struct embedding for composition, not inheritance
- io.Reader/io.Writer for streaming, not []byte`,

	"python": `## Python
- Type hints on all function signatures, use mypy/pyright
- Dataclasses or Pydantic for data structures, not raw dicts
- Virtual env per project (venv, poetry, uv)
- f-strings for formatting, pathlib for file paths
- Context managers for resource cleanup
- List/dict comprehensions over map/filter where clearer`,

	"terraform": `## Terraform
- Modules for reusable infrastructure components
- Remote state with locking (S3+DynamoDB, GCS, etc.)
- Variables with descriptions and validation blocks
- Outputs for cross-module references
- lifecycle { prevent_destroy } on stateful resources
- Use moved blocks for refactoring, not manual state surgery`,

	"docker": `## Docker
- Multi-stage builds to minimize image size
- Run as non-root user (USER directive)
- COPY specific files, not . (use .dockerignore)
- Pin base image versions, not :latest
- One process per container
- HEALTHCHECK for orchestrator integration`,

	"kubernetes": `## Kubernetes
- Resource requests AND limits on all containers
- Liveness + readiness probes, distinct endpoints
- ConfigMap/Secret for config, not baked into images
- NetworkPolicy to restrict pod-to-pod traffic
- PodDisruptionBudget for availability during updates
- Labels: app, version, component, managed-by`,

	"prisma": `## Prisma
- Schema-first: define models in schema.prisma
- Migrations: npx prisma migrate dev, never push to prod
- Use select/include to fetch only needed fields
- Transactions: prisma.$transaction for multi-step operations
- Middleware for logging, soft deletes`,

	"tailwind": `## Tailwind CSS
- Use @apply sparingly — prefer utility classes in markup
- Design tokens via theme.extend in tailwind.config
- Responsive: mobile-first (sm:, md:, lg: breakpoints)
- Dark mode: class strategy for manual toggle support
- Purge unused styles: content paths in config`,

	"github-actions": `## GitHub Actions
- Pin action versions to SHA, not @main or @v1
- Use GITHUB_TOKEN, not PATs, where possible
- Cache dependencies (actions/cache) for faster runs
- Matrix strategy for cross-platform/version testing
- Concurrency groups to cancel redundant runs`,
}

// Generate creates docs-index.md and .docs-meta.json in the project root.
func Generate(projectRoot string) ([]string, error) {
	techs := stack.DetectStack(projectRoot)
	depHash := stack.ComputeDependencyHash(projectRoot)

	var techNames []string
	for _, t := range techs {
		techNames = append(techNames, t.Name)
	}

	// Build the docs-index.md content
	var sb strings.Builder
	sb.WriteString("# Docs Index\n\n")
	sb.WriteString("<!-- Auto-generated by bmad docs. Do not edit manually. -->\n")
	sb.WriteString("<!-- Run `bmad docs --refresh` to regenerate. -->\n\n")

	if len(techs) == 0 {
		sb.WriteString("No stack detected. Add dependency files (package.json, go.mod, etc.) and re-run.\n")
	} else {
		sb.WriteString(fmt.Sprintf("**Detected stack:** %s\n\n", strings.Join(techNames, ", ")))

		// Group by category
		categories := map[string][]Tech{
			"language":  {},
			"framework": {},
			"runtime":   {},
			"tool":      {},
		}
		for _, t := range techs {
			categories[t.Category] = append(categories[t.Category], t)
		}

		// Summary table
		sb.WriteString("| Category | Technologies |\n")
		sb.WriteString("|----------|-------------|\n")
		for _, cat := range []string{"language", "runtime", "framework", "tool"} {
			if len(categories[cat]) > 0 {
				names := make([]string, len(categories[cat]))
				for i, t := range categories[cat] {
					names[i] = t.Name
				}
				sb.WriteString(fmt.Sprintf("| %s | %s |\n", cat, strings.Join(names, ", ")))
			}
		}
		sb.WriteString("\n---\n\n")

		// Framework-specific directives
		for _, t := range techs {
			if directive, ok := frameworkDirectives[t.Name]; ok {
				sb.WriteString(directive)
				sb.WriteString("\n\n")
			}
		}
	}

	// Write docs-index.md
	claudeDir := filepath.Join(projectRoot, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating .claude directory: %w", err)
	}

	indexPath := filepath.Join(claudeDir, DocsIndexFile)
	if err := os.WriteFile(indexPath, []byte(sb.String()), 0o644); err != nil {
		return nil, fmt.Errorf("writing docs-index.md: %w", err)
	}

	// Write .docs-meta.json
	meta := Meta{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		DependencyHash: depHash,
		DetectedStack:  techNames,
	}
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling meta: %w", err)
	}

	metaPath := filepath.Join(claudeDir, DocsMetaFile)
	if err := os.WriteFile(metaPath, metaData, 0o644); err != nil {
		return nil, fmt.Errorf("writing .docs-meta.json: %w", err)
	}

	return techNames, nil
}

// IsStale checks if the docs-index needs regeneration.
// Returns true if:
// - .docs-meta.json doesn't exist
// - Dependency hash has changed
// - Last generation was more than StaleDays ago
func IsStale(projectRoot string) (bool, string) {
	metaPath := filepath.Join(projectRoot, ".claude", DocsMetaFile)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return true, "docs-index not yet generated"
	}

	var meta Meta
	if err := json.Unmarshal(data, &meta); err != nil {
		return true, "corrupted .docs-meta.json"
	}

	// Check dependency hash
	currentHash := stack.ComputeDependencyHash(projectRoot)
	if currentHash != meta.DependencyHash {
		return true, "dependency files have changed"
	}

	// Check age
	genTime, err := time.Parse(time.RFC3339, meta.GeneratedAt)
	if err != nil {
		return true, "invalid timestamp in .docs-meta.json"
	}

	age := time.Since(genTime)
	if age > time.Duration(StaleDays)*24*time.Hour {
		return true, fmt.Sprintf("docs-index is %d days old (threshold: %d days)", int(age.Hours()/24), StaleDays)
	}

	return false, ""
}

// Type alias for external use
type Tech = stack.Tech
