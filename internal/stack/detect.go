package stack

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Tech represents a detected technology.
type Tech struct {
	Name     string
	Category string // "language", "framework", "runtime", "database", "cloud", "tool"
}

// depFileMap maps dependency files to the technologies they indicate.
var depFileMap = map[string][]Tech{
	"package.json": {
		{Name: "node", Category: "runtime"},
		{Name: "javascript", Category: "language"},
	},
	"tsconfig.json": {
		{Name: "typescript", Category: "language"},
	},
	"go.mod": {
		{Name: "go", Category: "language"},
	},
	"requirements.txt": {
		{Name: "python", Category: "language"},
	},
	"pyproject.toml": {
		{Name: "python", Category: "language"},
	},
	"Pipfile": {
		{Name: "python", Category: "language"},
	},
	"Gemfile": {
		{Name: "ruby", Category: "language"},
	},
	"Cargo.toml": {
		{Name: "rust", Category: "language"},
	},
	"pom.xml": {
		{Name: "java", Category: "language"},
	},
	"build.gradle": {
		{Name: "java", Category: "language"},
	},
	"build.gradle.kts": {
		{Name: "kotlin", Category: "language"},
	},
	"composer.json": {
		{Name: "php", Category: "language"},
	},
	"Dockerfile": {
		{Name: "docker", Category: "tool"},
	},
	"docker-compose.yml": {
		{Name: "docker-compose", Category: "tool"},
	},
	"docker-compose.yaml": {
		{Name: "docker-compose", Category: "tool"},
	},
	".terraform.lock.hcl": {
		{Name: "terraform", Category: "tool"},
	},
	"Makefile": {
		{Name: "make", Category: "tool"},
	},
}

// frameworkDetectors check for framework-specific markers in package.json or similar.
type frameworkDetector struct {
	File      string
	Contains  string
	Framework Tech
}

var frameworkDetectors = []frameworkDetector{
	// JavaScript/TypeScript frameworks (check package.json content)
	{"package.json", `"next"`, Tech{Name: "nextjs", Category: "framework"}},
	{"package.json", `"react"`, Tech{Name: "react", Category: "framework"}},
	{"package.json", `"vue"`, Tech{Name: "vue", Category: "framework"}},
	{"package.json", `"nuxt"`, Tech{Name: "nuxt", Category: "framework"}},
	{"package.json", `"svelte"`, Tech{Name: "svelte", Category: "framework"}},
	{"package.json", `"@angular/core"`, Tech{Name: "angular", Category: "framework"}},
	{"package.json", `"express"`, Tech{Name: "express", Category: "framework"}},
	{"package.json", `"fastify"`, Tech{Name: "fastify", Category: "framework"}},
	{"package.json", `"hono"`, Tech{Name: "hono", Category: "framework"}},
	{"package.json", `"nest"`, Tech{Name: "nestjs", Category: "framework"}},
	{"package.json", `"@nestjs/core"`, Tech{Name: "nestjs", Category: "framework"}},
	{"package.json", `"tailwindcss"`, Tech{Name: "tailwind", Category: "framework"}},
	{"package.json", `"prisma"`, Tech{Name: "prisma", Category: "tool"}},
	{"package.json", `"drizzle-orm"`, Tech{Name: "drizzle", Category: "tool"}},
	// Python frameworks (check requirements.txt or pyproject.toml)
	{"requirements.txt", "django", Tech{Name: "django", Category: "framework"}},
	{"requirements.txt", "flask", Tech{Name: "flask", Category: "framework"}},
	{"requirements.txt", "fastapi", Tech{Name: "fastapi", Category: "framework"}},
	{"pyproject.toml", "django", Tech{Name: "django", Category: "framework"}},
	{"pyproject.toml", "flask", Tech{Name: "flask", Category: "framework"}},
	{"pyproject.toml", "fastapi", Tech{Name: "fastapi", Category: "framework"}},
	// Ruby frameworks
	{"Gemfile", "rails", Tech{Name: "rails", Category: "framework"}},
	{"Gemfile", "sinatra", Tech{Name: "sinatra", Category: "framework"}},
	// PHP frameworks
	{"composer.json", "laravel", Tech{Name: "laravel", Category: "framework"}},
	{"composer.json", "symfony", Tech{Name: "symfony", Category: "framework"}},
}

// directoryDetectors check for the presence of directories.
var directoryDetectors = map[string]Tech{
	".github":    {Name: "github-actions", Category: "tool"},
	".gitlab-ci": {Name: "gitlab-ci", Category: "tool"},
	"k8s":        {Name: "kubernetes", Category: "tool"},
	"helm":       {Name: "helm", Category: "tool"},
	"infra":      {Name: "infrastructure", Category: "tool"},
}

// DetectStack scans a project root and returns detected technologies.
func DetectStack(projectRoot string) []Tech {
	seen := make(map[string]bool)
	var techs []Tech

	addTech := func(t Tech) {
		if !seen[t.Name] {
			seen[t.Name] = true
			techs = append(techs, t)
		}
	}

	// Check for dependency files
	for file, fileTechs := range depFileMap {
		if _, err := os.Stat(filepath.Join(projectRoot, file)); err == nil {
			for _, t := range fileTechs {
				addTech(t)
			}
		}
	}

	// Check for .tf files (Terraform)
	if matches, _ := filepath.Glob(filepath.Join(projectRoot, "*.tf")); len(matches) > 0 {
		addTech(Tech{Name: "terraform", Category: "tool"})
	}
	if matches, _ := filepath.Glob(filepath.Join(projectRoot, "infra", "*.tf")); len(matches) > 0 {
		addTech(Tech{Name: "terraform", Category: "tool"})
	}

	// Check for framework-specific markers
	for _, d := range frameworkDetectors {
		filePath := filepath.Join(projectRoot, d.File)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		if strings.Contains(string(content), d.Contains) {
			addTech(d.Framework)
		}
	}

	// Check for directories
	for dir, t := range directoryDetectors {
		if info, err := os.Stat(filepath.Join(projectRoot, dir)); err == nil && info.IsDir() {
			addTech(t)
		}
	}

	// Check for .github/workflows specifically
	if _, err := os.Stat(filepath.Join(projectRoot, ".github", "workflows")); err == nil {
		addTech(Tech{Name: "github-actions", Category: "tool"})
	}

	sort.Slice(techs, func(i, j int) bool {
		return techs[i].Name < techs[j].Name
	})

	return techs
}

// ListDependencyFiles returns the list of dependency files present in the project.
func ListDependencyFiles(projectRoot string) []string {
	var files []string
	for file := range depFileMap {
		if _, err := os.Stat(filepath.Join(projectRoot, file)); err == nil {
			files = append(files, file)
		}
	}
	sort.Strings(files)
	return files
}

// ComputeDependencyHash computes a SHA-256 hash of all dependency file contents.
// Used to detect when dependencies have changed and docs-index should be refreshed.
func ComputeDependencyHash(projectRoot string) string {
	files := ListDependencyFiles(projectRoot)
	if len(files) == 0 {
		return ""
	}

	h := sha256.New()
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(projectRoot, f))
		if err != nil {
			continue
		}
		h.Write([]byte(f))
		h.Write(data)
	}

	return hex.EncodeToString(h.Sum(nil))[:16]
}
