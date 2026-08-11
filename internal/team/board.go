package team

import (
	"fmt"
	"strings"
)

// RoleSpec is a fully-resolved role (post user-confirmation): its name,
// folder, one-line mission, and objectives. Produced by Propose (Claude) or
// the huh fallback, consumed by scaffolding.
type RoleSpec struct {
	Name       string   `json:"name"`
	Folder     string   `json:"folder"`
	Mission    string   `json:"mission"`
	Objectives []string `json:"objectives"`
}

// RenderBoard builds the initial BOARD.md: shared coordination state with a
// Current Focus, one section per role, and the Cross-Role Requests protocol
// (the `@role:` + Definition-of-Done convention the hooks key off).
func RenderBoard(project, focus string, roles []RoleSpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — Team Board\n\n", project)
	b.WriteString("Shared coordination state. Each role reads this **before acting** and appends updates under its own section. This file is the baton between panes — the only channel between sessions.\n\n")

	b.WriteString("## Current Focus\n_(set by the lead)_\n")
	if strings.TrimSpace(focus) != "" {
		fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(focus))
	}
	b.WriteString("\n## Roles\n\n")
	for _, r := range roles {
		fmt.Fprintf(&b, "### %s\n- Status: —\n- Last update: —\n", r.Name)
		if strings.TrimSpace(r.Mission) != "" {
			fmt.Fprintf(&b, "- Mission: %s\n", strings.TrimSpace(r.Mission))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Cross-Role Requests\n")
	b.WriteString("_Use `@<role>: request`. The target role's pane is poked live when the line is written, and re-catches any un-`[DONE]` `@<role>:` lines when it boots._\n")
	b.WriteString("_**Definition of Done:** every `@<role>:` request states its acceptance checks + the proof it must produce. The role does the work and **appends its proof** under the request. Only the lead writes `[DONE]` on a line, after verifying the proof. `[DONE]` stops the poke/rescan from re-triggering it._\n\n")
	// Example uses the literal `@<role>:` placeholder, never a real role
	// name — otherwise the board's own example would self-trigger the
	// notify/rescan hooks.
	b.WriteString("<!-- Example: @<role>: build the login API. DoD: passing test + curl output. -->\n\n")

	b.WriteString("## Next Actions\n_(each role appends concrete follow-ups here)_\n")
	return b.String()
}

// RenderResourceRequests builds the companion RESOURCE-REQUESTS.md — the
// place a role asks the lead/founder to approve spend, credentials, or
// writes to shared systems before acting on them.
func RenderResourceRequests(project string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — Resource Requests\n\n", project)
	b.WriteString("Requests for anything a role cannot grant itself: budget, API keys, access to shared systems, third-party accounts. The lead/founder reviews and writes `APPROVED` before the role acts.\n\n")
	b.WriteString("## Open\n_(append `- [ ] <role>: <request> — <why>` lines here)_\n\n")
	b.WriteString("## Approved\n_(moved here with `APPROVED <date>` once granted)_\n")
	return b.String()
}
