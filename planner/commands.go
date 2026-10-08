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

// Messages for the root: "the month/label changed, reload the grid".
type UpdateMonthMsg struct{ newDate backend.Date }

func (p *Planner) UpdateMonth(next bool) tea.Cmd {
	step := -1
	if next {
		step = 1
	}
	p.CalendarModel.SelectedDate = p.CalendarModel.SelectedDate.AddMonths(step)

	d := p.CalendarModel.SelectedDate // capture a value, not p
	return func() tea.Msg { return UpdateMonthMsg{d} }
}

type UpdateLabelMsg struct{ labelID backend.LabelID }

func (p *Planner) UpdateLabel(next bool) tea.Cmd {
	p.Labels.Step(next)
	p.ApplyLabelColor()

	id := p.Labels.Selected // capture a value, not p
	return func() tea.Msg { return UpdateLabelMsg{id} }
}

// ApplyLabelColor makes the calendar and notes accents follow the selected label.
func (p *Planner) ApplyLabelColor() {
	c := p.Labels.GetColor()
	p.CalendarModel.SetAccent(c)
	p.NotesModel.SetAccent(c)
}

// SyncNotes shows the selected day's entry in Notes. A day with no entry
// clears it, which shows the placeholder.
func (p *Planner) SyncNotes() {
	p.NotesModel.ChangeValue(p.Entries[p.CalendarModel.SelectedDate].Msg)
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
func (p *Planner) SaveEntry() tea.Cmd {
	if p.Labels.Selected == "" {
		return nil
	}
	msg := SaveEntryMsg{
		Date:  p.CalendarModel.SelectedDate,
		Label: p.Labels.Selected,
		Msg:   p.NotesModel.TextArea.Value(),
	}
	return func() tea.Msg { return msg }
}

// ToggleEntry marks or unmarks the selected day (x). Marking creates an entry
// from the current notes text, unmarking removes the entry.
func (p *Planner) ToggleEntry() tea.Cmd {
	date := p.CalendarModel.SelectedDate
	if !p.HasEntry(date) {
		return p.SaveEntry()
	}
	msg := RemoveEntryMsg{Date: date, Label: p.Labels.Selected}
	return func() tea.Msg { return msg }
}
