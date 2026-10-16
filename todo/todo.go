package todo

import (
	"database/sql"
	"fmt"
	"time"

	"todo-tui/backend"
	"todo-tui/config"
	"todo-tui/database"
	"todo-tui/leftbox"
	"todo-tui/planner"
	"todo-tui/preview"
	"todo-tui/taskform"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// Todo is the root model. It coordinates the dashboard panels and owns
// persistence, so child views never write to the database directly.
type Todo struct {
	db  *sql.DB
	cfg config.Config
	err error

	plannerModel planner.Planner
	leftboxModel leftbox.Leftbox
	previewModel preview.Preview

	form       taskform.Form
	formActive bool
	focus      Pane
}

type Pane int

const (
	PaneTasks Pane = iota
	PanePreview
	PanePlanner
	paneCount
)

// NewTodo creates the dashboard, ensures configured labels exist, and loads saved tasks.
func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	theme := cfg.Theme()
	t := &Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: planner.New(theme, cfg.Labels, backend.NewDate(time.Now())),
		leftboxModel: leftbox.New(theme, cfg.Labels, cfg.Layout.Tasks.Width, cfg.Layout.Tasks.Height),
		previewModel: preview.New(theme, cfg.Layout.Preview.Width, cfg.Layout.Preview.Height),
		form:         taskform.New(theme, cfg.Labels, cfg.Layout.Editor.Width, cfg.Layout.Editor.Height),
		focus:        PaneTasks,
	}

	if err := database.SyncLabels(t.db, t.cfg.Labels); err != nil {
		t.err = err
	}
	t.reloadTasks()
	t.syncPreview()
	return t
}

// Init asks the planner to load its initially selected date and label off the
// main loop.
func (t *Todo) Init() tea.Cmd {
	return t.loadPlanner(*t.plannerModel.SelectedLabel(), *t.plannerModel.SelectedDate())
}

// View joins the task list, selected-task preview (or editor), and planner.
func (t *Todo) View() tea.View {
	// Render copies so focus styling does not alter the configured model colors.
	tasks := t.leftboxModel
	previewPanel := t.previewModel
	plannerPanel := t.plannerModel

	tasks.Style.Theme.Border = t.panelBorder(PaneTasks, string(tasks.Style.Theme.Border))
	previewPanel.Style.Theme.Border = t.panelBorder(PanePreview, string(previewPanel.Style.Theme.Border))
	plannerPanel.Style.Theme.Border = t.panelBorder(PanePlanner, string(plannerPanel.Style.Theme.Border))

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		withHelp(tasks.View(), t.focus == PaneTasks, "tab/shift+tab focus · ctrl+1/2/3 pane · j/k move", "space complete · enter group · a add · e edit · f filter · [/] tabs"),
		withHelp(previewPanel.View(), t.focus == PanePreview, "tab/shift+tab focus · ctrl+1/2/3 pane", "e edit selected task · a add task"),
		withHelp(plannerPanel.View(), t.focus == PanePlanner, "tab/shift+tab focus · ctrl+1/2/3 pane · h/j/k/l move date", "ctrl+w calendar/notes · ctrl+h/l label · ctrl+j/k month"),
	)
	view := tea.NewView(content)
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

// panelBorder preserves the configured active color and mutes inactive panels.
func (t *Todo) panelBorder(pane Pane, activeColor string) backend.HexCode {
	if t.focus == pane {
		return backend.HexCode(activeColor)
	}
	return backend.HexCode("#737994")
}

// withHelp places movement hints and panel-specific shortcuts below the panel,
// clipped to the panel's width so they never widen the column past its border
// and open a gap between panels.
func withHelp(content string, focused bool, movement, specific string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#737994"))
	if focused {
		style = style.Foreground(lipgloss.Color("#C6D0F5"))
	}
	width := lipgloss.Width(content)
	clip := func(s string) string {
		return style.Render(ansi.Truncate(s, width, "…"))
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, clip(movement), clip(specific))
}

