package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
)

// ResolvePath converts p to an absolute path, expanding ~ and resolving relative paths.
func ResolvePath(p string) (string, error) {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if p == "~" {
			return home, nil
		}
		p = filepath.Join(home, p[2:])
	}
	return filepath.Abs(p)
}

// RunAdd implements the `terraview add <directory> [--tags ...]` subcommand.
func RunAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	tagsFlag := fs.String("tags", "", "comma-separated tags to assign to the workspace")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: terraview add <directory> [--tags tag1,tag2]")
	}

	resolved, err := ResolvePath(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("directory does not exist: %s", resolved)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", resolved)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	for _, ws := range cfg.Workspaces {
		if ws.Dir == resolved {
			fmt.Fprintf(os.Stderr, "warning: %s is already in the config\n", resolved)
			return nil
		}
	}

	var tags []string
	if *tagsFlag != "" {
		for _, t := range strings.Split(*tagsFlag, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	cfg.Workspaces = append(cfg.Workspaces, config.WorkspaceEntry{Dir: resolved, Tags: tags})
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	home, _ := os.UserHomeDir()
	display := resolved
	if home != "" && (resolved == home || strings.HasPrefix(resolved, home+string(os.PathSeparator))) {
		display = "~" + resolved[len(home):]
	}
	fmt.Printf("Added workspace: %s\n", display)
	return nil
}
