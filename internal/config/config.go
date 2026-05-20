package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var configDirOverride string

// SetConfigDir overrides the config directory. Call with "" to reset. For testing only.
func SetConfigDir(dir string) {
	configDirOverride = dir
}

// WorkspaceEntry holds a directory path and its associated tags.
type WorkspaceEntry struct {
	Dir  string   `json:"dir"`
	Tags []string `json:"tags,omitempty"`
}

// Config holds the persisted user configuration.
type Config struct {
	Workspaces []WorkspaceEntry `json:"workspaces"`
}

func configPath() (string, error) {
	if configDirOverride != "" {
		return filepath.Join(configDirOverride, "config.json"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "terraview", "config.json"), nil
}

// Load reads config from disk. Migrates legacy "directories" format automatically.
// Returns empty Config if file doesn't exist.
func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var raw struct {
		Workspaces  []WorkspaceEntry `json:"workspaces"`
		Directories []string         `json:"directories"` // legacy
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{}, err
	}
	if len(raw.Workspaces) > 0 {
		return Config{Workspaces: raw.Workspaces}, nil
	}
	ws := make([]WorkspaceEntry, len(raw.Directories))
	for i, d := range raw.Directories {
		ws[i] = WorkspaceEntry{Dir: d}
	}
	return Config{Workspaces: ws}, nil
}

// Save writes config to disk, creating directories as needed.
func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
