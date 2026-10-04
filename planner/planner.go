package planner

import (
	"strings"

	backend "todo-tui/backend_old"
	"todo-tui/calendar"
	"todo-tui/notes"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	Focus    Focus
}

func NewPlanner(c backend.Calendar) Planner {
	return Planner{
		NotesModel:    notes.NewNotes(""),
		CalendarModel: calendar.NewCalendar(),
		Style:         GetDefaultStyle(),

		Calendar: c,
		Focus:    FocusCalendar,
	}
}

func (p Planner) Init() tea.Cmd {
	return nil
}

func (p Planner) View() tea.View {
	headerStyle := p.CalendarModel.GetSelectedColor().ToStyle().Padding(0, 1)
	footerStyle := p.Calendar.GetCurrentLabel().Color.ToStyle().Padding(0, 1)

	headerText := p.CalendarModel.SelectedDate.Format()
	footerText := p.Calendar.GetCurrentLabel().Name

	headerFormatted := headerStyle.Render(headerText)
	footerFormatted := footerStyle.Render(footerText)

	width := 30
	padding := func(str string) string {
		return strings.Repeat("─", max(0, (width-lipgloss.Width(str))/2))
	}

	headerPadding := padding(headerFormatted)
	footerPadding := padding(footerFormatted)

	headerFull := "╭" + headerPadding + headerFormatted + headerPadding + "╮"
	footerFull := "├" + footerPadding + footerFormatted + footerPadding + "┤"

	notesView := p.NotesModel.View().Content
	calendarView := p.CalendarModel.View().Content

	calendarBorder := p.Style.BorderColor.ToStyle().
		Border(lipgloss.NormalBorder(), false, true).Render(calendarView)

	notesBorder := p.Style.BorderColor.ToStyle().Padding(1, 1).
		Border(lipgloss.NormalBorder(), false, true, true, true).Render(notesView)

	complete := lipgloss.JoinVertical(
		lipgloss.Left,
		headerFull,
		calendarBorder,
		footerFull,
		notesBorder,
	)

	return tea.NewView(complete)
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
