package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/catalog"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Claude account profiles",
	Long: `Manage Claude Code account profiles (multi-account via CLAUDE_CONFIG_DIR).

Without arguments, shows an interactive profile picker then launches claude.

Examples:
  ck profile                    # Interactive profile selector
  ck profile list               # List configured profiles
  ck profile use <name>         # Launch Claude with the given profile
  ck profile add <name> <path>  # Register a new profile
  ck profile remove <name>      # Remove a profile`,
	RunE: runProfileTUI,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured profiles",
	RunE:  runProfileList,
}

var profileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Launch Claude with the given profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileUse,
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name> <path>",
	Short: "Register a new profile",
	Args:  cobra.ExactArgs(2),
	RunE:  runProfileAdd,
}

var profileRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Short:   "Remove a profile",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runProfileRemove,
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRemoveCmd)
}

func runProfileTUI(cmd *cobra.Command, args []string) error {
	store, err := catalog.ReadProfiles()
	if err != nil {
		return fmt.Errorf("reading profiles: %w", err)
	}

	if len(store.Profiles) == 0 {
		fmt.Println(warnStyle.Render("  No profiles configured."))
		fmt.Println(dimStyle.Render("  Use 'ck profile add <name> <path>' to add one."))
		return nil
	}

	names := profileSortedKeys(store.Profiles)
	options := make([]huh.Option[string], len(names))
	for i, name := range names {
		label := fmt.Sprintf("%-15s %s", name, dimStyle.Render(store.Profiles[name]))
		options[i] = huh.NewOption(label, name)
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Claude profile").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(ckTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if selected == "" {
		return nil
	}

	return launchWithProfile(store.Profiles[selected])
}

func runProfileList(cmd *cobra.Command, args []string) error {
	store, err := catalog.ReadProfiles()
	if err != nil {
		return fmt.Errorf("reading profiles: %w", err)
	}

	if len(store.Profiles) == 0 {
		fmt.Println(dimStyle.Render("  No profiles configured."))
		fmt.Println(dimStyle.Render("  Use 'ck profile add <name> <path>' to add one."))
		return nil
	}

	fmt.Println(sectionHeader("Profiles"))
	for _, name := range profileSortedKeys(store.Profiles) {
		fmt.Printf("  %s %-15s %s\n", bullet, accentStyle.Render(name), dimStyle.Render(store.Profiles[name]))
	}
	return nil
}

func runProfileUse(cmd *cobra.Command, args []string) error {
	name := args[0]
	store, err := catalog.ReadProfiles()
	if err != nil {
		return fmt.Errorf("reading profiles: %w", err)
	}

	path, ok := store.Profiles[name]
	if !ok {
		return fmt.Errorf("profile %q not found — use 'ck profile list' to see available profiles", name)
	}

	return launchWithProfile(path)
}

func runProfileAdd(cmd *cobra.Command, args []string) error {
	name, path := args[0], args[1]

	expanded := profileExpandHome(path)
	if _, err := os.Stat(expanded); os.IsNotExist(err) {
		fmt.Println(warnStyle.Render(fmt.Sprintf("  Warning: path %q does not exist (adding anyway)", path)))
	}

	store, err := catalog.ReadProfiles()
	if err != nil {
		return fmt.Errorf("reading profiles: %w", err)
	}

	store.Profiles[name] = path
	if err := catalog.WriteProfiles(store); err != nil {
		return fmt.Errorf("writing profiles: %w", err)
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Profile %q added → %s", arrow, name, path)))
	return nil
}

func runProfileRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	store, err := catalog.ReadProfiles()
	if err != nil {
		return fmt.Errorf("reading profiles: %w", err)
	}

	if _, ok := store.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}

	delete(store.Profiles, name)
	if err := catalog.WriteProfiles(store); err != nil {
		return fmt.Errorf("writing profiles: %w", err)
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("  %s Profile %q removed", arrow, name)))
	return nil
}

// launchWithProfile replaces the current process with claude using CLAUDE_CONFIG_DIR.
func launchWithProfile(configDir string) error {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude not found in PATH: %w", err)
	}

	expanded := profileExpandHome(configDir)
	env := append(os.Environ(), "CLAUDE_CONFIG_DIR="+expanded)
	return syscall.Exec(claudePath, []string{"claude"}, env)
}

// profileExpandHome expands a leading ~ to the user home directory.
func profileExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}

// profileSortedKeys returns sorted keys of a string map.
func profileSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
