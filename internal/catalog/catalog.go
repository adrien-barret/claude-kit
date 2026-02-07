package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Component represents a single template component (agent, skill, command, rule).
type Component struct {
	Type        string // "agents", "skills", "commands", "rules"
	Name        string // e.g. "backend", "security/pentest-web"
	Description string // extracted from YAML frontmatter
	Path        string // absolute path in template dir
}

// Category groups components by type.
type Category struct {
	Name       string
	Components []Component
}

// ScanTemplate scans the template directory and returns categorized components.
func ScanTemplate(templateDir string) ([]Category, error) {
	if _, err := os.Stat(templateDir); err != nil {
		return nil, fmt.Errorf("template directory not found: %s", templateDir)
	}

	types := []string{"agents", "skills", "commands", "rules"}
	var categories []Category

	for _, t := range types {
		dir := filepath.Join(templateDir, t)
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		var components []Component
		switch t {
		case "skills":
			components = scanSkills(dir, t)
		case "agents", "commands", "rules":
			components = scanMarkdownDir(dir, t)
		}

		if len(components) > 0 {
			sort.Slice(components, func(i, j int) bool {
				return components[i].Name < components[j].Name
			})
			categories = append(categories, Category{Name: t, Components: components})
		}
	}

	return categories, nil
}

// scanSkills handles the nested skill directory structure.
// Skills can be flat (code-reviewer/) or nested (security/pentest-web/).
func scanSkills(dir, typeName string) []Component {
	var components []Component

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		skillDir := filepath.Join(dir, name)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillFile); err == nil {
			// Direct skill directory (e.g. code-reviewer/SKILL.md)
			desc := ExtractDescription(skillFile)
			components = append(components, Component{
				Type:        typeName,
				Name:        name,
				Description: desc,
				Path:        skillDir,
			})
		}

		// Check for nested sub-skills (e.g. security/pentest-web/)
		subEntries, err := os.ReadDir(skillDir)
		if err != nil {
			continue
		}
		for _, sub := range subEntries {
			if !sub.IsDir() {
				continue
			}
			subSkillFile := filepath.Join(skillDir, sub.Name(), "SKILL.md")
			if _, err := os.Stat(subSkillFile); err == nil {
				subName := name + "/" + sub.Name()
				desc := ExtractDescription(subSkillFile)
				components = append(components, Component{
					Type:        typeName,
					Name:        subName,
					Description: desc,
					Path:        filepath.Join(skillDir, sub.Name()),
				})
			}
		}
	}

	return components
}

// scanMarkdownDir scans a directory of .md files.
func scanMarkdownDir(dir, typeName string) []Component {
	var components []Component

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		filePath := filepath.Join(dir, entry.Name())
		desc := ExtractDescription(filePath)

		components = append(components, Component{
			Type:        typeName,
			Name:        name,
			Description: desc,
			Path:        filePath,
		})
	}

	return components
}

// ExtractDescription reads the YAML frontmatter description from a file.
func ExtractDescription(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if inFrontmatter {
				break // end of frontmatter
			}
			inFrontmatter = true
			continue
		}

		if inFrontmatter && strings.HasPrefix(line, "description:") {
			desc := strings.TrimPrefix(line, "description:")
			desc = strings.TrimSpace(desc)
			// Remove surrounding quotes if present
			if len(desc) >= 2 && ((desc[0] == '"' && desc[len(desc)-1] == '"') || (desc[0] == '\'' && desc[len(desc)-1] == '\'')) {
				desc = desc[1 : len(desc)-1]
			}
			return desc
		}
	}

	return ""
}

// ExtractSkillDeps reads an agent file's frontmatter and returns its skills list.
func ExtractSkillDeps(agentPath string) []string {
	f, err := os.Open(agentPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	inSkills := false
	var skills []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}

		if !inFrontmatter {
			continue
		}

		// Check for skills key
		if strings.HasPrefix(line, "skills:") {
			rest := strings.TrimPrefix(line, "skills:")
			rest = strings.TrimSpace(rest)
			if rest != "" && rest != "|" {
				// Inline list: skills: skill1, skill2
				for _, s := range strings.Split(rest, ",") {
					s = strings.TrimSpace(s)
					if s != "" {
						skills = append(skills, s)
					}
				}
				inSkills = false
			} else {
				inSkills = true
			}
			continue
		}

		if inSkills {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- ") {
				skill := strings.TrimPrefix(trimmed, "- ")
				skill = strings.TrimSpace(skill)
				if skill != "" {
					skills = append(skills, skill)
				}
			} else if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				// End of skills list
				inSkills = false
			}
		}
	}

	return skills
}

