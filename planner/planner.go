package planner

import (
	backend "todo-tui/backend"
	"todo-tui/calendar"
	"todo-tui/notes"
	"todo-tui/ui"

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

// New builds a planner pane composing the calendar and notes sub-panes.
func New(theme backend.Theme, labels backend.Labels, date backend.Date) Planner {
	p := Planner{
		NotesModel:    notes.New(theme, ""),
		CalendarModel: calendar.New(theme, date),
		Style:         NewStyle(theme),

		Labels: labels,
		Focus:  FocusCalendar,
	}
	p.ApplyLabelColor()
	return p
}

// Init asks the root to load the starting month, the same way a month change does.
func (p Planner) Init() tea.Cmd {
	return nil
}

func (p Planner) View() string {
	width := 28
	border := ui.Fg(p.Style.Theme.Border).BorderForeground(lipgloss.Color(string(p.Style.Theme.Border)))

	headerFull := ui.GenBorder(
		ui.Fg(p.CalendarModel.GetSelectedColor()).Padding(0, 1).Bold(true).Render(p.CalendarModel.SelectedDate.Format()),
		border,
		ui.Top,
		width,
	)

	calendarBorder := border.
		Border(lipgloss.RoundedBorder(), false, true).
		Render(p.CalendarModel.View())

	footerFull := ui.GenBorder(
		ui.Fg(p.Labels.GetColor()).Padding(0, 1).Bold(true).Render(p.Labels.GetName()),
		border,
		ui.Mid,
		width,
	)

	notesBorder := border.Padding(1, 1).
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Render(p.NotesModel.View())

	complete := lipgloss.JoinVertical(
		lipgloss.Left,
		headerFull,
		calendarBorder,
		footerFull,
		notesBorder,
	)

	return complete
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
