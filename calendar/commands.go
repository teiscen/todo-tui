package calendar

// Mark/Unmark
// Movement
type direction int

const (
	Left direction = iota
	Down
	Up
	Right
)

func (c *Calendar) Movement(dir direction) {
	switch dir {
	case Left:
		c.SelectedDate = c.SelectedDate.AddDate(0, 0, -1)
	case Down:
		c.SelectedDate = c.SelectedDate.AddDate(0, 0, 7)
	case Up:
		c.SelectedDate = c.SelectedDate.AddDate(0, 0, -7)
	case Right:
		c.SelectedDate = c.SelectedDate.AddDate(0, 0, 1)
	}
}

// Update changes maybe
func (c *Calendar) ToggleMark() {
	for row := range c.Grid {
		for col := range c.Grid[row] {
			cell := &c.Grid[row][col]
			if cell.Date == c.SelectedDate {
				cell.IsMarked = !(*cell).IsMarked
			}
		}
	}
}
