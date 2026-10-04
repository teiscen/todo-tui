package planner

import (
	backend "todo-tui/backend_old"
	"todo-tui/calendar"
)

func TestPlanner() Planner {
	cModel := calendar.TestCalendar()
	c := backend.MockCalendar()
	p := NewPlanner(c)
	p.CalendarModel = cModel
	return p
}
