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

// func CalendarTestPrint() {
// 	c := TestCalendar()
// 	println(c.View())
// }
//
// func CalendarTestUpdate() {
// 	// fmt.Print("\034[H\033[2J") // Clear console ANSI escape sequence
// 	locn := time.Now().Location()
// 	sep27 := time.Date(2026, time.September, 27, 0, 0, 0, 0, locn)
// 	oct13 := time.Date(2026, time.October, 13, 0, 0, 0, 0, locn)
// 	oct17 := time.Date(2026, time.October, 17, 0, 0, 0, 0, locn)
//
// 	c := TestCalendar()
//
// 	for range 10 {
// 		// Reset back to start
// 		c.SelectedDate = backend.TimeToDate(oct13)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Left
// 		fmt.Print("\034[H\033[2J")
// 		// c.Update()
// 		c.Movement(Up)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Right
// 		fmt.Print("\034[H\033[2J")
// 		// c.Update()
// 		c.Movement(Down)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Down
// 		fmt.Print("\034[H\033[2J")
// 		// c.Update()
// 		c.Movement(Right)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Up
// 		fmt.Print("\034[H\033[2J")
// 		// c.Update()
// 		c.Movement(Left)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Toggle Marks
// 		fmt.Print("\034[H\033[2J")
// 		// c.Update()
// 		c.ToggleMark()
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
// 		c.ToggleMark()
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// First Index in grid [0][0]
// 		c.SelectedDate = backend.TimeToDate(sep27)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
//
// 		// Weekend Example
// 		c.SelectedDate = backend.TimeToDate(oct17)
// 		println(c.View())
// 		time.Sleep(1 * time.Second)
// 	}
// }
