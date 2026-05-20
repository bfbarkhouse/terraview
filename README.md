# terraview
### A terminal UI for your Terraform state files

Never forget a deployed resource again. terraview shows all your Terraform state files at a glance — green means resources are tracked, red means empty. Add state files with `terraview add .`, bulk-scan existing projects, browse resources, tag and filter, and drop into a shell in any directory, all without leaving your terminal.

## Demo

![terraview demo](demo-v2.gif)

## Installation

```bash
git clone https://github.com/bfbarkhouse/terraview
cd terraview
go install .
```

Requires Go 1.24+.

## Usage

```bash
terraview          # launch the TUI
terraview --help   # show help
```

On first launch the state file list is empty. Press `e` to add your Terraform directories, or use `terraview add` from the command line:

```bash
terraview add .                              # add the current directory
terraview add ~/infra/prod                   # add a specific directory
terraview add ~/infra/prod --tags prod,aws   # add with tags
```

If the directory is already in the config, terraview prints a warning and exits cleanly.

## Features

### State file list + resource browser

The default view splits the screen: state files on the left, resources for the selected state file on the right. Use `tab` to move focus between the two panels, and arrow keys (or `j`/`k`) to navigate within each.

A green dot next to a state file means a `terraform.tfstate` file was found with at least one resource. A red dot means no state file or no resources.

### Summary view

Press `s` to switch to a card grid showing all state files at a glance. Cards wrap to fill the window width. Use arrow keys to navigate, `enter` to jump to a state file in split view, and `s` again to go back.

### Adding and editing state files

Press `e` to open the state file editor. From there:

- **Add a single directory** — press `enter`, type a path (tab-expands `~`), confirm with `enter`.
- **Recursive scan** — press `enter`, type a root path, then press `tab` to toggle to recursive mode. terraview walks the directory tree and adds every directory containing a `terraform.tfstate` file as a separate state file.
- **Delete a state file** — navigate to it and press `d`.
- **Save** — `ctrl+s`. Changes are written to `~/.config/terraview/config.json`.

### Tags

Tags let you group and filter state files. Press `t` on any selected state file to open the tag editor:

- `enter` or `tab` — add a new tag (comma-separated for multiple)
- `d` — remove the highlighted tag
- `ctrl+s` — save

The tag input offers prefix-match autocomplete from all tags currently in use across your state files. Press `tab` to accept a suggestion.

### Filtering

Press `f` to open the filter bar. Type a tag and press `enter` to add it as an active filter. Multiple filters are ANDed — only state files matching all active tags are shown. Press `backspace` on an empty input to remove the last filter. Press `esc` to close the filter bar (active filters remain in effect; the header shows a count).

### Open a shell

Press `o` to suspend terraview and open `$SHELL` with its working directory set to the selected state file's directory. When you `exit` the shell, terraview resumes exactly where you left off.

### Refresh

Press `r` to re-read the `terraform.tfstate` file for the currently selected state file and update the resource list.

## Key reference

| Key | Action |
|-----|--------|
| `↑` / `k`, `↓` / `j` | Navigate up/down |
| `←` / `h`, `→` / `l` | Navigate left/right (summary view) |
| `tab` | Switch focus between state file and resource panels |
| `s` | Toggle summary view |
| `e` | Open state file editor |
| `t` | Edit tags for selected state file |
| `f` | Open filter bar |
| `r` | Refresh selected state file |
| `o` | Open shell in selected directory |
| `enter` | Select / confirm |
| `esc` | Cancel / close modal |
| `q` / `ctrl+c` | Quit |

## Configuration

Config is stored at `~/.config/terraview/config.json` (or `$XDG_CONFIG_HOME/terraview/config.json` where set). It is written automatically by the state file editor — you should not need to edit it by hand. Format:

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
