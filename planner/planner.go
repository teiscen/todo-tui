package planner

import (
	"strings"

	backend "todo-tui/backend"
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

type Entries map[backend.Date]backend.Entry

type Planner struct {
	NotesModel    notes.Notes
	CalendarModel calendar.Calendar
	Style         PlannerStyle

	Focus Focus

	Entries Entries
	Labels  backend.Labels
}

func NewPlanner() Planner {
	return Planner{
		NotesModel:    notes.NewNotes(""),
		CalendarModel: calendar.NewCalendar(),
		Style:         GetDefaultStyle(),

		Focus: FocusCalendar,
	}
}

// Init asks the root to load the starting month, the same way a month change does.
func (p Planner) Init() tea.Cmd {
	d := p.CalendarModel.SelectedDate
	return func() tea.Msg { return UpdateMonthMsg{d} }
}

func (p Planner) View() tea.View {
	headerStyle := p.CalendarModel.GetSelectedColor().ToStyle().Padding(0, 1).Bold(true)
	footerStyle := p.Labels.GetColor().ToStyle().Padding(0, 1).Bold(true)

	headerText := p.CalendarModel.SelectedDate.Format()
	footerText := p.Labels.GetName()

	headerFormatted := headerStyle.Render(headerText)
	footerFormatted := footerStyle.Render(footerText)

	width := 30
	padding := func(str string) (left, right string) {
		pad := strings.Repeat("─", max(0, (width-lipgloss.Width(str))/2))
		if lipgloss.Width(str)%2 == 1 {
			right = "─"
		}
		left += pad
		right += pad
		return
	}

	headerPaddingL, headerPaddingR := padding(headerFormatted)
	footerPaddingL, footerPaddingR := padding(footerFormatted)

	headerFull := "╭" + headerPaddingL + headerFormatted + headerPaddingR + "╮"
	footerFull := "├" + footerPaddingL + footerFormatted + footerPaddingR + "┤"

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

func (p Planner) Update(msg tea.Msg) (Planner, tea.Cmd) {
	// Planner-level keys are handled here and not forwarded, otherwise the
	// focused child also sees them (e.g. ctrl+h and ctrl+w delete text in Notes).
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+w": // toggle focus between Calendar and Notes
			p.ToggleFocus()
			return p, nil
		case "ctrl+h": // previous label
			return p, p.UpdateLabel(false)
		case "ctrl+l": // next label
			return p, p.UpdateLabel(true)
		case "ctrl+j": // next month
			return p, p.UpdateMonth(true)
		case "ctrl+k": // previous month
			return p, p.UpdateMonth(false)
		case "ctrl+r": // save the entry
			return p, p.SaveEntry()
		case "x": // mark/unmark the day, only when the calendar has focus
			if p.Focus == FocusCalendar {
				return p, p.ToggleEntry()
			}
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	prev := p.CalendarModel.SelectedDate

	var cmd tea.Cmd
	switch p.Focus {
	case FocusCalendar:
		p.CalendarModel, cmd = p.CalendarModel.Update(msg)
	case FocusNotes:
		p.NotesModel, cmd = p.NotesModel.Update(msg)
	}

	// The calendar moved on its own (h/j/k/l). If it left the month, ask the
	// root to reload the grid, which also refreshes Notes. Otherwise only the
	// selected day changed, so show that day's entry.
	if cur := p.CalendarModel.SelectedDate; cur != prev {
		if cur.Year != prev.Year || cur.Month != prev.Month {
			cmd = tea.Batch(cmd, func() tea.Msg { return UpdateMonthMsg{cur} })
		} else {
			p.SyncNotes()
		}
	}

	return p, cmd
}