// Update routes messages to the active panel and handles database-backed task actions.
func (t *Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ctrl+q exits from anywhere, including while the task editor is open.
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+q" {
		return t, tea.Quit
	}

	switch m := msg.(type) {
	case tea.MouseMotionMsg:
		if !t.formActive && m.X < t.leftboxModel.Style.Width && t.leftboxModel.HoverAt(m.Y) {
			t.focus = PaneTasks
			t.syncPreview()
		}
		return t, nil

	case taskform.CancelMsg:
		t.formActive = false
		t.syncPreview()
		return t, nil

	case taskform.SavedMsg:
		if _, err := database.SaveTask(t.db, dbTaskFromForm(m.Task)); err != nil {
			t.err = err
			t.form.Error = err.Error()
			t.previewModel.SetForm(&t.form)
			return t, nil
		}
		t.err = nil
		t.formActive = false
		t.reloadTasks()
		t.syncPreview()
		return t, nil

	case leftbox.ToggleTodoMsg:
		t.err = database.ToggleTask(t.db, m.ID, m.Recurring, m.Done)
		t.reloadTasks()
		t.syncPreview()
		return t, nil

	case planner.ChangeCalendarMsg:
		return t, t.loadPlanner(m.LabelID, m.NewDate)

	case planner.SaveEntryMsg:
		t.AddEntry(m.Label, m.Date, m.Msg)
		return t, t.loadPlanner(m.Label, m.Date)

	case planner.RemoveEntryMsg:
		t.RemoveEntry(m.Label, m.Date)
		return t, t.loadPlanner(m.Label, m.Date)

	case plannerLoadedMsg:
		// Drop stale results: the user may have moved on while the query ran.
		if m.LabelID != *t.plannerModel.SelectedLabel() || m.NewDate != *t.plannerModel.SelectedDate() {
			return t, nil
		}
		t.err = nil
		t.plannerModel.ReloadState(m.sDate, m.entries)
		return t, nil

	case errMsg:
		t.err = m.err
		return t, nil
	}

	// While a text field owns the keyboard, forward every key to it instead of
	// interpreting navigation shortcuts.
	if t.typing() {
		return t, t.forwardToFocused(msg)
	}

	// Global navigation and task shortcuts.
	if key, ok := msg.(tea.KeyPressMsg); ok {
		if handled, cmd := t.handleKey(key); handled {
			return t, cmd
		}
	}

	return t, t.forwardToFocused(msg)
}

// typing reports whether a text field (the task editor or the planner's notes)
// owns the keyboard.
func (t *Todo) typing() bool {
	return t.formActive || t.plannerModel.NotesModel.TextArea.Focused()
}

