package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var (
	listAvailable bool
	listInstalled bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available and installed components",
	Long: `Show a table of BMAD template components.

By default, shows both available and installed components side by side.
Use --available or --installed to filter.`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVar(&listAvailable, "available", false, "Show available components only")
	listCmd.Flags().BoolVar(&listInstalled, "installed", false, "Show installed components only")
}

func runList(cmd *cobra.Command, args []string) error {
	tmplDir := resolveTemplateDir()

	fmt.Println(banner())

	available, err := catalog.ScanTemplate(tmplDir)
	if err != nil {
		return fmt.Errorf("scanning templates: %w", err)
	}

	targetDir := resolveTarget()
	installed, _ := catalog.GetInstalled(targetDir)

	// Build a set of installed component keys
	installedSet := make(map[string]bool)
	installedCount := 0
	for _, cat := range installed {
		for _, c := range cat.Components {
			installedSet[cat.Name+"/"+c.Name] = true
			installedCount++
		}
	}

	totalCount := 0
	for _, cat := range available {
		totalCount += len(cat.Components)
	}

	// Summary line
	summary := fmt.Sprintf("  %s %d installed  %s %d available",
		lipgloss.NewStyle().Foreground(green).Bold(true).Render("●"),
		installedCount,
		lipgloss.NewStyle().Foreground(dim).Render("●"),
		totalCount-installedCount,
	)
	fmt.Println(summary)

	for _, cat := range available {
		if listInstalled {
			hasInstalled := false
			for _, c := range cat.Components {
				if installedSet[cat.Name+"/"+c.Name] {
					hasInstalled = true
					break
				}
			}
			if !hasInstalled {
				continue
			}
		}

		fmt.Println(sectionHeader(strings.ToUpper(cat.Name)))

		rows := [][]string{}
		for _, c := range cat.Components {
			isInst := installedSet[cat.Name+"/"+c.Name]

			if listInstalled && !isInst {
				continue
			}
			if listAvailable && isInst {
				continue
			}

			status := dot
			nameRendered := dimStyle.Render(c.Name)
			if isInst {
				status = checkMark
				nameRendered = lipgloss.NewStyle().Foreground(white).Bold(true).Render(c.Name)
			}

			desc := c.Description
			if len(desc) > 55 {
				desc = desc[:52] + "..."
			}
			if isInst {
				desc = lipgloss.NewStyle().Foreground(pink).Render(desc)
			} else {
				desc = dimStyle.Render(desc)
			}

			rows = append(rows, []string{status, nameRendered, desc})
		}

		if len(rows) == 0 {
			fmt.Println(dimStyle.Render("    (none)"))
			continue
		}

		t := table.New().
			Border(lipgloss.HiddenBorder()).
			Headers(
				"",
				tableHeaderStyle.Render("Name"),
				tableHeaderStyle.Render("Description"),
			).
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				s := lipgloss.NewStyle().PaddingRight(2)
				if col == 0 {
					s = s.PaddingLeft(4).Width(3)
				} else if col == 1 {
					s = s.Width(30)
				}
				return s
			})

		fmt.Println(t)
	}

	fmt.Println()
	fmt.Println(dimStyle.Render("  Run 'ck add' to install agents interactively"))
	fmt.Println()
	return nil
}
