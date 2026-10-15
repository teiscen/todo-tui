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

type Entries map[string]string

type Planner struct {
	NotesModel    notes.Notes
	CalendarModel calendar.Calendar
	Style         PlannerStyle

	Focus Focus

	Entries Entries
	Labels  backend.Labels
}

func NewPlanner(date backend.Date) Planner {
	return Planner{
		NotesModel:    notes.NewNotes(""),
		CalendarModel: calendar.NewCalendar(date),
		Style:         GetDefaultStyle(),

		Focus: FocusCalendar,
	}
}

// Init asks the root to load the starting month, the same way a month change does.
func (p Planner) Init() tea.Cmd {
	return nil
}

func (p Planner) View() tea.View {
	type posn int
	const (
		top posn = iota
		mid
		bot
	)
	width := 28
	genBreak := func(text string, style lipgloss.Style, posn posn) string {
		text = style.Render(text)

		var paddingL, paddingR, left, right string
		pad := strings.Repeat("─", max(0, (width-lipgloss.Width(text))/2))
		if lipgloss.Width(text)%2 == 1 {
			paddingR = "─"
		}
		paddingL += pad
		paddingR += pad

		switch posn {
		case top:
			left = "╭" + paddingL
			right = paddingR + "╮"
		case mid:
			left = "├" + paddingL
			right = paddingR + "┤"
		case bot:
			left = "╰" + paddingL
			right = paddingR + "╯"
		}

		borderStyle := p.Style.BorderColor.ToStyle()
		return borderStyle.Render(left) + text + borderStyle.Render(right)
	}

	headerFull := genBreak(
		p.CalendarModel.SelectedDate.Format(),
		p.CalendarModel.GetSelectedColor().ToStyle().Padding(0, 1).Bold(true),
		top,
	)

	calendarBorder := p.Style.BorderColor.ToStyle().
		Border(lipgloss.RoundedBorder(), false, true).
		Render(p.CalendarModel.View().Content)

	footerFull := genBreak(
		p.Labels.GetName(),
		p.Labels.GetColor().ToStyle().Padding(0, 1).Bold(true),
		mid,
	)

	notesBorder := p.Style.BorderColor.ToStyle().Padding(1, 1).
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Render(p.NotesModel.View().Content)

	complete := lipgloss.JoinVertical(
		lipgloss.Left,
		headerFull,
		calendarBorder,
		footerFull,
		notesBorder,
	)

	return tea.NewView(complete)
}

// Handle on inputs:
// ctrl h/j/k/l and r
// Then delegate to the submodels (planner/notes)
// Read the output and determine if actions are needed
func (p Planner) Update(msg tea.Msg) (Planner, tea.Cmd) {
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
		case "ctrl+c":
			return p, tea.Quit
		}
	}

	var cmd tea.Cmd // initializes to nil
	switch p.Focus {
	case FocusCalendar:
		p.CalendarModel, cmd = p.CalendarModel.Update(msg)
	case FocusNotes:
		p.NotesModel, cmd = p.NotesModel.Update(msg)
	}

	// Capture the changes returned by notes or calendar model
	switch m := msg.(type) {
	case calendar.ChangeMonthMsg:
		return p, func() tea.Msg {
			p.SyncNotes()
			return ChangeCalendarMsg{
				*p.SelectedLabel(),
				*p.SelectedDate(),
			}
		}
	case calendar.ChangeDateMsg:
		p.SyncNotes()
	case calendar.ToggleMarkMsg:
		if m.Status {
			return p, func() tea.Msg {
				return SaveEntryMsg{
					Date:  *p.SelectedDate(),
					Label: *p.SelectedLabel(),
					Msg:   p.NotesModel.TextArea.Value(),
				}
			}
		} else {
			return p, func() tea.Msg {
				return RemoveEntryMsg{
					Date:  *p.SelectedDate(),
					Label: *p.SelectedLabel(),
				}
			}
		}
	case notes.WriteTextMsg:
		return p, func() tea.Msg {
			return SaveEntryMsg{
				Date:  *p.SelectedDate(),
				Label: *p.SelectedLabel(),
				Msg:   p.NotesModel.TextArea.Value(),
			}
		}
	}
	return p, cmd
}