// handleKey applies pane navigation and task shortcuts. It reports false for
// keys that should be forwarded to the focused pane instead.
func (t *Todo) handleKey(key tea.KeyPressMsg) (bool, tea.Cmd) {
	switch key.String() {
	case "tab":
		t.focus = (t.focus + 1) % paneCount
		return true, nil
	case "shift+tab":
		t.focus = (t.focus + paneCount - 1) % paneCount
		return true, nil
	case "ctrl+1":
		t.focus = PaneTasks
		return true, nil
	case "ctrl+2":
		t.focus = PanePreview
		return true, nil
	case "ctrl+3":
		t.focus = PanePlanner
		return true, nil
	case "a":
		if t.focus == PaneTasks || t.focus == PanePreview {
			t.beginAdd()
			return true, nil
		}
		return false, nil
	case "e":
		if t.focus == PaneTasks || t.focus == PanePreview {
			t.beginEdit()
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

// forwardToFocused passes a message to whatever owns the keyboard: the editor
// while it is open, then the focused pane.
func (t *Todo) forwardToFocused(msg tea.Msg) tea.Cmd {
	if t.formActive {
		var cmd tea.Cmd
		t.form, cmd = t.form.Update(msg)
		t.previewModel.SetForm(&t.form)
		return cmd
	}
	var cmd tea.Cmd
	switch t.focus {
	case PaneTasks:
		t.leftboxModel, cmd = t.leftboxModel.Update(msg)
		t.syncPreview()
	case PanePlanner:
		t.plannerModel, cmd = t.plannerModel.Update(msg)
	}
	return cmd
}

// beginAdd opens a blank form, defaulting to the currently selected label.
func (t *Todo) beginAdd() {
	t.focus = PanePreview
	task := taskform.Task{
		Horizon: "small",
		LabelID: t.cfg.Labels.Selected,
	}
	t.form.Begin(task, false)
	t.formActive = true
	t.previewModel.SetForm(&t.form)
}

// beginEdit loads the selected task's persisted fields into the shared form.
func (t *Todo) beginEdit() {
	row, ok := t.leftboxModel.Selected()
	if !ok {
		return
	}

	tasks, err := database.FetchTasks(t.db)
	if err != nil {
		t.err = err
		return
	}
	for _, task := range tasks {
		if task.ID == row.ID {
			t.focus = PanePreview
			t.form.Begin(formTaskFromDB(task), true)
			t.formActive = true
			t.previewModel.SetForm(&t.form)
			return
		}
	}
}

// reloadTasks converts persisted tasks into display-ready rows for the left panel.
func (t *Todo) reloadTasks() {
	tasks, err := database.FetchTasks(t.db)
	if err != nil {
		t.err = err
		return
	}

	today := backend.NewDate(time.Now())
	rows := make([]leftbox.Row, 0, len(tasks))
	for _, task := range tasks {
		row, err := t.rowFromTask(task, today)
		if err != nil {
			t.err = err
			continue
		}
		rows = append(rows, row)
	}
	t.leftboxModel.ReloadState(rows)
}

// rowFromTask resolves task data into the grouping, label color, schedule, and due text used by the UI.
func (t *Todo) rowFromTask(task database.Task, today backend.Date) (leftbox.Row, error) {
	done, err := database.IsTaskDoneToday(t.db, task.ID, task.Recurring, task.DoneAt)
	if err != nil {
		return leftbox.Row{}, fmt.Errorf("check completion for task %d: %w", task.ID, err)
	}

	tab := leftbox.TabSmall
	if task.Horizon == "long" {
		tab = leftbox.TabLong
	}
	if task.Recurring {
		tab = leftbox.TabRecurring
	}

	row := leftbox.Row{
		ID:      task.ID,
		Tab:     tab,
		Title:   task.Title,
		Done:    done,
		LabelID: backend.LabelID(task.LabelID.String),
	}
	if info, ok := t.cfg.Labels.Info[row.LabelID]; ok {
		row.Color = info.Color
	}

	if task.Recurring {
		row.Schedule = leftbox.ScheduleText(task.Weekdays)
		switch task.Kind {
		case backend.RecurDaily:
			row.Schedule = "daily"
		case backend.RecurWeekdays:
			row.Schedule = leftbox.ScheduleText(task.Weekdays)
		case backend.RecurInterval:
			row.Schedule = fmt.Sprintf("every %dd", task.IntervalDays.Int64)
		}
		row.Group = leftbox.GroupToday
		if done {
			row.Group = leftbox.GroupDone
		}
		return row, nil
	}

	if task.DueDate.Valid && task.DueDate.String != "" {
		due, err := backend.ParseSQL(task.DueDate.String)
		if err != nil {
			return row, fmt.Errorf("parse due date for task %d: %w", task.ID, err)
		}
		row.Group = leftbox.GroupFor(&due, done, today)
		if !done {
			row.Due, row.DueTone = leftbox.DueInfo(due, today)
		}
	} else {
		row.Group = leftbox.GroupFor(nil, done, today)
	}
	return row, nil
}

// syncPreview keeps the preview panel consistent with the active form or selected task.
func (t *Todo) syncPreview() {
	if t.formActive {
		t.previewModel.SetForm(&t.form)
		return
	}
	if row, ok := t.leftboxModel.Selected(); ok {
		t.previewModel.SetTask(detailFromRow(row), true)
		return
	}
	t.previewModel.SetTask(preview.Detail{}, false)
}
