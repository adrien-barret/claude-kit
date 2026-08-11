package team

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// archetypeRel is the template subdir (under the resolved template dir)
// holding one folder per role archetype, plus a `generic` fallback.
const archetypeRel = "team-roles"

// RoleTemplateData is the data passed to a role archetype template.
type RoleTemplateData struct {
	Project    string
	Role       string
	Mission    string
	Objectives []string
	// BoardRel is the board path relative to this role's folder (e.g.
	// ../board/BOARD.md), so the rendered CLAUDE.md links the real file.
	BoardRel string
	// Roles is every role name on the team (for cross-role awareness).
	Roles []string
}

// RenderRole renders a role's CLAUDE.md from its archetype template. It
// looks for <templateDir>/team-roles/<role>/CLAUDE.md.tmpl and falls back
// to team-roles/generic/CLAUDE.md.tmpl for unknown roles.
func RenderRole(templateDir string, data RoleTemplateData) (string, error) {
	path := filepath.Join(templateDir, archetypeRel, data.Role, "CLAUDE.md.tmpl")
	if _, err := os.Stat(path); err != nil {
		path = filepath.Join(templateDir, archetypeRel, "generic", "CLAUDE.md.tmpl")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("role").Parse(string(raw))
	if err != nil {
		return "", err
	}
	var out strings.Builder
	if err := tmpl.Execute(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

// BoardRelFrom computes the board path relative to a role's folder, in
// forward-slash form suitable for a Markdown link.
func BoardRelFrom(roleFolder, boardRel string) string {
	rel, err := filepath.Rel(roleFolder, filepath.FromSlash(boardRel))
	if err != nil {
		return boardRel
	}
	return filepath.ToSlash(rel)
}
