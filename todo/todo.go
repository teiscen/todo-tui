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

// Init asks the planner to load its initially selected date and label.
func (t *Todo) Init() tea.Cmd {
	return func() tea.Msg {
		return planner.ChangeCalendarMsg{
			LabelID: *t.plannerModel.SelectedLabel(),
			NewDate: *t.plannerModel.SelectedDate(),
		}
	}
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

// withHelp places movement hints and panel-specific shortcuts below the panel.
func withHelp(content string, focused bool, movement, specific string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#737994"))
	if focused {
		style = style.Foreground(lipgloss.Color("#C6D0F5"))
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, style.Render(movement), style.Render(specific))
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

	case leftbox.SelectTodoMsg:
		if !t.formActive {
			t.syncPreview()
		}
		return t, nil

	case planner.ChangeCalendarMsg:
		t.ReloadPlannerState(m.LabelID, m.NewDate)
		return t, nil

	case planner.SaveEntryMsg:
		t.AddEntry(m.Label, m.Date, m.Msg)
		t.ReloadPlannerState(m.Label, m.Date)
		return t, nil

	case planner.RemoveEntryMsg:
		t.RemoveEntry(m.Label, m.Date)
		t.ReloadPlannerState(m.Label, m.Date)
		return t, nil
	}

	// While editing, reserve keyboard input for the form instead of the other panels.
	if t.formActive {
		var cmd tea.Cmd
		t.form, cmd = t.form.Update(msg)
		t.previewModel.SetForm(&t.form)
		return t, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "tab":
			t.focus = (t.focus + 1) % paneCount
			return t, nil
		case "shift+tab":
			t.focus = (t.focus + paneCount - 1) % paneCount
			return t, nil
		case "ctrl+1":
			t.focus = PaneTasks
			return t, nil
		case "ctrl+2":
			t.focus = PanePreview
			return t, nil
		case "ctrl+3":
			t.focus = PanePlanner
			return t, nil
		case "a":
			if t.focus == PaneTasks || t.focus == PanePreview {
				t.beginAdd()
				return t, nil
			}
		case "e":
			if t.focus == PaneTasks || t.focus == PanePreview {
				t.beginEdit()
				return t, nil
			}
		}
		var cmd tea.Cmd
		switch t.focus {
		case PaneTasks:
			t.leftboxModel, cmd = t.leftboxModel.Update(msg)
			t.syncPreview()
		case PanePlanner:
			t.plannerModel, cmd = t.plannerModel.Update(msg)
		}
		return t, cmd
	}

	var cmd tea.Cmd
	if t.focus == PanePlanner {
		t.plannerModel, cmd = t.plannerModel.Update(msg)
	}
	return t, cmd
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
