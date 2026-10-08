package planner

import (
	backend "todo-tui/backend"

	tea "charm.land/bubbletea/v2"
)

// Command:
// Switch between the labels (updates calendar)
// Swtich between the months (updates calendar)

func (p *Planner) ToggleFocus() {
	if p.Focus == FocusCalendar {
		p.Focus = FocusNotes
	} else {
		p.Focus = FocusCalendar
	}
	p.NotesModel.ToggleFocus()
}

func (p *Planner) HasEntry(date backend.Date) bool {
	_, ok := p.Entries[date]
	return ok
}

type UpdateMonthMsg struct{ newDate backend.Date }

func (p *Planner) UpdateMonth(next bool) tea.Cmd {
	curr := p.CalendarModel.SelectedDate
	if next {
		p.CalendarModel.SelectedDate = curr.AddDate(0, 1, 0)
	} else {
		p.CalendarModel.SelectedDate = curr.AddDate(0, -1, 0)
	}
	return func() tea.Msg { return UpdateMonthMsg{p.CalendarModel.SelectedDate} }
}

type UpdateLabelMsg struct{ labelID backend.LabelID }

func (p *Planner) UpdateLabel(next bool) tea.Cmd {
	if next {
		p.Labels.Step(true)
	} else {
		p.Labels.Step(false)
	}
	return func() tea.Msg { return UpdateLabelMsg{p.Labels.Selected} }
}

func (p *Planner) WriteChange() {
	// d := p.CalendarModel.SelectedDate
	// e := backend.Entry{
	// 	Status: backend.Full,
	// 	Msg:    p.NotesModel.TextArea.Value(),
	// }
}
