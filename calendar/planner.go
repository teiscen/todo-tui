package calendar

import (
	"todo-tui/backend"
)

// Might need to track window size to adjust render
type Planner struct {
	Cal  Calendar
	Note Note
}

func NewPlanner(s backend.State) Planner {
	return Planner{
		NewCalendar(s),
		NewNote(s),
	}
}

func (p Planner) Render() string {
	calRender := p.Cal.Render()
	noteRender := p.Note.Render()
	planRender := p.RenderStr(calRender, noteRender)

	return planRender
}
