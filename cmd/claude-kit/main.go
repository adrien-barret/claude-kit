package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/config"
)

var version = "dev"

var (
	templateDir string
	projectDir  string
)

var rootCmd = &cobra.Command{
	Use:   "claude-kit",
	Short: "claude-kit (ck) — manage Claude Code project templates",
	Long: `claude-kit is a CLI for managing BMAD project templates.

It provides interactive setup, component management, stack-aware docs
generation, and template synchronization for Claude Code projects.

Alias: ck`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the claude-kit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&templateDir, "template-dir", "", "Override template directory path")
	rootCmd.PersistentFlags().StringVarP(&projectDir, "project", "f", "", "Project directory (default: current directory)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(docsCmd)
}

func resolveTemplateDir() string {
	if templateDir != "" {
		return templateDir
	}
	return config.TemplateDir()
}

// resolveProjectRoot returns the project root directory.
// Uses -C flag if set, otherwise the current working directory.
func resolveProjectRoot() string {
	if projectDir != "" {
		abs, err := filepath.Abs(projectDir)
		if err != nil {
			return projectDir
		}
		return abs
	}
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return "."
}

// resolveTarget returns the .claude target directory within the project root.
func resolveTarget() string {
	return filepath.Join(resolveProjectRoot(), ".claude")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
