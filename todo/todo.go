package todo

import (
	"database/sql"

	"todo-tui/config"
	"todo-tui/planner"

	tea "charm.land/bubbletea/v2"
)

type Todo struct {
	db           *sql.DB
	cfg          config.Config
	plannerModel planner.Planner

	Style TodoStyle

	// err is the last database error. Not displayed anywhere yet.
}

func NewTodo(cfg config.Config, db *sql.DB) *Todo {
	return &Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
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
	return t.plannerModel.View()
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
