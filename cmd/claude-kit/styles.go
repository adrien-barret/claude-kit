package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Zywoo-inspired color palette — hot pink + gold yellow
var (
	pink   = lipgloss.Color("#FF69B4")
	rose   = lipgloss.Color("#E91E8B")
	yellow = lipgloss.Color("#FFD700")
	gold   = lipgloss.Color("#FFC107")
	green  = lipgloss.Color("#00FF87")
	dim    = lipgloss.Color("#6C6C6C")
	white  = lipgloss.Color("#FAFAFA")
	red    = lipgloss.Color("#FF5555")

	// Base styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellow).
			Background(rose).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(pink).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(gold)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(pink)

	dimStyle = lipgloss.NewStyle().
			Foreground(dim)

	accentStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	// Box for banners
	bannerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(rose).
			Padding(0, 2).
			MarginBottom(1)

	// Table header
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(yellow).
				Underline(true)

	// Checkmark and bullet styles
	checkMark = lipgloss.NewStyle().Foreground(green).Bold(true).Render("✓")
	bullet    = lipgloss.NewStyle().Foreground(rose).Bold(true).Render("●")
	arrow     = lipgloss.NewStyle().Foreground(yellow).Bold(true).Render("→")
	dot       = lipgloss.NewStyle().Foreground(dim).Render("·")
)

func banner() string {
	logo := lipgloss.NewStyle().Foreground(rose).Bold(true).Render(`
   _____ _                 _        _  ___ _
  / ____| |               | |      | |/ (_) |
 | |    | | __ _ _   _  __| | ___  | ' / _| |_
 | |    | |/ _` + "`" + ` | | | |/ _` + "`" + ` |/ _ \ |  < | | __|
 | |____| | (_| | |_| | (_| |  __/ | . \| | |_
  \_____|_|\__,_|\__,_|\__,_|\___| |_|\_\_|\__|`)

	tagline := lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true).
		Render("  ⚡ Claude Code Project Templates")

	return fmt.Sprintf("%s\n%s\n", logo, tagline)
}

func sectionHeader(title string) string {
	return fmt.Sprintf("\n%s %s\n",
		lipgloss.NewStyle().Foreground(rose).Bold(true).Render("▸"),
		lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(title),
	)
}

func installedBadge() string {
	return lipgloss.NewStyle().
		Foreground(green).
		Background(lipgloss.Color("#1a3a1a")).
		Padding(0, 1).
		Render("installed")
}

// ckTheme returns a custom huh theme matching the zywoo color palette.
func ckTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Yellow checkmark for selected items
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(yellow).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(dim).SetString("  ")

	// Rose/pink for selected option text
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(rose)
	t.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(white)

	// Rose cursor selector (the > indicator)
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(rose).SetString("> ")

	// Pink titles
	t.Focused.Title = lipgloss.NewStyle().Foreground(pink).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Foreground(dim)

	// Yellow focused button, dim blurred button
	t.Focused.FocusedButton = lipgloss.NewStyle().Foreground(yellow).Background(rose).Bold(true).Padding(0, 1)
	t.Focused.BlurredButton = lipgloss.NewStyle().Foreground(dim).Padding(0, 1)

	// Blurred state — dimmed versions
	t.Blurred.SelectedPrefix = lipgloss.NewStyle().Foreground(gold).SetString("✓ ")
	t.Blurred.UnselectedPrefix = lipgloss.NewStyle().Foreground(dim).SetString("  ")
	t.Blurred.SelectedOption = lipgloss.NewStyle().Foreground(pink)
	t.Blurred.UnselectedOption = lipgloss.NewStyle().Foreground(dim)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(dim)
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().Foreground(dim).SetString("  ")

	return t
}
