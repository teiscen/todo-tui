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

func NewPlanner(c backend.Calendar) Planner {
	return Planner{
		NotesModel:    notes.NewNotes(""),
		CalendarModel: calendar.NewCalendar(),
		Style:         GetDefaultStyle(),

		Focus: FocusCalendar,
	}
}

// Make the model use the current things
func (p Planner) Init() tea.Cmd {
	return nil
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
	// prevDate := p.CalendarModel.SelectedDate
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+w":
			// toggle focus between Calendar and Notes Models
			p.ToggleFocus()
		case "ctrl+h":
			// Swap Calendar for previous label
			cmd = p.UpdateLabel(false)
		case "ctrl+j":
			// Swap Calendar to next month
			cmd = p.UpdateMonth(true)
		case "ctrl+k":
			// Swap Calendar to previous month
			cmd = p.UpdateMonth(false)
		case "ctrl+l":
			// Swap Calendar for next label
			cmd = p.UpdateLabel(true)
		case "ctrl+r":
			// cmd = p.WriteChange()
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	switch p.Focus {
	case FocusCalendar:
		p.CalendarModel, cmd = p.CalendarModel.Update(msg)
	case FocusNotes:
		p.NotesModel, cmd = p.NotesModel.Update(msg)
	}

	// TODO: refactor so its not ungly
	// Update Calendar if its a new month
	// Update Notes value if the day changed
	// Need to update calendar via result from x CalendarModel
	//
	// newStr := ""
	// if newText, ok := p.Calendar.GetEntry(p.CalendarModel.SelectedDate); ok {
	// 	newStr = newText.Msg
	// }
	// if prevDate.Month != p.CalendarModel.SelectedDate.Month {
	// 	p.CalendarModel.UpdateGrid(p.Calendar)
	// 	p.NotesModel.ChangeValue(newStr)
	// } else if prevDate != p.CalendarModel.SelectedDate {
	// 	p.NotesModel.ChangeValue(newStr)
	// }

	return p, cmd
}
