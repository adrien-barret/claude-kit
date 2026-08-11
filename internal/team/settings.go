package team

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// hookSpec is one (event, matcher, command) we ensure is wired.
type hookSpec struct {
	event   string
	matcher string // "" = no matcher key (SessionStart-style)
	command string
}

// teamHooks are the two generic hooks that drive the live poke + boot
// rescan. Both invoke this same binary; they no-op outside a `ck team`
// project (Discover returns ok=false), so they coexist with any other
// global hooks (e.g. Kerios's) without collision.
var teamHooks = []hookSpec{
	{event: "SessionStart", matcher: "", command: "ck team hook rescan"},
	{event: "PostToolUse", matcher: "Write|Edit", command: "ck team hook notify"},
}

// EnsureHooks idempotently merges teamHooks into the settings file at path
// (typically ~/.claude/settings.json), preserving every existing key. It
// backs the file up (.bak) before writing. Returns whether it changed
// anything.
func EnsureHooks(path string) (bool, error) {
	root := map[string]any{}
	orig, err := os.ReadFile(path)
	if err == nil && len(orig) > 0 {
		if err := json.Unmarshal(orig, &root); err != nil {
			return false, err
		}
	}

	hooks, _ := root["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}

	changed := false
	for _, h := range teamHooks {
		entries, _ := hooks[h.event].([]any)
		if hookCommandPresent(entries, h.command) {
			continue
		}
		entries = append(entries, newHookEntry(h))
		hooks[h.event] = entries
		changed = true
	}
	if !changed {
		return false, nil
	}
	root["hooks"] = hooks

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	if len(orig) > 0 {
		_ = os.WriteFile(path+".bak", orig, 0o644) // best-effort backup
	}
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return false, err
	}
	if err := os.WriteFile(path, append(out, '\n'), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// hookCommandPresent reports whether command is already registered in any of
// the event's entries (so re-running never duplicates it).
func hookCommandPresent(entries []any, command string) bool {
	for _, e := range entries {
		entry, _ := e.(map[string]any)
		inner, _ := entry["hooks"].([]any)
		for _, h := range inner {
			hm, _ := h.(map[string]any)
			if cmd, _ := hm["command"].(string); cmd == command {
				return true
			}
		}
	}
	return false
}

func newHookEntry(h hookSpec) map[string]any {
	entry := map[string]any{
		"hooks": []any{
			map[string]any{"type": "command", "command": h.command},
		},
	}
	if h.matcher != "" {
		entry["matcher"] = h.matcher
	}
	return entry
}
