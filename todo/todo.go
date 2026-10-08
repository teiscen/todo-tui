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

	// err is the last database error. Not displayed anywhere yet.
	err error
}

func NewTodo(cfg config.Config, db *sql.DB) Todo {
	return Todo{
		db:           db,
		cfg:          cfg,
		plannerModel: cfg.LoadPlannerModel(),
	}
}

func (t Todo) Init() tea.Cmd {
	return t.plannerModel.Init()
}

func (t Todo) View() tea.View {
	return t.plannerModel.View()
}

func (t Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case planner.UpdateMonthMsg, planner.UpdateLabelMsg:
		t.reloadGrid()
		return t, nil
	case planner.SaveEntryMsg:
		t.saveEntry(msg)
		return t, nil
	case planner.RemoveEntryMsg:
		t.removeEntry(msg)
		return t, nil
	}

	var cmd tea.Cmd
	t.plannerModel, cmd = t.plannerModel.Update(msg)
	return t, cmd
}

// loadEntries fetches the visible month for the selected label and rebuilds
// the calendar grid. On error the grid is rebuilt with no marks.
func (t *Todo) loadEntries() {
	cal := &t.plannerModel.CalendarModel

	entries, err := database.FetchGrid(t.db, cal.SelectedDate, t.plannerModel.Labels.Selected)
	t.err = err

	t.plannerModel.Entries = entries
	cal.UpdateGrid(t.plannerModel.HasEntry)
}

// reloadGrid is loadEntries plus showing the selected day's entry in Notes.
func (t *Todo) reloadGrid() {
	t.loadEntries()
	t.plannerModel.SyncNotes()
}

// saveEntry keeps Notes as it is, so the cursor doesn't jump while typing.
func (t *Todo) saveEntry(m planner.SaveEntryMsg) {
	err := database.AddEntry(t.db, m.Date.FormatSQL(), string(m.Label), m.Msg)
	if t.err = err; err != nil {
		return
	}
	t.loadEntries()
}

func (t *Todo) removeEntry(m planner.RemoveEntryMsg) {
	err := database.RemoveEntry(t.db, m.Date.FormatSQL(), string(m.Label))
	if t.err = err; err != nil {
		return
	}
	t.reloadGrid()
}
