// Package team implements `ck team` — project-agnostic multi-role pane
// orchestration (a shared board + per-role CLAUDE.md + generic hooks +
// tmux), a reusable generalization of the Kerios department board system.
package team

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigDir is the per-project directory holding the discovery anchor.
const ConfigDir = ".team"

// ConfigFile is the discovery anchor filename inside ConfigDir.
const ConfigFile = "config.json"

// Role is one team member: a name (used in `@name:` board requests and as
// the tmux window name) and the project-relative folder it works in.
type Role struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
}

// Config is the contract every hook reads. It lives at
// <root>/.team/config.json and is the sole source of truth for the board
// location and the role→folder map.
type Config struct {
	Project     string `json:"project"`
	TmuxSession string `json:"tmux_session"`
	// Board is the board file path, relative to the project root.
	Board string `json:"board"`
	Roles []Role `json:"roles"`
}

// BoardPath returns the absolute board path for a project rooted at root.
func (c Config) BoardPath(root string) string {
	return filepath.Join(root, filepath.FromSlash(c.Board))
}

// WriteTo writes the config to <root>/.team/config.json (creating the dir).
func (c Config) WriteTo(root string) error {
	dir := filepath.Join(root, ConfigDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ConfigFile), append(data, '\n'), 0o644)
}

// Discover walks up from startDir looking for a .team/config.json. On
// success it returns the parsed config and the project root (the directory
// that contains .team). ok is false when no config is found up to the
// filesystem root — the signal for a hook to no-op.
func Discover(startDir string) (cfg Config, root string, ok bool) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ConfigDir, ConfigFile)
		if data, err := os.ReadFile(candidate); err == nil {
			if json.Unmarshal(data, &cfg) == nil {
				return cfg, dir, true
			}
			return Config{}, "", false // present but unparseable → treat as absent
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return Config{}, "", false // reached filesystem root
		}
		dir = parent
	}
}

// RoleForDir returns the role whose folder contains dir (dir is at or below
// <root>/<folder>). Used by the rescan hook to map a pane's cwd to its role.
func RoleForDir(cfg Config, root, dir string) (Role, bool) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}
	for _, r := range cfg.Roles {
		roleRoot, err := filepath.Abs(filepath.Join(root, r.Folder))
		if err != nil {
			continue
		}
		if absDir == roleRoot || isUnder(absDir, roleRoot) {
			return r, true
		}
	}
	return Role{}, false
}

// isUnder reports whether path is inside base (a strict descendant).
func isUnder(path, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel != ".." && !filepath.IsAbs(rel) && !hasDotDotPrefix(rel)
}

func hasDotDotPrefix(rel string) bool {
	return len(rel) >= 2 && rel[0] == '.' && rel[1] == '.'
}