// GetInstalled returns components installed in the target directory.
func GetInstalled(targetDir string) ([]Category, error) {
	if _, err := os.Stat(targetDir); err != nil {
		return nil, nil
	}

	types := []string{"agents", "skills", "commands", "rules"}
	var categories []Category

	for _, t := range types {
		dir := filepath.Join(targetDir, t)
		if _, err := os.Stat(dir); err != nil {
			continue
		}

		var components []Component
		switch t {
		case "skills":
			components = scanSkills(dir, t)
		case "agents", "commands", "rules":
			components = scanMarkdownDir(dir, t)
		}

		if len(components) > 0 {
			sort.Slice(components, func(i, j int) bool {
				return components[i].Name < components[j].Name
			})
			categories = append(categories, Category{Name: t, Components: components})
		}
	}

	return categories, nil
}

// CopyComponent copies a component from template to target directory.
func CopyComponent(templateDir, targetDir, compType, name string) error {
	switch compType {
	case "skills":
		return copySkill(templateDir, targetDir, name)
	case "agents", "commands", "rules":
		return copyMarkdown(templateDir, targetDir, compType, name)
	default:
		return fmt.Errorf("unknown component type: %s", compType)
	}
}

func copySkill(templateDir, targetDir, name string) error {
	srcDir := filepath.Join(templateDir, "skills", name)
	if _, err := os.Stat(srcDir); err != nil {
		return fmt.Errorf("skill not found: %s", name)
	}

	dstDir := filepath.Join(targetDir, "skills", name)
	if err := os.MkdirAll(filepath.Dir(dstDir), 0o755); err != nil {
		return err
	}

	return copyDir(srcDir, dstDir)
}

func copyMarkdown(templateDir, targetDir, compType, name string) error {
	srcFile := filepath.Join(templateDir, compType, name+".md")
	if _, err := os.Stat(srcFile); err != nil {
		return fmt.Errorf("%s not found: %s", compType, name)
	}

	dstDir := filepath.Join(targetDir, compType)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	return copyFile(srcFile, filepath.Join(dstDir, name+".md"))
}

// RemoveComponent removes a component from the target directory.
func RemoveComponent(targetDir, compType, name string) error {
	switch compType {
	case "skills":
		dir := filepath.Join(targetDir, "skills", name)
		return os.RemoveAll(dir)
	case "agents", "commands", "rules":
		file := filepath.Join(targetDir, compType, name+".md")
		return os.Remove(file)
	default:
		return fmt.Errorf("unknown component type: %s", compType)
	}
}

// IsInstalled checks if a specific component is installed.
func IsInstalled(targetDir, compType, name string) bool {
	switch compType {
	case "skills":
		skillFile := filepath.Join(targetDir, "skills", name, "SKILL.md")
		_, err := os.Stat(skillFile)
		return err == nil
	case "agents", "commands", "rules":
		file := filepath.Join(targetDir, compType, name+".md")
		_, err := os.Stat(file)
		return err == nil
	}
	return false
}

// FindReferencingAgents returns agent names that reference the given skill.
func FindReferencingAgents(targetDir, skillName string) []string {
	agentsDir := filepath.Join(targetDir, "agents")
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil
	}

	var refs []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		agentPath := filepath.Join(agentsDir, entry.Name())
		deps := ExtractSkillDeps(agentPath)
		for _, d := range deps {
			if d == skillName {
				refs = append(refs, strings.TrimSuffix(entry.Name(), ".md"))
				break
			}
		}
	}
	return refs
}

// CopyBaseFiles copies CLAUDE.md and settings.json from template to target.
func CopyBaseFiles(templateDir, targetDir string) error {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	// Copy CLAUDE.md (lives one level up from the .claude/ dir in template)
	claudeMd := filepath.Join(templateDir, "CLAUDE.md")
	if _, err := os.Stat(claudeMd); err == nil {
		if err := copyFile(claudeMd, filepath.Join(targetDir, "CLAUDE.md")); err != nil {
			return err
		}
	}

	// Copy settings.json
	settingsJson := filepath.Join(templateDir, "settings.json")
	if _, err := os.Stat(settingsJson); err == nil {
		if err := copyFile(settingsJson, filepath.Join(targetDir, "settings.json")); err != nil {
			return err
		}
	}

	return nil
}

// ReadSettingsTeammateMode reads settings.json and returns the current
// teammateMode value (defaults to "auto" if absent).
func ReadSettingsTeammateMode(targetDir string) (string, error) {
	path := filepath.Join(targetDir, "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading settings.json: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return "", fmt.Errorf("parsing settings.json: %w", err)
	}

	if mode, ok := settings["teammateMode"].(string); ok && mode != "" {
		return mode, nil
	}
	return "auto", nil
}

// PatchSettingsTeammateMode reads the installed settings.json and sets the
// teammateMode field to the given value.
func PatchSettingsTeammateMode(targetDir, mode string) error {
	path := filepath.Join(targetDir, "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading settings.json: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("parsing settings.json: %w", err)
	}

	settings["teammateMode"] = mode

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling settings.json: %w", err)
	}
	out = append(out, '\n')

	return os.WriteFile(path, out, 0o644)
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
