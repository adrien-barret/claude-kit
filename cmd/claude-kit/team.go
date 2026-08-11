package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/AdeptMind/infra-tool/claude-cli/internal/team"
)

var teamCmd = &cobra.Command{
	Use:   "team <ProjectName>",
	Short: "Scaffold a multi-role team: role subfolders, a shared board, hooks, and tmux",
	Long: `team creates a project with one subfolder per role, a shared board, a
per-role CLAUDE.md, generic auto-reload hooks, and a launched tmux session —
a project-agnostic version of the Kerios department board system.

A temporary headless Claude session proposes the roles + objectives from a
short project description; you confirm; ck does all the file writes. Falls
back to a manual picker if the Claude CLI isn't available.

Roles coordinate through board/BOARD.md using ` + "`@role:`" + ` requests. The
hooks (installed once, globally) poke a role's tmux pane when a request for
it lands and re-list pending requests when a pane boots. They discover the
project via .team/config.json and no-op everywhere else, so they coexist
with other global hooks.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTeamScaffold,
}

// hidden hook dispatcher: `ck team hook notify|rescan`, wired into global
// settings.json. Reads the Claude hook JSON on stdin.
var teamHookCmd = &cobra.Command{
	Use:    "hook",
	Short:  "Internal: team board hooks (invoked by Claude Code)",
	Hidden: true,
}

var teamHookNotifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "PostToolUse: poke role panes when a board request lands",
	RunE:  runTeamHookNotify,
}

var teamHookRescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "SessionStart: list pending board requests for this pane's role",
	RunE:  runTeamHookRescan,
}

func init() {
	teamHookCmd.AddCommand(teamHookNotifyCmd, teamHookRescanCmd)
	teamCmd.AddCommand(teamHookCmd)
}

// hookEvent is the subset of the Claude Code hook JSON we consume.
type hookEvent struct {
	Cwd       string `json:"cwd"`
	ToolInput struct {
		FilePath string `json:"file_path"`
	} `json:"tool_input"`
}

func readHookEvent() hookEvent {
	var ev hookEvent
	_ = json.NewDecoder(os.Stdin).Decode(&ev)
	if ev.Cwd == "" {
		ev.Cwd, _ = os.Getwd()
	}
	return ev
}

func runTeamHookNotify(cmd *cobra.Command, _ []string) error {
	ev := readHookEvent()
	cfg, root, ok := team.Discover(ev.Cwd)
	if !ok {
		return nil // not a team project → no-op
	}
	if ev.ToolInput.FilePath == "" {
		return nil
	}
	written := ev.ToolInput.FilePath
	if !filepath.IsAbs(written) {
		written = filepath.Join(ev.Cwd, written)
	}
	if filepath.Clean(written) != filepath.Clean(cfg.BoardPath(root)) {
		return nil // only board writes poke panes
	}
	content, err := os.ReadFile(cfg.BoardPath(root))
	if err != nil {
		return nil
	}
	for _, target := range team.NotifyTargets(cfg, string(content)) {
		role := target[strings.LastIndex(target, ":")+1:]
		msg := fmt.Sprintf("[team] new @%s: request on the board — re-read %s", role, cfg.Board)
		_ = exec.Command("tmux", "send-keys", "-t", target, msg, "C-m").Run() // best-effort
	}
	return nil
}

func runTeamHookRescan(cmd *cobra.Command, _ []string) error {
	ev := readHookEvent()
	cfg, root, ok := team.Discover(ev.Cwd)
	if !ok {
		return nil
	}
	role, ok := team.RoleForDir(cfg, root, ev.Cwd)
	if !ok {
		return nil
	}
	content, err := os.ReadFile(cfg.BoardPath(root))
	if err != nil {
		return nil
	}
	pending := team.PendingRequests(string(content), role.Name)
	if len(pending) == 0 {
		return nil
	}
	fmt.Printf("[team] %d pending @%s: request(s) on the board (%s):\n", len(pending), role.Name, cfg.Board)
	for _, l := range pending {
		fmt.Println(l)
	}
	return nil
}

func runTeamScaffold(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ck team <ProjectName>")
	}
	project := args[0]
	root := filepath.Join(resolveProjectRoot(), project)

	fmt.Println(banner())
	fmt.Println(sectionHeader("Team: " + project))

	// 1. Project description → headless-Claude proposal (or fallback).
	desc := askDescription()
	proposal, err := team.Propose(desc)
	if err != nil {
		fmt.Println(dimStyle.Render("  Claude proposal unavailable — falling back to the role picker."))
		proposal = fallbackProposal(desc)
	}
	ensureLead(&proposal)

	// 2. Confirm / trim the role set.
	roles := confirmRoles(proposal.Roles)
	if len(roles) == 0 {
		return fmt.Errorf("no roles selected")
	}

	// 3. Scaffold.
	if err := scaffoldTeam(root, project, proposal.BoardFocus, roles); err != nil {
		return err
	}

	// 4. Wire the generic hooks globally (idempotent).
	if home, herr := os.UserHomeDir(); herr == nil {
		settings := filepath.Join(home, ".claude", "settings.json")
		if changed, err := team.EnsureHooks(settings); err != nil {
			fmt.Println(warnStyle.Render("  ! could not wire hooks: " + err.Error()))
		} else if changed {
			fmt.Printf("  %s wired team hooks into %s\n", checkMark, dimStyle.Render(settings))
		} else {
			fmt.Printf("  %s team hooks already wired\n", checkMark)
		}
	}

	// 5. Launch tmux (or print the command).
	launchOrHint(root, project, roles)

	fmt.Printf("\n%s team %s ready — %d roles\n", checkMark, accentStyle.Render(project), len(roles))
	return nil
}

func askDescription() string {
	desc := ""
	form := huh.NewForm(huh.NewGroup(
		huh.NewText().
			Title("Describe the project in a sentence or two").
			Description("Used to propose the roles and their objectives.").
			Value(&desc),
	)).WithTheme(ckTheme())
	if err := form.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(desc)
}

// confirmRoles shows a preselected multi-select so the user can drop any
// proposed role before scaffolding.
func confirmRoles(proposed []team.RoleSpec) []team.RoleSpec {
	opts := make([]huh.Option[int], 0, len(proposed))
	for i, r := range proposed {
		label := r.Name
		if strings.TrimSpace(r.Mission) != "" {
			label = fmt.Sprintf("%s — %s", r.Name, r.Mission)
		}
		opts = append(opts, huh.NewOption(label, i).Selected(true))
	}
	var keep []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title("Confirm the roles for this team").
			Options(opts...).
			Value(&keep),
	)).WithTheme(ckTheme())
	if err := form.Run(); err != nil || len(keep) == 0 {
		return proposed // non-interactive or aborted → keep all
	}
	out := make([]team.RoleSpec, 0, len(keep))
	for _, i := range keep {
		out = append(out, proposed[i])
	}
	return out
}

// ensureLead guarantees a lead role exists (it owns Current Focus + [DONE]).
func ensureLead(p *team.Proposal) {
	for _, r := range p.Roles {
		if r.Name == "lead" {
			return
		}
	}
	p.Roles = append([]team.RoleSpec{{
		Name:       "lead",
		Folder:     "lead",
		Mission:    "Coordinate the team, dispatch work, and verify Definitions of Done.",
		Objectives: []string{"Keep Current Focus current", "Dispatch @role: requests with a DoD", "Verify proof and mark [DONE]"},
	}}, p.Roles...)
}

// fallbackProposal builds a proposal from the common archetypes when the
// Claude CLI is unavailable.
func fallbackProposal(desc string) team.Proposal {
	names := []string{"lead", "eng", "marketing", "sales", "design", "research", "ops"}
	roles := make([]team.RoleSpec, 0, len(names))
	for _, n := range names {
		roles = append(roles, team.RoleSpec{
			Name:       n,
			Folder:     n,
			Mission:    fmt.Sprintf("Own the %s workstream.", n),
			Objectives: []string{fmt.Sprintf("Advance %s goals", n), "Coordinate via the board"},
		})
	}
	focus := strings.TrimSpace(desc)
	if focus == "" {
		focus = "Set the initial team focus."
	}
	return team.Proposal{BoardFocus: focus, Roles: roles}
}

func scaffoldTeam(root, project, focus string, roles []team.RoleSpec) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	templateDir := resolveTemplateDir()
	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = r.Name
	}

	// Per-role CLAUDE.md (skip existing — idempotent re-run).
	for _, r := range roles {
		folder := filepath.Join(root, r.Folder)
		if err := os.MkdirAll(folder, 0o755); err != nil {
			return err
		}
		dst := filepath.Join(folder, "CLAUDE.md")
		if _, err := os.Stat(dst); err == nil {
			fmt.Printf("  %s %s/CLAUDE.md exists — kept\n", dimStyle.Render("•"), r.Folder)
			continue
		}
		content, err := team.RenderRole(templateDir, team.RoleTemplateData{
			Project:    project,
			Role:       r.Name,
			Mission:    r.Mission,
			Objectives: r.Objectives,
			BoardRel:   team.BoardRelFrom(r.Folder, "board/BOARD.md"),
			Roles:      names,
		})
		if err != nil {
			return fmt.Errorf("render role %s: %w", r.Name, err)
		}
		if err := os.WriteFile(dst, []byte(content), 0o644); err != nil {
			return err
		}
		fmt.Printf("  %s %s/CLAUDE.md\n", checkMark, r.Folder)
	}

	// Board (never clobber an existing BOARD.md) + resource requests.
	boardDir := filepath.Join(root, "board")
	if err := os.MkdirAll(boardDir, 0o755); err != nil {
		return err
	}
	boardPath := filepath.Join(boardDir, "BOARD.md")
	if _, err := os.Stat(boardPath); err != nil {
		if err := os.WriteFile(boardPath, []byte(team.RenderBoard(project, focus, roles)), 0o644); err != nil {
			return err
		}
		fmt.Printf("  %s board/BOARD.md\n", checkMark)
	} else {
		fmt.Printf("  %s board/BOARD.md exists — kept\n", dimStyle.Render("•"))
	}
	rrPath := filepath.Join(boardDir, "RESOURCE-REQUESTS.md")
	if _, err := os.Stat(rrPath); err != nil {
		_ = os.WriteFile(rrPath, []byte(team.RenderResourceRequests(project)), 0o644)
	}

	// Discovery anchor.
	cfg := team.Config{
		Project:     project,
		TmuxSession: team.Slug(project),
		Board:       "board/BOARD.md",
		Roles:       toRoles(roles),
	}
	if err := cfg.WriteTo(root); err != nil {
		return err
	}
	fmt.Printf("  %s .team/config.json\n", checkMark)
	return nil
}

func toRoles(specs []team.RoleSpec) []team.Role {
	out := make([]team.Role, len(specs))
	for i, s := range specs {
		out[i] = team.Role{Name: s.Name, Folder: s.Folder}
	}
	return out
}

// launchOrHint starts the tmux session (one window per role + a board
// window) or, if tmux is missing / the session exists, prints what to run.
func launchOrHint(root, project string, roles []team.RoleSpec) {
	session := team.Slug(project)
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		fmt.Println(dimStyle.Render("  tmux not found — launch it yourself when ready."))
		return
	}
	if exec.Command(tmuxPath, "has-session", "-t", session).Run() == nil {
		fmt.Printf("  %s tmux session %q already exists — left as-is\n", dimStyle.Render("•"), session)
		return
	}
	first := roles[0]
	if err := exec.Command(tmuxPath, "new-session", "-d", "-s", session,
		"-n", first.Name, "-c", filepath.Join(root, first.Folder)).Run(); err != nil {
		fmt.Println(warnStyle.Render("  ! tmux launch failed: " + err.Error()))
		return
	}
	for _, r := range roles[1:] {
		_ = exec.Command(tmuxPath, "new-window", "-t", session, "-n", r.Name,
			"-c", filepath.Join(root, r.Folder)).Run()
	}
	_ = exec.Command(tmuxPath, "new-window", "-t", session, "-n", "board",
		"-c", filepath.Join(root, "board")).Run()
	fmt.Printf("  %s tmux session %s launched (attach: tmux attach -t %s)\n",
		checkMark, accentStyle.Render(session), session)
}
