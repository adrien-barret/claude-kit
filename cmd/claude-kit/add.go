package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var addCmd = &cobra.Command{
	Use:   "add [names...] | add <type> <name> | add new <description>",
	Short: "Add agents, components, or discover new ones with AI",
	Long: `Add components to the current project's .claude/ directory.

When called with no arguments, shows an interactive agent picker.
Selected agents automatically install their linked skills and rules.

When called with agent names, installs those agents + their dependencies.
When called with a type prefix, installs that specific component type.

Use "new" to trigger Smart Add: searches local templates, VoltAgent,
and aitmpl.com using Claude CLI, then lets you pick and install.

Examples:
  ck add                                  # Interactive agent picker
  ck add bmad                             # Install all BMAD agents, commands, and rules
  ck add backend devops                   # Add backend + devops agents + deps
  ck add skill code-reviewer              # Add a specific skill
  ck add command review                   # Add a specific command
  ck add rule testing                     # Add a specific rule
  ck add new database review              # Smart add — AI finds matching components
  ck add new performance auditing         # Smart add — natural language query`,
	RunE: runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	tmplDir := resolveTemplateDir()
	targetDir := resolveTarget()

	// No args → interactive agent picker
	if len(args) == 0 {
		return runInteractiveAdd(tmplDir, targetDir)
	}

	// "bmad" bundle → install all BMAD agents + commands + rules
	if len(args) == 1 && strings.ToLower(args[0]) == "bmad" {
		return addBmadBundle(tmplDir, targetDir)
	}

	// "new" keyword → smart add: ck add new <description>
	if strings.ToLower(args[0]) == "new" {
		if len(args) < 2 {
			return fmt.Errorf("usage: ck add new <description>\n  Example: ck add new database review")
		}
		query := strings.Join(args[1:], " ")
		return runSmartAdd(tmplDir, targetDir, query)
	}

	// Check if first arg is a type prefix
	firstNorm := normalizeType(args[0])
	if isComponentType(firstNorm) && len(args) >= 2 {
		// Explicit type: ck add skill code-reviewer
		return addExplicitComponents(tmplDir, targetDir, firstNorm, args[1:])
	}

	// Otherwise treat all args as agent names: ck add backend devops sre
	return addAgents(tmplDir, targetDir, args)
}

func isComponentType(t string) bool {
	switch t {
	case "agents", "skills", "commands", "rules":
		return true
	}
	return false
}

// runInteractiveAdd shows a multi-select of available agents.
func runInteractiveAdd(tmplDir, targetDir string) error {
	fmt.Println(banner())
	fmt.Println(subtitleStyle.Render("  Add agents (skills & rules are installed automatically)"))
	fmt.Println()

	categories, err := catalog.ScanTemplate(tmplDir)
	if err != nil {
		return fmt.Errorf("scanning templates: %w", err)
	}

	// Find the agents category
	var agentComps []catalog.Component
	for _, cat := range categories {
		if cat.Name == "agents" {
			agentComps = cat.Components
			break
		}
	}

	if len(agentComps) == 0 {
		return fmt.Errorf("no agents found in template directory")
	}

	// Build options — only show agents not yet installed
	options := make([]huh.Option[string], 0, len(agentComps))
	for _, c := range agentComps {
		if catalog.IsInstalled(targetDir, "agents", c.Name) {
			continue
		}
		label := c.Name
		if c.Description != "" {
			desc := c.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			label = fmt.Sprintf("%s -- %s", c.Name, desc)
		}
		options = append(options, huh.NewOption(label, c.Name))
	}

	if len(options) == 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  %s All agents already installed!", arrow)))
		fmt.Println(dimStyle.Render("  Use 'ck remove' to remove agents."))
		return nil
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select agents to add").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if len(selected) == 0 {
		fmt.Println("No agents selected.")
		return nil
	}

	// Ensure base files exist
	ensureBaseFiles(tmplDir, targetDir)

	for _, name := range selected {
		installAgent(tmplDir, targetDir, name)
	}

	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Done!", arrow)))
	return nil
}

// addAgents installs agents by name with their dependencies.
func addAgents(tmplDir, targetDir string, names []string) error {
	fmt.Println(banner())

	ensureBaseFiles(tmplDir, targetDir)

	for _, name := range names {
		installAgent(tmplDir, targetDir, name)
	}

	fmt.Println()
	return nil
}

