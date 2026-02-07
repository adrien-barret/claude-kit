package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

// Recommendation represents a component suggested by Claude.
type Recommendation struct {
	Source      string `json:"source"`                // "local", "voltagent", "aitmpl"
	Type       string `json:"type"`                  // "agents", "skills", "commands", "rules"
	Name       string `json:"name"`                  // component name
	Description string `json:"description"`           // what it does
	URL        string `json:"url,omitempty"`          // source URL for external components
}

func runSmartAdd(tmplDir, targetDir, query string) error {
	fmt.Println(banner())
	fmt.Println(subtitleStyle.Render("  Smart Add — AI-powered component discovery"))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Query: %s", query)))
	fmt.Println()

	// Check claude CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH; smart add requires the Claude Code CLI")
	}

	// Build local catalog
	localCatalog := buildLocalCatalog(tmplDir)

	// Fetch external catalogs with spinner
	var voltAgentCatalog string
	var fetchErr error
	_ = spinner.New().
		Title("Fetching external catalogs...").
		Action(func() {
			voltAgentCatalog, fetchErr = fetchVoltAgent()
		}).
		Run()

	if fetchErr != nil {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  Could not fetch VoltAgent catalog: %v", fetchErr)))
		voltAgentCatalog = ""
	}

	// Build prompt and run Claude
	prompt := buildSmartAddPrompt(query, localCatalog, voltAgentCatalog)

	var recommendations []Recommendation
	var claudeErr error
	_ = spinner.New().
		Title("Asking Claude for recommendations...").
		Action(func() {
			recommendations, claudeErr = runClaudeRecommend(prompt)
		}).
		Run()

	if claudeErr != nil {
		return fmt.Errorf("recommendation failed: %w", claudeErr)
	}

	if len(recommendations) == 0 {
		fmt.Println(warnStyle.Render("  No matching components found for your query."))
		return nil
	}

	// Show recommendations and let user pick
	return presentRecommendations(tmplDir, targetDir, recommendations)
}

// buildLocalCatalog produces a text summary of all local template components.
func buildLocalCatalog(tmplDir string) string {
	categories, err := catalog.ScanTemplate(tmplDir)
	if err != nil {
		return "(error scanning local templates)"
	}

	var sb strings.Builder
	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf("\n### %s\n", cat.Name))
		for _, c := range cat.Components {
			desc := c.Description
			if desc == "" {
				desc = "(no description)"
			}
			sb.WriteString(fmt.Sprintf("- %s: %s\n", c.Name, desc))
		}
	}
	return sb.String()
}

// fetchVoltAgent downloads the VoltAgent awesome-agent-skills README.
func fetchVoltAgent() (string, error) {
	url := "https://raw.githubusercontent.com/VoltAgent/awesome-agent-skills/main/README.md"
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	content := string(body)
	// Truncate if too large to fit in prompt
	if len(content) > 20000 {
		content = content[:20000] + "\n... (truncated)"
	}
	return content, nil
}

// buildSmartAddPrompt creates the prompt sent to Claude for component recommendation.
func buildSmartAddPrompt(query, localCatalog, voltAgent string) string {
	var sb strings.Builder

	sb.WriteString(`You are a Claude Code component recommender for the "claude-kit" CLI tool.

Based on the user's description, recommend components to install from the available sources.

## USER REQUEST
`)
	sb.WriteString(query)
	sb.WriteString(`

## SOURCE 1: LOCAL TEMPLATES (can be installed directly)
`)
	sb.WriteString(localCatalog)

	if voltAgent != "" {
		sb.WriteString(`

## SOURCE 2: VOLTAGENT (external GitHub skills)
Each entry below links to a GitHub repository containing Claude Code skills.
`)
		sb.WriteString(voltAgent)
	}

	sb.WriteString(`

## SOURCE 3: AITMPL.COM
Browse https://www.aitmpl.com for additional Claude Code components (skills, agents, commands, settings, hooks, MCPs).
Recommend components from this source based on your knowledge of what is available there.

## INSTRUCTIONS

1. Analyze the user's request
2. Find matching components from ALL sources (local first, then external)
3. Prefer local components when a good match exists
4. For VoltAgent entries, include the GitHub repository URL from the README
5. For aitmpl.com entries, include "https://www.aitmpl.com" as the URL
6. type must be one of: agents, skills, commands, rules

Respond with ONLY a valid JSON array — no markdown fences, no explanation, no text before or after:
[
  {
    "source": "local",
    "type": "skills",
    "name": "code-reviewer",
    "description": "Reviews code for quality and best practices",
    "url": ""
  },
  {
    "source": "voltagent",
    "type": "skills",
    "name": "some-skill",
    "description": "What it does",
    "url": "https://github.com/user/repo"
  }
]

If nothing matches, respond with: []
`)

	return sb.String()
}

