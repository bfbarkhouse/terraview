package main

import (
	"fmt"
	"os"

	"github.com/bfbarkhouse-redpanda/terraview/internal/cmd"
	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
	"github.com/bfbarkhouse-redpanda/terraview/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

const usage = `terraview — terminal UI for monitoring Terraform workspaces

Usage:
  terraview                         launch the TUI
  terraview add <dir> [--tags ...]  add a workspace to the config

Flags:
  -h, --help  show this help

TUI keys:
  e        open workspace editor
  t        edit tags for selected workspace
  f        filter by tag
  s        toggle summary view
  r        refresh selected workspace
  o        open shell in workspace directory
  q        quit
`

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		fmt.Print(usage)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "add" {
		if err := cmd.RunAdd(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	app := tui.New(cfg)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
