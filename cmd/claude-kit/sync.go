package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
	"github.com/AdeptMind/infra-tool/claude-cli/internal/docsindex"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update installed components and refresh docs-index",
	Long: `Sync updates installed components from the template catalog.

Only components that are already installed are updated — no new components
are added. After updating, the docs-index is refreshed if stale.`,
	RunE: runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	tmplDir := resolveTemplateDir()
	targetDir := resolveTarget()

	fmt.Println(banner())

	if _, err := os.Stat(targetDir); err != nil {
		return fmt.Errorf("no .claude directory found at %s — run 'ck init' first", targetDir)
	}

	var updated int
	var syncErr error

	action := func() {
		// Get installed components
		installed, err := catalog.GetInstalled(targetDir)
		if err != nil {
			syncErr = fmt.Errorf("scanning installed: %w", err)
			return
		}

		// Update base files
		if err := catalog.CopyBaseFiles(tmplDir, targetDir); err != nil {
			syncErr = fmt.Errorf("updating base files: %w", err)
			return
		}
		updated++

		// Update each installed component from template
		for _, cat := range installed {
			for _, comp := range cat.Components {
				err := catalog.CopyComponent(tmplDir, targetDir, cat.Name, comp.Name)
				if err != nil {
					// Skip components not in template (user-created)
					continue
				}
				updated++
			}
		}
	}

	if err := spinner.New().
		Title("Syncing components...").
		Action(action).
		Run(); err != nil {
		return err
	}

	if syncErr != nil {
		return syncErr
	}

	fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render(fmt.Sprintf("Updated %d components", updated))))

	// Refresh docs-index
	projectRoot := filepath.Dir(targetDir)
	if strings.HasSuffix(targetDir, ".claude") {
		stale, reason := docsindex.IsStale(projectRoot)
		if stale {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  %s Docs-index needs refresh: %s", bullet, reason)))

			var techs []string
			var docsErr error

			docsAction := func() {
				techs, docsErr = docsindex.Generate(projectRoot)
			}

			if err := spinner.New().
				Title("Refreshing docs-index...").
				Action(docsAction).
				Run(); err != nil {
				return err
			}

			if docsErr != nil {
				fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("  Docs refresh failed: %v", docsErr)))
			} else {
				fmt.Println(fmt.Sprintf("  %s %s", checkMark, infoStyle.Render(fmt.Sprintf("Docs-index refreshed (stack: %s)", strings.Join(techs, ", ")))))
			}
		} else {
			fmt.Println(fmt.Sprintf("  %s %s", checkMark, dimStyle.Render("Docs-index is up to date")))
		}
	}

	return nil
}