// addExplicitComponents installs components of a specific type.
func addExplicitComponents(tmplDir, targetDir, compType string, names []string) error {
	fmt.Println(banner())

	ensureBaseFiles(tmplDir, targetDir)

	for _, name := range names {
		if catalog.IsInstalled(targetDir, compType, name) {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  %s/%s already installed, updating", compType, name)))
		}

		if err := catalog.CopyComponent(tmplDir, targetDir, compType, name); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  %s/%s: %v", compType, name, err)))
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Added %s/%s", compType, name))))

		// Auto-pull deps for agents
		if compType == "agents" {
			pullAgentDeps(tmplDir, targetDir, name)
		}
	}

	fmt.Println()
	return nil
}

// installAgent installs an agent + all its skill dependencies.
func installAgent(tmplDir, targetDir, name string) {
	// Install the agent itself
	if err := catalog.CopyComponent(tmplDir, targetDir, "agents", name); err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  Agent %s: %v", name, err)))
		return
	}
	fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Added agent: %s", name))))

	// Pull skill deps
	pullAgentDeps(tmplDir, targetDir, name)

	// Auto-install related rules based on agent role
	autoRules := guessRulesForAgent(name)
	if len(autoRules) > 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("    %s Installing %d related rules:", bullet, len(autoRules))))
		for _, rule := range autoRules {
			if catalog.IsInstalled(targetDir, "rules", rule) {
				continue
			}
			if err := catalog.CopyComponent(tmplDir, targetDir, "rules", rule); err != nil {
				continue // silently skip missing rules
			}
			fmt.Println(fmt.Sprintf("      %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added rule: %s", rule))))
		}
	}

	// Auto-install extra skills (not in frontmatter)
	extraSkills := guessExtraSkillsForAgent(name)
	if len(extraSkills) > 0 {
		for _, skill := range extraSkills {
			if catalog.IsInstalled(targetDir, "skills", skill) {
				continue
			}
			if err := catalog.CopyComponent(tmplDir, targetDir, "skills", skill); err != nil {
				continue
			}
			fmt.Println(fmt.Sprintf("      %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added extra skill: %s", skill))))
		}
	}

	// Auto-install related commands
	autoCmds := guessCommandsForAgent(name)
	if len(autoCmds) > 0 {
		for _, cmd := range autoCmds {
			if catalog.IsInstalled(targetDir, "commands", cmd) {
				continue
			}
			if err := catalog.CopyComponent(tmplDir, targetDir, "commands", cmd); err != nil {
				continue
			}
			fmt.Println(fmt.Sprintf("      %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added command: %s", cmd))))
		}
	}
}

// pullAgentDeps reads an agent's frontmatter and installs its skill dependencies.
func pullAgentDeps(tmplDir, targetDir, agentName string) {
	agentPath := filepath.Join(tmplDir, "agents", agentName+".md")
	deps := catalog.ExtractSkillDeps(agentPath)

	if len(deps) == 0 {
		return
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("    %s Installing %d skill dependencies:", bullet, len(deps))))

	for _, dep := range deps {
		if catalog.IsInstalled(targetDir, "skills", dep) {
			fmt.Println(dimStyle.Render(fmt.Sprintf("      %s %s (already installed)", dot, dep)))
			continue
		}
		if err := catalog.CopyComponent(tmplDir, targetDir, "skills", dep); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("      Skill %s: %v", dep, err)))
			continue
		}
		fmt.Println(fmt.Sprintf("      %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added skill: %s", dep))))
	}
}

