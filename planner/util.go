package planner

import (
	"todo-tui/backend"
	"todo-tui/calendar"
)

// GridStart is the first cell of a 6x7 (Sunday-first) month grid: the Sunday
// on or before the 1st of d's month.
// func GridStart(d backend.Date) backend.Date {
// 	first := d.FirstDay()
// 	return first.AddDate(0, 0, -int(first.Weekday()))
// }

func (p *Planner) SelectedDate() *backend.Date {
	return &p.CalendarModel.SelectedDate
}

func (p *Planner) CalendarState() *[6][7]calendar.Info {
	return &p.CalendarModel.State
}

func (p *Planner) SelectedLabel() *backend.LabelID {
	return &p.Labels.Selected
}

// ApplyLabelColor makes the calendar and notes accents follow the selected label.
func (p *Planner) ApplyLabelColor() {
	c := p.Labels.GetColor()
	p.CalendarModel.SetAccent(c)
	p.NotesModel.SetAccent(c)
}

func (p *Planner) HasEntry(date backend.Date) bool {
	_, ok := p.Entries[date.Key()]
	return ok
}

// Take in the result for entries
func (p *Planner) ReloadState(sDate backend.Date, entries map[string]string) {
	// Update Calendar state
	// Update Notes
	p.Entries = entries
	p.UpdateCalendarState(sDate)
	p.SyncNotes()
}

func (p *Planner) UpdateCalendarState(sDate backend.Date) {
	state := p.CalendarState()

	for row := range state {
		for col := range state[row] {
			_, ok := p.Entries[sDate.Key()]
			state[row][col] = calendar.Info{Date: sDate, IsMarked: ok}
			sDate.ChangeDate(0, 0, 1)
		}
	}
}

// SyncNotes shows the selected day's entry in Notes. A day with no entry
// clears it, which shows the placeholder.
func (p *Planner) SyncNotes() {
	sel := *p.SelectedDate()
	p.NotesModel.ChangeValue(p.Entries[sel.Key()])
}
