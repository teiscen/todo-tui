# todo-tui

A keyboard-driven todo dashboard for the terminal. Tasks, subtasks, recurrence,
time tracking, a calendar grid, and notes — in one multi-pane TUI backed by a
local SQLite database.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- **Task list** with horizons (small / long-term), colored labels, due dates, and
  collapsible groups.
- **Preview / editor pane** for task details: notes, subtasks, links, and time log.
- **Subtasks (steps)** with statuses (`todo` / `started` / `done`), step notes, and
  per-step due dates.
- **Recurring tasks** — daily, specific weekdays, or every N days — with a
  completion history.
- **Time tracking** — log minutes per task or step.
- **Planner** — a calendar grid with per-label entries.
- **Notes** pane.
- **Theming & layout** via a single YAML config (colors, panel sizes, labels).
- **Pure-Go SQLite** (no cgo) — trivially cross-compilable.

## Requirements

- Go 1.25+

## Install

```sh
git clone https://github.com/teiscen/todo-tui.git
cd todo-tui
go build -o todo-tui .
```

## Run

```sh
./todo-tui
# or
go run .
```

On first launch the SQLite database is created automatically and starts empty.

> **Note:** the app reads `./artifacts/config.yaml` and writes
> `./artifacts/db.sqlite3`, so it must be run from the repository root (or a
> directory containing an `artifacts/` folder).

## Configuration

`artifacts/config.yaml` controls:

- `layout` — outer dimensions (in terminal cells) of the tasks, preview, and
  editor panels.
- `colors` — text, border, weekday, and weekend colors (hex).
- `labels` — the default selected label, label ordering, and each label's display
  name + color.

## Keybindings

<!-- TODO: document the pane-focus and per-panel shortcuts once finalized. -->

## License

<!-- TODO: choose a license. -->
