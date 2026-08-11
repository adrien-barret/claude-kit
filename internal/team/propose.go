package team

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Proposal is what a temporary headless Claude session returns: a focus line
// for the board plus a proposed role set with per-role objectives. Claude is
// a pure advisor here — it never touches the filesystem; ck does all writes.
type Proposal struct {
	BoardFocus string     `json:"board_focus"`
	Roles      []RoleSpec `json:"roles"`
}

// ErrClaudeUnavailable signals the caller should fall back to the manual
// (huh) flow: the `claude` CLI is missing, errored, or returned no usable
// proposal.
var ErrClaudeUnavailable = fmt.Errorf("claude proposal unavailable")

// Propose asks a temporary headless Claude session to turn a project
// description into a role proposal. Mirrors smartadd.go's invocation
// (`claude -p … --output-format text`) and extracts the JSON object from the
// response. Any failure returns ErrClaudeUnavailable so the caller can fall
// back to manual entry.
func Propose(description string) (Proposal, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return Proposal{}, ErrClaudeUnavailable
	}
	out, err := exec.Command("claude", "-p", proposePrompt(description), "--output-format", "text").Output()
	if err != nil {
		return Proposal{}, ErrClaudeUnavailable
	}
	jsonStr := extractJSONObject(strings.TrimSpace(string(out)))
	if jsonStr == "" {
		return Proposal{}, ErrClaudeUnavailable
	}
	var p Proposal
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil || len(p.Roles) == 0 {
		return Proposal{}, ErrClaudeUnavailable
	}
	NormalizeFolders(&p)
	return p, nil
}

// NormalizeFolders fills in a slug folder for any role missing one.
func NormalizeFolders(p *Proposal) {
	for i := range p.Roles {
		if strings.TrimSpace(p.Roles[i].Folder) == "" {
			p.Roles[i].Folder = Slug(p.Roles[i].Name)
		} else {
			p.Roles[i].Folder = Slug(p.Roles[i].Folder)
		}
	}
}

// Slug lowercases and hyphenates a name into a safe folder/window token.
func Slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func proposePrompt(description string) string {
	return fmt.Sprintf(`You are helping bootstrap a multi-role "team" for a project — one Claude Code session per role, coordinating through a shared board.

From the project description below, propose the SMALLEST useful set of roles. Almost every team wants a "lead" (coordination) and a "board" is created automatically — do NOT include "board" as a role. Prefer common archetypes when they fit: lead, eng, marketing, sales, design, research, ops. Add project-specific roles only when clearly warranted.

Reply with ONLY a JSON object, no prose, in exactly this shape:
{
  "board_focus": "one sentence describing the team's current top priority",
  "roles": [
    { "name": "lead", "folder": "lead", "mission": "one sentence", "objectives": ["...", "..."] }
  ]
}

PROJECT DESCRIPTION:
%s`, description)
}

// extractJSONObject returns the first balanced {...} object in s (Claude may
// wrap it in prose or a markdown fence).
func extractJSONObject(s string) string {
	start := strings.Index(s, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	inStr := false
	esc := false
	for i := start; i < len(s); i++ {
		c := s[i]
		switch {
		case esc:
			esc = false
		case c == '\\' && inStr:
			esc = true
		case c == '"':
			inStr = !inStr
		case inStr:
			// skip
		case c == '{':
			depth++
		case c == '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}
