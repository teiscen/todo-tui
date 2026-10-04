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

func NewPlanner(c backend.Calendar) Planner {Fr  Sa  │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
│  ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   ⬤   │
│                              │
├──────── Application ────────┤
│                              │
│ hlkjhflskadjhjjjlkjl;skjlkjl │
│ kj                           │
│                              │
│                              │
│ kljlkj                       │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
│ ~                            │
	return Planner{
		NotesModel:    notes.NewNotes(""),
		CalendarModel: calendar.NewCalendar(),
		Style:         GetDefaultStyle(),

		Calendar: c,
		Focus:    FocusCalendar,
	}
}

// Make the model use the current things
func (p Planner) Init() tea.Cmd {
	return nil
}

func (p Planner) View() tea.View {
	headerStyle := p.CalendarModel.GetSelectedColor().ToStyle().Padding(0, 1).Bold(true)
	footerStyle := p.Calendar.GetCurrentLabel().Color.ToStyle().Padding(0, 1).Bold(true)

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
		Border(lipgloss.RoundedBorder(), false, true).Render(calendarView)

	notesBorder := p.Style.BorderColor.ToStyle().Padding(1, 1).
		Border(lipgloss.RoundedBorder(), false, true, true, true).Render(notesView)

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
			p.ToggleFocus()
		case "ctrl+h":
			// Swap Calendar for previous label
			p.UpdateLabel(false)
		case "ctrl+j":
			// Swap Calendar to next month
			p.UpdateMonth(true)
		case "ctrl+k":
			// Swap Calendar to previous month
			p.UpdateMonth(false)
		case "ctrl+l":
			// Swap Calendar for next label
			p.UpdateLabel(true)
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	switch p.Focus {
	case FocusCalendar:
		prevDate := p.CalendarModel.SelectedDate
		p.CalendarModel, cmd = p.CalendarModel.Update(msg)
		if prevDate.Month != p.CalendarModel.SelectedDate.Month {
			p.CalendarModel.UpdateGrid(p.Calendar)
		}
	case FocusNotes:
		p.NotesModel, cmd = p.NotesModel.Update(msg)
	}

	return p, cmd
}
