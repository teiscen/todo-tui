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

	plannerModel planner.Planner

	leftboxModel leftbox.Leftbox
	previewModel preview.Preview

	Style TodoStyle

	// err is the last database error. Not displayed anywhere yet.
}

func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	t := &Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
		leftboxModel: cfg.LoadLeftboxModel(),
		previewModel: cfg.LoadPreviewModel(),
	}
	t.previewModel.SetRows(t.leftboxModel.Rows)
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
	return tea.NewView(combined)
}

func (t *Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
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

	var cmd tea.Cmd
	t.plannerModel, cmd = t.plannerModel.Update(msg)
	return t, cmd
}
