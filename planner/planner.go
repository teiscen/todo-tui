package planner

import (
	backend "todo-tui/backend_old"
	"todo-tui/calendar"
	"todo-tui/notes"

	tea "charm.land/bubbletea/v2"
)

type Focus int

const (
	FocusCalendar Focus = iota
	FocusNotes
)

type Planner struct {
	NotesModel    notes.Notes
	CalendarModel calendar.Calendar
	Style         PlannerStyle

	Calendar backend.Calendar
	Labels   backend.Labels

	Focus Focus
}

func NewPlanner(c backend.Calendar) Planner {
	return Planner{}
}

func (p Planner) Init() tea.Cmd {
	return nil
}

func (p Planner) View() tea.View {
	return tea.NewView("NEW VIEW FOR PLANNER")
}

func (p Planner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+w":
			// toggle focus between Calendar and Notes Models
		case "ctrl+h":
			// Swap Calendar for previous label
		case "ctrl+j":
			// Swap Calendar to next month
		case "ctrl+k":
			// Swap Calendar to previous month
		case "ctrl+l":
			// Swap Calendar for next label
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	return p, cmd
}
