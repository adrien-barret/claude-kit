package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultTemplateDirName = "templates"
	BmadDirName            = ".bmad"
	ClaudeDirName          = ".claude"
)

// TemplateDir resolves the template directory using this priority:
//  1. $BMAD_TEMPLATE_DIR environment variable
//  2. ~/.bmad/templates/
//  3. Adjacent project-template/ (for development)
func TemplateDir() string {
	if dir := os.Getenv("BMAD_TEMPLATE_DIR"); dir != "" {
		return dir
	}

	home, err := os.UserHomeDir()
	if err == nil {
		installed := filepath.Join(home, BmadDirName, DefaultTemplateDirName)
		if info, err := os.Stat(installed); err == nil && info.IsDir() {
			return installed
		}
	}

	// Dev mode: look relative to the binary
	exe, err := os.Executable()
	if err == nil {
		adjacent := filepath.Join(filepath.Dir(exe), "project-template", ClaudeDirName)
		if info, err := os.Stat(adjacent); err == nil && info.IsDir() {
			return adjacent
		}
	}

	// Fallback: look relative to working directory
	if cwd, err := os.Getwd(); err == nil {
		adjacent := filepath.Join(cwd, "project-template", ClaudeDirName)
		if info, err := os.Stat(adjacent); err == nil && info.IsDir() {
			return adjacent
		}
	}

	// Default path even if it doesn't exist yet
	if home != "" {
		return filepath.Join(home, BmadDirName, DefaultTemplateDirName)
	}
	return filepath.Join("~", BmadDirName, DefaultTemplateDirName)
}

// InstalledTemplatesDir returns ~/.bmad/templates/.
func InstalledTemplatesDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, BmadDirName, DefaultTemplateDirName)
}

// IsWindows returns true if running on Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
