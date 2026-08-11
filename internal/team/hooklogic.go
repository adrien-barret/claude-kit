package team

import (
	"fmt"
	"strings"
)

// requestLine reports whether a board line is an active request for role —
// it contains `@<role>:` and is not yet marked [DONE].
func requestLine(line, role string) bool {
	if !strings.Contains(line, "@"+role+":") {
		return false
	}
	// Done markers carry a date inside the brackets in practice
	// (`[DONE 2026-08-11]`), so match the opening `[DONE`, not `[DONE]`.
	return !strings.Contains(strings.ToUpper(line), "[DONE")
}

// NotifyTargets returns the tmux window targets (`session:role`) to poke,
// given the current board content and config. A role is poked when the
// board has at least one active (un-[DONE]) `@role:` request. Empty result
// → nothing to poke. Pure: the caller performs the tmux send-keys.
func NotifyTargets(cfg Config, boardContent string) []string {
	lines := strings.Split(boardContent, "\n")
	var targets []string
	for _, r := range cfg.Roles {
		for _, l := range lines {
			if requestLine(l, r.Name) {
				targets = append(targets, fmt.Sprintf("%s:%s", cfg.TmuxSession, r.Name))
				break
			}
		}
	}
	return targets
}

// PendingRequests returns the active (un-[DONE]) `@role:` request lines in
// the board, trimmed. Used by the boot rescan to surface what a role owes.
func PendingRequests(boardContent, role string) []string {
	var out []string
	for _, l := range strings.Split(boardContent, "\n") {
		if requestLine(l, role) {
			out = append(out, strings.TrimSpace(l))
		}
	}
	return out
}
