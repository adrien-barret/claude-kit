package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// DepType determines how a dependency is installed.
type DepType string

const (
	DepTypePlugin DepType = "plugin"
	DepTypeBrew   DepType = "brew"
	DepTypeNpm    DepType = "npm"
	DepTypeGo     DepType = "go"
	DepTypeShell  DepType = "shell"
)

// Dependency represents an installable tool or plugin.
type Dependency struct {
	Name                 string
	Description          string
	Type                 DepType
	Source               string // repo or package identifier
	PluginMarketplaceCmd string // for plugin type
	PluginInstallCmd     string // for plugin type
	InstallCmd           string // for auto-installable types
}

var depRegistry = []Dependency{
	{
		Name:                 "claude-mem",
		Description:          "Persistent memory compression system for Claude Code",
		Type:                 DepTypePlugin,
		Source:               "thedotmack/claude-mem",
		PluginMarketplaceCmd: "/plugin marketplace add thedotmack/claude-mem",
		PluginInstallCmd:     "/plugin install claude-mem",
	},
}

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Manage recommended dependencies and tools",
	Long: `Manage recommended dependencies and tools for Claude Code projects.

Use subcommands to install plugins, CLI tools, and other dependencies.`,
}

var depInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install recommended dependencies interactively",
	Long: `Interactively select and install recommended dependencies.

Plugins require manual steps (slash commands in Claude Code).
Other dependency types (brew, npm, go, shell) are installed automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDepInstall()
	},
}

func init() {
	depCmd.AddCommand(depInstallCmd)
}

func runDepInstall() error {
	fmt.Println(banner())
	fmt.Println(subtitleStyle.Render("  Install recommended dependencies"))
	fmt.Println()

	// Build multi-select options from registry
	options := make([]huh.Option[int], 0, len(depRegistry))
	for i, dep := range depRegistry {
		label := fmt.Sprintf("%s -- %s", dep.Name, dep.Description)
		options = append(options, huh.NewOption(label, i))
	}

	var selected []int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[int]().
				Title("Select dependencies to install").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if len(selected) == 0 {
		fmt.Println("No dependencies selected.")
		return nil
	}

	// Separate plugin deps (manual) from auto-installable ones
	var pluginDeps []Dependency
	var autoDeps []Dependency
	for _, idx := range selected {
		dep := depRegistry[idx]
		if dep.Type == DepTypePlugin {
			pluginDeps = append(pluginDeps, dep)
		} else {
			autoDeps = append(autoDeps, dep)
		}
	}

	// Auto-install non-plugin deps
	for _, dep := range autoDeps {
		fmt.Println(sectionHeader(fmt.Sprintf("Installing %s", dep.Name)))
		if err := autoInstallDep(dep); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("  Failed to install %s: %v", dep.Name, err)))
		} else {
			fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(dep.Name)))
		}
	}

	// Print manual plugin instructions
	if len(pluginDeps) > 0 {
		fmt.Println(sectionHeader("Plugin Setup (manual steps in Claude Code)"))
		fmt.Println(dimStyle.Render("  Run these slash commands inside a Claude Code session:"))
		fmt.Println()

		step := 1
		for _, dep := range pluginDeps {
			fmt.Println(fmt.Sprintf("  %s %s",
				accentStyle.Render(fmt.Sprintf("%d.", step)),
				infoStyle.Render(dep.PluginMarketplaceCmd),
			))
			step++
			fmt.Println(fmt.Sprintf("  %s %s",
				accentStyle.Render(fmt.Sprintf("%d.", step)),
				infoStyle.Render(dep.PluginInstallCmd),
			))
			step++
		}
		fmt.Println()
	}

	// Summary
	installed := len(autoDeps)
	manual := len(pluginDeps)
	parts := make([]string, 0, 2)
	if installed > 0 {
		parts = append(parts, fmt.Sprintf("%d installed", installed))
	}
	if manual > 0 {
		parts = append(parts, fmt.Sprintf("%d require manual setup", manual))
	}
	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Done! %s", arrow, strings.Join(parts, ", "))))

	return nil
}

func autoInstallDep(dep Dependency) error {
	var cmd string
	var args []string

	switch dep.Type {
	case DepTypeBrew:
		cmd = "brew"
		args = []string{"install", dep.Source}
	case DepTypeNpm:
		cmd = "npm"
		args = []string{"install", "-g", dep.Source}
	case DepTypeGo:
		cmd = "go"
		args = []string{"install", dep.Source}
	case DepTypeShell:
		cmd = "sh"
		args = []string{"-c", dep.InstallCmd}
	default:
		return fmt.Errorf("unsupported dep type: %s", dep.Type)
	}

	// Check that the installer binary exists
	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("%s not found in PATH", cmd)
	}

	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