// runClaudeRecommend invokes the Claude CLI non-interactively and parses the JSON response.
func runClaudeRecommend(prompt string) ([]Recommendation, error) {
	cmd := exec.Command("claude", "-p", prompt, "--output-format", "text")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("claude CLI exited with code %d: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("claude CLI: %w", err)
	}

	responseStr := strings.TrimSpace(string(output))

	// Extract JSON array from response (Claude might wrap in markdown)
	jsonStr := extractJSONArray(responseStr)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON array found in Claude response:\n%s", responseStr)
	}

	var recs []Recommendation
	if err := json.Unmarshal([]byte(jsonStr), &recs); err != nil {
		return nil, fmt.Errorf("invalid JSON in Claude response: %w\nExtracted: %s", err, jsonStr)
	}

	return recs, nil
}

// extractJSONArray finds the first JSON array in a string.
func extractJSONArray(s string) string {
	start := strings.Index(s, "[")
	if start < 0 {
		return ""
	}
	// Find matching closing bracket
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

// presentRecommendations shows a multi-select form and installs chosen components.
func presentRecommendations(tmplDir, targetDir string, recs []Recommendation) error {
	fmt.Println(sectionHeader("Recommendations"))

	options := make([]huh.Option[int], 0, len(recs))
	for i, rec := range recs {
		sourceTag := sourceLabel(rec.Source)
		desc := rec.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		label := fmt.Sprintf("%s %s/%s — %s", sourceTag, rec.Type, rec.Name, desc)
		options = append(options, huh.NewOption(label, i))
	}

	var selected []int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[int]().
				Title("Select components to install").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if len(selected) == 0 {
		fmt.Println("No components selected.")
		return nil
	}

	ensureBaseFiles(tmplDir, targetDir)

	for _, idx := range selected {
		rec := recs[idx]
		switch rec.Source {
		case "local":
			installLocalRec(tmplDir, targetDir, rec)
		default:
			installExternalRec(targetDir, rec)
		}
	}

	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Done!", arrow)))
	return nil
}

func sourceLabel(source string) string {
	switch source {
	case "local":
		return "[local]"
	case "voltagent":
		return "[VoltAgent]"
	case "aitmpl":
		return "[aitmpl]"
	default:
		return fmt.Sprintf("[%s]", source)
	}
}

// installLocalRec installs a component from the local template catalog.
func installLocalRec(tmplDir, targetDir string, rec Recommendation) {
	if rec.Type == "agents" {
		installAgent(tmplDir, targetDir, rec.Name)
		return
	}

	if catalog.IsInstalled(targetDir, rec.Type, rec.Name) {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  %s/%s already installed, updating", rec.Type, rec.Name)))
	}

	if err := catalog.CopyComponent(tmplDir, targetDir, rec.Type, rec.Name); err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  %s/%s: %v", rec.Type, rec.Name, err)))
		return
	}
	fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Added %s/%s", rec.Type, rec.Name))))
}

// installExternalRec installs a component from an external source.
func installExternalRec(targetDir string, rec Recommendation) {
	if rec.URL == "" {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  %s/%s [%s]: no URL provided, skipping", rec.Type, rec.Name, rec.Source)))
		return
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("  %s Fetching %s/%s from %s...", bullet, rec.Type, rec.Name, rec.Source)))

	if strings.Contains(rec.URL, "github.com") {
		if err := installFromGitHub(targetDir, rec); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  %s/%s: %v", rec.Type, rec.Name, err)))
			return
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Added %s/%s [%s]", rec.Type, rec.Name, rec.Source))))
		return
	}

	// Non-GitHub external — show URL for manual install
	fmt.Println(warnStyle.Render(fmt.Sprintf("  %s/%s: install manually from %s", rec.Type, rec.Name, rec.URL)))
}

// installFromGitHub clones a GitHub repo as a skill directory, or fetches a raw file.
func installFromGitHub(targetDir string, rec Recommendation) error {
	// Skills are directories — clone the repo
	if rec.Type == "skills" {
		destDir := filepath.Join(targetDir, "skills", rec.Name)
		if err := os.MkdirAll(filepath.Dir(destDir), 0o755); err != nil {
			return err
		}

		cmd := exec.Command("git", "clone", "--depth=1", "--quiet", rec.URL, destDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}

		// Remove .git directory from cloned skill
		_ = os.RemoveAll(filepath.Join(destDir, ".git"))
		return nil
	}

	// Agents, commands, rules are single markdown files — try to fetch raw content
	rawURL := githubToRaw(rec.URL)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching %s: status %d", rawURL, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	destFile := filepath.Join(targetDir, rec.Type, rec.Name+".md")
	if err := os.MkdirAll(filepath.Dir(destFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destFile, content, 0o644)
}

// githubToRaw converts a GitHub blob/tree URL to a raw.githubusercontent.com URL.
func githubToRaw(url string) string {
	url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
	url = strings.Replace(url, "/blob/", "/", 1)
	url = strings.Replace(url, "/tree/", "/", 1)
	return url
}

