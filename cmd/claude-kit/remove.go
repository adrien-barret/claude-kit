package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var removeCmd = &cobra.Command{
	Use:   "remove [names...] | remove <type> <name>",
	Short: "Remove installed components",
	Long: `Remove components from the current project's .claude/ directory.

When called with no arguments, shows an interactive picker of installed components.
When called with names, removes those agents (or type + name for other components).

Examples:
  ck remove                     # Interactive picker
  ck remove backend             # Remove the backend agent
  ck remove skill code-reviewer # Remove a specific skill`,
	RunE: runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	targetDir := resolveTarget()

	fmt.Println(banner())

	// No args → interactive
	if len(args) == 0 {
		return runInteractiveRemove(targetDir)
	}

	// Check if first arg is a type prefix
	firstNorm := normalizeType(args[0])
	if isComponentType(firstNorm) && len(args) >= 2 {
		return removeComponents(targetDir, firstNorm, args[1:])
	}

	// Otherwise treat as agent names
	return removeComponents(targetDir, "agents", args)
}

func runInteractiveRemove(targetDir string) error {
	installed, err := catalog.GetInstalled(targetDir)
	if err != nil || len(installed) == 0 {
		return fmt.Errorf("no components installed in %s", targetDir)
	}

	fmt.Println(subtitleStyle.Render("  Remove installed components"))
	fmt.Println()

	// Build options from all installed components
	type componentRef struct {
		compType string
		name     string
	}

	var options []huh.Option[string]
	refMap := make(map[string]componentRef)

	for _, cat := range installed {
		for _, c := range cat.Components {
			key := cat.Name + "/" + c.Name
			label := fmt.Sprintf("[%s] %s", cat.Name, c.Name)
			options = append(options, huh.NewOption(label, key))
			refMap[key] = componentRef{compType: cat.Name, name: c.Name}
		}
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select components to remove").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if len(selected) == 0 {
		fmt.Println("Nothing selected.")
		return nil
	}

	// Confirm
	var confirm bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Remove %d components?", len(selected))).
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

	for _, key := range selected {
		ref := refMap[key]

		// Warn if skill is referenced
		if ref.compType == "skills" {
			refs := catalog.FindReferencingAgents(targetDir, ref.name)
			if len(refs) > 0 {
				fmt.Println(warnStyle.Render(fmt.Sprintf("    %s is used by agents: %v", ref.name, refs)))
			}
		}

		if err := catalog.RemoveComponent(targetDir, ref.compType, ref.name); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  Could not remove %s: %v", key, err)))
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Removed %s", key))))
	}

	fmt.Println()
	return nil
}

func removeComponents(targetDir, compType string, names []string) error {
	for _, name := range names {
		if !catalog.IsInstalled(targetDir, compType, name) {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  %s/%s is not installed", compType, name)))
			continue
		}

		// Warn if skill is referenced
		if compType == "skills" {
			refs := catalog.FindReferencingAgents(targetDir, name)
			if len(refs) > 0 {
				fmt.Println(warnStyle.Render(fmt.Sprintf("    %s is used by agents: %v", name, refs)))
			}
		}

		if err := catalog.RemoveComponent(targetDir, compType, name); err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  Could not remove %s/%s: %v", compType, name, err)))
			continue
		}
		fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Removed %s/%s", compType, name))))
	}
	fmt.Println()
	return nil
}
