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

// Messages for the root: "the month/label changed, reload the grid".
type ChangeCalendarMsg struct {
	LabelID backend.LabelID
	NewDate backend.Date
}

func (p *Planner) UpdateMonth(next bool) tea.Cmd {
	step := -1
	if next {
		step = 1
	}
	p.SelectedDate().ChangeDate(0, step, 0)
	return func() tea.Msg {
		return ChangeCalendarMsg{
			*p.SelectedLabel(),
			*p.SelectedDate(),
		}
	}
}

func (p *Planner) UpdateLabel(next bool) tea.Cmd {
	p.Labels.Step(next)
	p.ApplyLabelColor()
	return func() tea.Msg {
		return ChangeCalendarMsg{
			*p.SelectedLabel(),
			*p.SelectedDate(),
		}
	}
}

// Messages for the root: "persist this change".
type SaveEntryMsg struct {
	Date  backend.Date
	Label backend.LabelID
	Msg   string
}

type RemoveEntryMsg struct {
	Date  backend.Date
	Label backend.LabelID
}

// SaveEntry saves the notes text as the selected day's entry (ctrl+r).
// Saving also marks the day.
// func (p *Planner) SaveEntry() tea.Cmd {
// 	if p.Labels.Selected == "" {
// 		return nil
// 	}
// 	msg := SaveEntryMsg{
// 		Date:  p.CalendarModel.SelectedDate,
// 		Label: p.Labels.Selected,
// 		Msg:   p.NotesModel.TextArea.Value(),
// 	}
// 	return func() tea.Msg { return msg }
// }

// ToggleEntry marks or unmarks the selected day (x). Marking creates an entry
// from the current notes text, unmarking removes the entry.
// func (p *Planner) ToggleEntry() tea.Cmd {
// 	date := p.CalendarModel.SelectedDate
// 	if !p.HasEntry(date) {
// 		return p.SaveEntry()
// 	}
// 	msg := RemoveEntryMsg{Date: date, Label: p.Labels.Selected}
// 	return func() tea.Msg { return msg }
// }
