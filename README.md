# terraview

A terminal UI for monitoring Terraform workspaces. Browse deployed resources across multiple state files, tag and filter workspaces, and drop into a shell in any directory — all without leaving your terminal.

## Screenshot

```
terraview  tab=focus  s=summary  e=edit  t=tags  f=filter  r=refresh  o=shell  q=quit

╭─────────────────────────╮ ╭──────────────────────────────────────────────────────╮
│ WORKSPACES              │ │ ~/infra/prod — 12 resources                          │
│                         │ │                                                      │
│ ▶ ● ~/infra/prod        │ │ ● aws_instance.web                                   │
│   prod aws              │ │ ● aws_instance.worker                                │
│   ● ~/infra/staging     │ │ ● aws_s3_bucket.assets                               │
│   staging aws           │ │ ● aws_rds_instance.main                              │
│   ● ~/infra/dev         │ │ ● aws_security_group.web                             │
│   dev                   │ │ ● aws_vpc.main                                       │
│                         │ │                                                      │
│ e = add/edit            │ │                                                      │
╰─────────────────────────╯ ╰──────────────────────────────────────────────────────╯
```

## Installation

```bash
git clone https://github.com/bfbarkhouse-redpanda/terraview
cd terraview
go install .
```

Requires Go 1.24+.

## Usage

```bash
terraview
```

On first launch the workspace list is empty. Press `e` to add your Terraform directories.

## Features

### Workspace list + resource browser

The default view splits the screen: workspaces on the left, resources for the selected workspace on the right. Use `tab` to move focus between the two panels, and arrow keys (or `j`/`k`) to navigate within each.

A green dot next to a workspace means a `terraform.tfstate` file was found with at least one resource. A red dot means no state file or no resources.

### Summary view

Press `s` to switch to a card grid showing all workspaces at a glance. Cards wrap to fill the window width. Use arrow keys to navigate, `enter` to jump to a workspace in split view, and `s` again to go back.

### Adding and editing workspaces

Press `e` to open the workspace editor. From there:

- **Add a single directory** — press `enter`, type a path (tab-expands `~`), confirm with `enter`.
- **Recursive scan** — press `enter`, type a root path, then press `tab` to toggle to recursive mode. terraview walks the directory tree and adds every folder containing a `terraform.tfstate` file as a separate workspace.
- **Delete a workspace** — navigate to it and press `d`.
- **Save** — `ctrl+s`. Changes are written to `~/.config/terraview/config.json`.

### Tags

Tags let you group and filter workspaces. Press `t` on any selected workspace to open the tag editor:

- `enter` or `tab` — add a new tag (comma-separated for multiple)
- `d` — remove the highlighted tag
- `ctrl+s` — save

The tag input offers prefix-match autocomplete from all tags currently in use across your workspaces. Press `tab` to accept a suggestion.

### Filtering

Press `f` to open the filter bar. Type a tag and press `enter` to add it as an active filter. Multiple filters are ANDed — only workspaces matching all active tags are shown. Press `backspace` on an empty input to remove the last filter. Press `esc` to close the filter bar (active filters remain in effect; the header shows a count).

### Open a shell

Press `o` to suspend terraview and open `$SHELL` with its working directory set to the selected workspace. When you `exit` the shell, terraview resumes exactly where you left off.

### Refresh

Press `r` to re-read the `terraform.tfstate` file for the currently selected workspace and update the resource list.

## Key reference

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Navigate up/down |
| `←` / `h`, `→` / `l` | Navigate left/right (summary view) |
| `tab` | Switch focus between workspace and resource panels |
| `s` | Toggle summary view |
| `e` | Open workspace editor |
| `t` | Edit tags for selected workspace |
| `f` | Open filter bar |
| `r` | Refresh selected workspace state |
| `o` | Open shell in selected workspace directory |
| `enter` | Select / confirm |
| `esc` | Cancel / close modal |
| `q` / `ctrl+c` | Quit |

## Configuration

Config is stored at `~/.config/terraview/config.json` (or `$XDG_CONFIG_HOME/terraview/config.json` where set). It is written automatically by the workspace editor — you should not need to edit it by hand. Format:

```json
{
  "workspaces": [
    { "dir": "~/infra/prod", "tags": ["prod", "aws"] },
    { "dir": "~/infra/staging", "tags": ["staging", "aws"] }
  ]
}
```

Legacy configs using a `"directories"` array are migrated automatically on first load.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling and layout
