package todo

import (
	"database/sql"

	"todo-tui/config"
	"todo-tui/leftbox"
	"todo-tui/planner"
	"todo-tui/preview"
	"todo-tui/weekly"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Todo struct {
	db  *sql.DB
	cfg config.Config

	plannerModel planner.Planner

	leftboxModel leftbox.Leftbox
	weeklyModel  weekly.Weekly
	previewModel preview.Preview

	Style TodoStyle

	// err is the last database error. Not displayed anywhere yet.
}

func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	return &Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
		leftboxModel: cfg.LoadLeftboxModel(),
		weeklyModel:  cfg.LoadWeeklyModel(),
		previewModel: cfg.LoadPreviewModel(),
	}
}

func (t *Todo) Init() tea.Cmd {
	// return t.plannerModel.Init()
	return func() tea.Msg {
		return planner.ChangeCalendarMsg{
			LabelID: *t.plannerModel.SelectedLabel(),
			NewDate: *t.plannerModel.SelectedDate(),
		}
	}
}

func (t *Todo) View() tea.View {
	left := t.leftboxModel.View().Content
	middle := lipgloss.JoinVertical(
		lipgloss.Left,
		t.previewModel.View().Content,
		t.weeklyModel.View().Content,
	)
	right := t.plannerModel.View().Content
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left, middle, right,
	)
	return tea.NewView(combined)
}

func (t *Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case planner.ChangeCalendarMsg:
		l := m.LabelID
		d := m.NewDate
		t.ReloadPlannerState(l, d)
		return t, nil
	case planner.SaveEntryMsg:
		l := m.Label
		d := m.Date
		msg := m.Msg
		t.AddEntry(l, d, msg)
		t.ReloadPlannerState(l, d)
		return t, nil
	case planner.RemoveEntryMsg:
		l := m.Label
		d := m.Date
		t.RemoveEntry(l, d)
		t.ReloadPlannerState(l, d)
		return t, nil
	}

	var cmd tea.Cmd
	t.plannerModel, cmd = t.plannerModel.Update(msg)
	return t, cmd
}
