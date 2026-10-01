package calendar

import (
	"time"

	backend "todo-tui/backend_old"
)

func TestCalendar() Calendar {
	locn := time.Now().Location()
	sep27 := time.Date(2026, time.September, 27, 0, 0, 0, 0, locn)
	oct13 := time.Date(2026, time.October, 13, 0, 0, 0, 0, locn)

	c := NewCalendar()
	c.SelectedDate = backend.TimeToDate(oct13)

	date := backend.TimeToDate(sep27)
	// initialize a Grid
	for row := range c.Grid {
		for col := range c.Grid[row] {
			c.Grid[row][col] = Grid{
				Date:     date,
				IsValid:  date.Month == 10,
				IsMarked: (row+col)%3 == 0,
			}
			date = date.AddDate(0, 0, 1)
		}
	}

	c.Grid[0][0].IsMarked = true
	c.Grid[0][1].IsMarked = true
	return c
}

func CalendarTestPrint() {
	c := TestCalendar()
	println(c.Render())
}
