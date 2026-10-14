package todo

import (
	"database/sql"

	"todo-tui/config"
	"todo-tui/leftbox"
	"todo-tui/planner"
	"todo-tui/preview"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Todo struct {
	db  *sql.DB
	cfg config.Config
	err error

	plannerModel planner.Planner

	leftboxModel leftbox.Leftbox
	previewModel preview.Preview

	Style TodoStyle
}

func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	t := &Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
		leftboxModel: cfg.LoadLeftboxModel(),
		previewModel: cfg.LoadPreviewModel(),
	}
	if row, ok := t.leftboxModel.Selected(); ok {
		t.previewModel.SetTask(row, true)
	}
	return t
}

func (t *Todo) Init() tea.Cmd {
	return func() tea.Msg {
		return planner.ChangeCalendarMsg{
			LabelID: *t.plannerModel.SelectedLabel(),
			NewDate: *t.plannerModel.SelectedDate(),
		}
	}
}

func (t *Todo) View() tea.View {
	left := t.leftboxModel.View().Content
	middle := t.previewModel.View().Content
	right := t.plannerModel.View().Content
	combined := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)
	view := tea.NewView(combined)
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

func (t *Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.MouseMotionMsg:
		mouse := m.Mouse()
		if mouse.X < t.leftboxModel.Style.Width && t.leftboxModel.HoverAt(mouse.Y) {
			if row, ok := t.leftboxModel.Selected(); ok {
				t.previewModel.SetTask(row, true)
			}
		}
		return t, nil
	case leftbox.SelectTodoMsg:
		if row, ok := t.leftboxModel.Selected(); ok {
			t.previewModel.SetTask(row, true)
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

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "j", "down", "k", "up", "[", "]", "f":
			var cmd tea.Cmd
			t.leftboxModel, cmd = t.leftboxModel.Update(msg)
			if row, ok := t.leftboxModel.Selected(); ok {
				t.previewModel.SetTask(row, true)
			} else {
				t.previewModel.SetTask(leftbox.Row{}, false)
			}
			return t, cmd
		}
	}

	var cmd tea.Cmd
	t.plannerModel, cmd = t.plannerModel.Update(msg)
	return t, cmd
}
