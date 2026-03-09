package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ProfileStore holds named Claude config directory paths.
type ProfileStore struct {
	Profiles map[string]string `json:"profiles"`
}

// ProfilesPath returns the path to ~/.bmad/profiles.json.
func ProfilesPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bmad", "profiles.json")
}

// ReadProfiles reads profiles.json, seeding defaults on first run.
func ReadProfiles() (ProfileStore, error) {
	store := ProfileStore{Profiles: make(map[string]string)}
	path := ProfilesPath()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		home, _ := os.UserHomeDir()
		if _, e := os.Stat(filepath.Join(home, ".claude")); e == nil {
			store.Profiles["perso"] = "~/.claude"
		}
		if _, e := os.Stat(filepath.Join(home, ".claude-work")); e == nil {
			store.Profiles["work"] = "~/.claude-work"
		}
		return store, nil
	}
	if err != nil {
		return store, err
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return store, err
	}
	if store.Profiles == nil {
		store.Profiles = make(map[string]string)
	}
	return store, nil
}

// WriteProfiles writes the store to ~/.bmad/profiles.json.
func WriteProfiles(store ProfileStore) error {
	path := ProfilesPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
