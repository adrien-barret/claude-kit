package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/docsindex"
)

var docsRefresh bool

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate docs-index.md via stack detection",
	Long: `Detect the project's tech stack and generate a compressed docs-index.md
with framework-specific directives that Claude can use for context.

The docs-index is written to .claude/docs-index.md along with metadata
in .claude/.docs-meta.json for staleness tracking.

Use --refresh to force regeneration even if the current index is fresh.`,
	RunE: runDocs,
}

func init() {
	docsCmd.Flags().BoolVar(&docsRefresh, "refresh", false, "Force regenerate even if fresh")
}

func runDocs(cmd *cobra.Command, args []string) error {
	projectRoot := resolveProjectRoot()

	fmt.Println(banner())

	// Check staleness unless --refresh
	if !docsRefresh {
		stale, reason := docsindex.IsStale(projectRoot)
		if !stale {
			fmt.Println(fmt.Sprintf("  %s %s", checkMark, dimStyle.Render("Docs-index is up to date.")))
			fmt.Println(dimStyle.Render("    Use --refresh to force regeneration."))
			fmt.Println()
			return nil
		}
		fmt.Println(warnStyle.Render(fmt.Sprintf("  %s Regenerating: %s", bullet, reason)))
	}

	var techs []string
	var genErr error

	action := func() {
		techs, genErr = docsindex.Generate(projectRoot)
	}

	if err := spinner.New().
		Title("Detecting stack and generating docs-index...").
		Action(action).
		Run(); err != nil {
		return err
	}

	if genErr != nil {
		return fmt.Errorf("generating docs-index: %w", genErr)
	}

	fmt.Println(fmt.Sprintf("  %s %s", checkMark, accentStyle.Render("Generated .claude/docs-index.md")))

	if len(techs) > 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("    %s Detected stack: %s", arrow, strings.Join(techs, ", "))))
	} else {
		fmt.Println(dimStyle.Render("    No stack detected. Add dependency files and re-run."))
	}

	fmt.Println(dimStyle.Render("    Metadata: .claude/.docs-meta.json"))
	fmt.Println()

	return nil
}
