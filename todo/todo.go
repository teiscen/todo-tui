package todo

import (
	"database/sql"

	"todo-tui/config"
	"todo-tui/database"
	"todo-tui/planner"

	tea "charm.land/bubbletea/v2"
)

type Todo struct {
	db           *sql.DB
	cfg          config.Config
	plannerModel planner.Planner

	Style TodoStyle
}

func NewTodo(cfgPath string, db *sql.DB) Todo {
	cfg, _ := config.LoadConfig(cfgPath)
	return Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
	}
}

func (t Todo) Init() tea.Cmd {
	return nil
}

func (t Todo) View() tea.View {
	return t.plannerModel.View()
}

func (t Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case planner.UpdateLabelMsg:
		sDate := t.plannerModel.CalendarModel.SelectedDate
		sLabel := t.plannerModel.Labels.Selected
		g, _ := database.FetchGrid(t.db, sDate, sLabel)
		t.plannerModel.Entries = g
		t.plannerModel.CalendarModel.UpdateGrid(t.plannerModel.HasEntry)
	case *planner.UpdateMonthMsg:
		sDate := t.plannerModel.CalendarModel.SelectedDate
		sLabel := t.plannerModel.Labels.Selected
		g, _ := database.FetchGrid(t.db, sDate, sLabel)
		t.plannerModel.Entries = g
		t.plannerModel.CalendarModel.UpdateGrid(t.plannerModel.HasEntry)
	default:
		t.plannerModel, cmd = t.plannerModel.Update(msg)
	}
	return t, cmd
}
