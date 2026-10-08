package calendar

import backend "todo-tui/backend"

// TODO:
// Change how this works, wether by indexing directly
// or getting rid of the need for it
func (c Calendar) GetSelectedColor() backend.HexCode {
	for row := range c.Grid {
		for col := range c.Grid[row] {
			if c.SelectedDate == c.Grid[row][col].Date {
				cell := c.Grid[row][col]

				var color backend.HexCode
				if cell.IsMarked {
					color = c.Style.ColorAccent
				} else if col == 0 || col == 6 {
					color = c.Style.ColorWeekend
				} else {
					color = c.Style.ColorWeekday
				}

				if !cell.IsValid {
					color, _ = color.MuteColorTarget(backend.HexCode("#303446"), 0.8)
				}
				return color
			}
		}
	}

	return backend.HexCode("#FFFFFF")
}

func (c *Calendar) SetAccent(color backend.HexCode) {
	c.Style.ColorAccent = color
}

func (c *Calendar) UpdateGrid(hasEntry func(backend.Date) bool) {
	checkValid := func(target backend.Date) bool {
		return target.Month == c.SelectedDate.Month
	}

	checkMarked := func(target backend.Date) bool {
		return hasEntry(target)
	}

	date := c.SelectedDate.GridStart()

	for row := range c.Grid {
		for col := range c.Grid[row] {
			cell := &c.Grid[row][col]

			*cell = Grid{
				Date:     date,
				IsValid:  checkValid(date),
				IsMarked: checkMarked(date),
			}

			date = date.AddDate(0, 0, 1)
		}
	}
}
