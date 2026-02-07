package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var teammateModeCmd = &cobra.Command{
	Use:   "teammate-mode",
	Short: "View or change the teammate display mode",
	RunE:  runTeammateMode,
}

func runTeammateMode(cmd *cobra.Command, args []string) error {
	targetDir := resolveTarget()

	current, err := catalog.ReadSettingsTeammateMode(targetDir)
	if err != nil {
		return err
	}

	fmt.Println(banner())
	fmt.Println(fmt.Sprintf("  %s Current teammate mode: %s",
		bullet,
		accentStyle.Render(current),
	))
	fmt.Println()

	newMode := current
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Teammate display mode").
				Options(
					huh.NewOption("auto -- split panes in tmux, otherwise in-process", "auto"),
					huh.NewOption("in-process -- all teammates in one terminal", "in-process"),
					huh.NewOption("tmux -- each teammate in its own pane (requires tmux/iTerm2)", "tmux"),
				).
				Value(&newMode),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if newMode == current {
		fmt.Println(dimStyle.Render("  No change."))
		return nil
	}

	if err := catalog.PatchSettingsTeammateMode(targetDir, newMode); err != nil {
		return fmt.Errorf("updating teammate mode: %w", err)
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Teammate mode set to %s", arrow, newMode)))
	return nil
}