// addBmadBundle installs the BMAD methodology: core agents, workflow commands, and base rules.
// Project-specific agents (backend, mobile, etc.) are added separately via ck add <agent>.
func addBmadBundle(tmplDir, targetDir string) error {
	fmt.Println(banner())
	fmt.Println(subtitleStyle.Render("  BMAD — Break, Model, Act, Deliver"))
	fmt.Println(dimStyle.Render("  Core methodology + workflow commands"))
	fmt.Println(dimStyle.Render("  Add project agents separately: ck add backend, ck add mobile-ios, etc."))
	fmt.Println()

	ensureBaseFiles(tmplDir, targetDir)

	// Core BMAD agents (methodology roles, always the same)
	agents := []string{"product-owner", "architect", "tech-lead"}

	fmt.Println(sectionHeader("Core Agents"))
	for _, name := range agents {
		installAgent(tmplDir, targetDir, name)
	}

	// BMAD workflow commands + utilities
	commands := []string{
		"bmad-run", "bmad-break", "bmad-model", "bmad-act", "bmad-deliver",
		"principles", "clarify", "analyze", "checklist",
		"ralph", "ralph-loop", "ralph-cancel",
		"r", "p", "c", "g",
		"gsd-prep",
		"role-product-owner", "role-architect", "role-tech-lead",
		"review", "test-gen", "security-check", "commit-msg",
		"code-only", "docs-gen", "pr-review",
	}

	fmt.Println(sectionHeader("Commands"))
	for _, name := range commands {
		if catalog.IsInstalled(targetDir, "commands", name) {
			continue
		}
		if err := catalog.CopyComponent(tmplDir, targetDir, "commands", name); err != nil {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added command: %s", name))))
	}

	// Base rules (universal regardless of project type)
	rules := []string{"code-style", "testing", "security", "documentation"}

	fmt.Println(sectionHeader("Rules"))
	for _, name := range rules {
		if catalog.IsInstalled(targetDir, "rules", name) {
			continue
		}
		if err := catalog.CopyComponent(tmplDir, targetDir, "rules", name); err != nil {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Added rule: %s", name))))
	}

	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("  %s BMAD methodology installed!", arrow)))
	fmt.Println(dimStyle.Render("  Now add your project agents: ck add backend, ck add frontend, etc."))
	return nil
}

// guessCommandsForAgent returns commands that should be auto-installed for a given agent.
func guessCommandsForAgent(name string) []string {
	devCmds := []string{"commit-msg", "review", "test-gen", "pr-review", "code-only"}
	roleCmd := "role-" + name

	switch strings.ToLower(name) {
	case "backend":
		return append(devCmds, roleCmd)
	case "frontend":
		return append(devCmds, roleCmd)
	case "mobile-react-native", "mobile-flutter", "mobile-ios", "mobile-android":
		return append(devCmds, roleCmd)
	case "tech-lead":
		return append(devCmds, "security-check", roleCmd)
	case "devops":
		return []string{"commit-msg", roleCmd}
	case "finops":
		return []string{"cost-review", roleCmd}
	case "security":
		return []string{"security-check", roleCmd}
	case "pentester":
		return []string{"pentest", roleCmd}
	case "architect":
		return []string{"docs-gen", roleCmd}
	case "product-owner":
		return []string{roleCmd}
	case "ui-designer":
		return []string{roleCmd}
	case "ux-designer":
		return []string{roleCmd}
	default:
		return []string{roleCmd}
	}
}

// guessExtraSkillsForAgent returns skills not in agent frontmatter but logically related.
func guessExtraSkillsForAgent(name string) []string {
	switch strings.ToLower(name) {
	case "backend", "frontend", "mobile-react-native", "mobile-flutter",
		"mobile-ios", "mobile-android", "tech-lead", "devops":
		return []string{"git-commit-helper"}
	case "finops":
		return []string{"finops"} // parent orchestrator
	case "security":
		return []string{"security"} // parent orchestrator
	case "architect":
		return []string{"terraform-review"}
	default:
		return nil
	}
}

// guessRulesForAgent returns rules that make sense for a given agent role.
func guessRulesForAgent(agentName string) []string {
	agentName = strings.ToLower(agentName)
	switch agentName {
	case "backend", "tech-lead":
		return []string{"code-style", "testing", "security", "api"}
	case "frontend":
		return []string{"code-style", "testing", "security", "frontend"}
	case "mobile-react-native", "mobile-flutter", "mobile-ios", "mobile-android":
		return []string{"code-style", "testing", "security"}
	case "ui-designer", "ux-designer":
		return []string{"frontend"}
	case "architect":
		return []string{"code-style", "security", "api", "infrastructure"}
	case "product-owner":
		return []string{"documentation"}
	case "devops":
		return []string{"infrastructure", "security", "documentation"}
	case "security", "pentester":
		return []string{"security"}
	case "finops":
		return []string{"finops", "infrastructure"}
	default:
		return []string{"code-style", "security"}
	}
}

// ensureBaseFiles copies CLAUDE.md + settings.json if they don't exist.
func ensureBaseFiles(tmplDir, targetDir string) {
	claudeMd := filepath.Join(targetDir, "CLAUDE.md")
	if _, err := os.Stat(claudeMd); os.IsNotExist(err) {
		_ = catalog.CopyBaseFiles(tmplDir, targetDir)
	}
}

func normalizeType(t string) string {
	t = strings.ToLower(t)
	switch t {
	case "agent":
		return "agents"
	case "skill":
		return "skills"
	case "command":
		return "commands"
	case "rule":
		return "rules"
	}
	return t
}
