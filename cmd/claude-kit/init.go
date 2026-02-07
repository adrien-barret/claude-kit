package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup — pick agents, everything else is automatic",
	Long: `Initialize a project with an interactive agent picker.

Pick the agents you need and Claude Kit automatically installs
all related skills, commands, and rules.

Only agents not yet installed are shown. Use 'ck remove' to
remove installed agents and their dependencies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractiveInit()
	},
}

func runInteractiveInit() error {
	tmplDir := resolveTemplateDir()
	targetDir := resolveTarget()

	fmt.Println(banner())

	// Scan available components
	categories, err := catalog.ScanTemplate(tmplDir)
	if err != nil {
		return fmt.Errorf("scanning templates: %w", err)
	}

	if len(categories) == 0 {
		return fmt.Errorf("no components found in template directory: %s", tmplDir)
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

	// Check for existing installation
	isExisting := false
	installedAgents := make(map[string]bool)
	if _, err := os.Stat(targetDir); err == nil {
		isExisting = true
		installed, _ := catalog.GetInstalled(targetDir)
		for _, cat := range installed {
			if cat.Name == "agents" {
				for _, c := range cat.Components {
					installedAgents[c.Name] = true
				}
			}
		}
	}

	if isExisting {
		fmt.Println(subtitleStyle.Render("  Project Setup (existing .claude/ detected)"))
		if len(installedAgents) > 0 {
			fmt.Println(dimStyle.Render(fmt.Sprintf("  %d agents already installed", len(installedAgents))))
		}
	} else {
		fmt.Println(subtitleStyle.Render("  Project Setup"))
	}
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Template: %s", tmplDir)))
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Target:   %s", targetDir)))
	fmt.Println()

	// Step 1: Ask if user wants BMAD methodology (skip if already installed)
	useBmad := false
	bmadAlreadyInstalled := installedAgents["product-owner"] && installedAgents["architect"] && installedAgents["tech-lead"]
	if !bmadAlreadyInstalled {
		bmadForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add BMAD methodology? (Break -> Model -> Act -> Deliver)").
					Description("Pre-selects core agents (product-owner, architect, tech-lead) + workflow commands.").
					Value(&useBmad),
			),
		).WithTheme(ckTheme())
		if err := bmadForm.Run(); err != nil {
			return err
		}
	}

	// Step 2: Agent-only picker — only show agents not yet installed
	// BMAD agents that aren't installed yet get pre-selected
	bmadAgents := map[string]bool{"product-owner": true, "architect": true, "tech-lead": true}

	var preselected []string
	options := make([]huh.Option[string], 0, len(agentComps))
	for _, c := range agentComps {
		if installedAgents[c.Name] {
			continue // skip already installed
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
		if useBmad && bmadAgents[c.Name] {
			preselected = append(preselected, c.Name)
		}
	}

	if len(options) == 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("  %s All agents already installed!", arrow)))
		fmt.Println(dimStyle.Render("  Use 'ck remove' to remove agents."))
		return nil
	}

	selectedAgents := preselected
	agentForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select agents to add (skills, commands & rules are automatic)").
				Options(options...).
				Value(&selectedAgents),
		),
	).WithTheme(ckTheme())

	if err := agentForm.Run(); err != nil {
		return err
	}

	if len(selectedAgents) == 0 {
		fmt.Println("No agents selected.")
		return nil
	}

	// Step 2b: Ask for teammate mode if 2+ agents selected
	teammateMode := "auto"
	if len(selectedAgents) >= 2 {
		teammateModeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Teammate display mode").
					Options(
						huh.NewOption("auto -- split panes in tmux, otherwise in-process", "auto"),
						huh.NewOption("in-process -- all teammates in one terminal", "in-process"),
						huh.NewOption("tmux -- each teammate in its own pane (requires tmux/iTerm2)", "tmux"),
					).
					Value(&teammateMode),
			),
		).WithTheme(ckTheme())
		if err := teammateModeForm.Run(); err != nil {
			return err
		}
	}

	// Step 3: Auto-compute all defaults from selected agents
	selectedSet := make(map[string]bool)
	for _, name := range selectedAgents {
		selectedSet[name] = true
	}

	// Compute skills, commands, rules from all selected agents
	skillSet := make(map[string]bool)
	commandSet := make(map[string]bool)
	ruleSet := make(map[string]bool)

	for _, name := range selectedAgents {
		// Skills from frontmatter
		agentPath := filepath.Join(tmplDir, "agents", name+".md")
		for _, dep := range catalog.ExtractSkillDeps(agentPath) {
			skillSet[dep] = true
		}

		// Extra skills
		for _, skill := range guessExtraSkillsForAgent(name) {
			skillSet[skill] = true
		}

		// Rules
		for _, rule := range guessRulesForAgent(name) {
			ruleSet[rule] = true
		}

		// Commands
		for _, cmd := range guessCommandsForAgent(name) {
			commandSet[cmd] = true
		}
	}

	// If BMAD accepted, add BMAD workflow commands
	if useBmad {
		for _, cmd := range []string{
			"bmad-run", "bmad-break", "bmad-model", "bmad-act", "bmad-deliver",
			"ralph", "ralph-loop", "ralph-cancel",
		} {
			commandSet[cmd] = true
		}
		// Base rules for BMAD
		for _, rule := range []string{"code-style", "testing", "security", "documentation"} {
			ruleSet[rule] = true
		}
	}

	// Always add ck-sync command
	commandSet["ck-sync"] = true

	// Convert sets to sorted slices for display
	skills := sortedKeys(skillSet)
	commands := sortedKeys(commandSet)
	rules := sortedKeys(ruleSet)

	// Step 4: Show summary
	fmt.Println()
	fmt.Println(subtitleStyle.Render("  Will install:"))
	fmt.Println(fmt.Sprintf("    %s %s: %s",
		bullet,
		accentStyle.Render(fmt.Sprintf("%d agents", len(selectedAgents))),
		dimStyle.Render(strings.Join(selectedAgents, ", ")),
	))
	fmt.Println(fmt.Sprintf("    %s %s: %s",
		bullet,
		accentStyle.Render(fmt.Sprintf("%d skills", len(skills))),
		dimStyle.Render(strings.Join(skills, ", ")),
	))
	fmt.Println(fmt.Sprintf("    %s %s: %s",
		bullet,
		accentStyle.Render(fmt.Sprintf("%d commands", len(commands))),
		dimStyle.Render(strings.Join(commands, ", ")),
	))
	fmt.Println(fmt.Sprintf("    %s %s: %s",
		bullet,
		accentStyle.Render(fmt.Sprintf("%d rules", len(rules))),
		dimStyle.Render(strings.Join(rules, ", ")),
	))
	if len(selectedAgents) >= 2 {
		fmt.Println(fmt.Sprintf("    %s %s: %s",
			bullet,
			accentStyle.Render("teammate mode"),
			dimStyle.Render(teammateMode),
		))
	}

	fmt.Println()

	// Step 5: Confirm
	var confirm bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Apply changes?").
				Value(&confirm),
		),
	).WithTheme(ckTheme())
	if err := confirmForm.Run(); err != nil {
		return err
	}
	if !confirm {
		fmt.Println("Aborted.")
		return nil
	}

	// Install base files
	if err := catalog.CopyBaseFiles(tmplDir, targetDir); err != nil {
		return fmt.Errorf("copying base files: %w", err)
	}
	if err := catalog.PatchSettingsTeammateMode(targetDir, teammateMode); err != nil {
		return fmt.Errorf("patching teammate mode: %w", err)
	}
	if !isExisting {
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render("Installed CLAUDE.md + settings.json")))
	}

	// Install agents
	fmt.Println(sectionHeader("Agents"))
	for _, name := range selectedAgents {
		if err := catalog.CopyComponent(tmplDir, targetDir, "agents", name); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  Agent %s: %v", name, err)))
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("agent: %s", name))))
	}

	// Install skills
	fmt.Println(sectionHeader("Skills"))
	for _, name := range skills {
		if err := catalog.CopyComponent(tmplDir, targetDir, "skills", name); err != nil {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("skill: %s", name))))
	}

	// Install commands
	fmt.Println(sectionHeader("Commands"))
	for _, name := range commands {
		if err := catalog.CopyComponent(tmplDir, targetDir, "commands", name); err != nil {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("command: %s", name))))
	}

	// Install rules
	fmt.Println(sectionHeader("Rules"))
	for _, name := range rules {
		if err := catalog.CopyComponent(tmplDir, targetDir, "rules", name); err != nil {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("rule: %s", name))))
	}

	// Install agent-teams rule if more than one agent selected
	if len(selectedAgents) > 1 {
		if err := catalog.CopyComponent(tmplDir, targetDir, "rules", "agent-teams"); err == nil {
			fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render("rule: agent-teams")))
		}
	}

	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Setup complete!", arrow)))
	fmt.Println(dimStyle.Render("  Run 'ck add' for more agents, 'ck remove' to remove components."))

	return nil
}

// sortedKeys returns the keys of a map sorted alphabetically.
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

