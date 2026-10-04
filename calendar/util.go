package calendar

import backend "todo-tui/backend"

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

func (c *Calendar) UpdateGrid(cal backend.Calendar) {
	checkValid := func(target backend.Date) bool {
		return target.Month == c.SelectedDate.Month
	}

	checkMarked := func(target backend.Date) bool {
		_, ok := cal.GetEntry(target)
		return ok
	}

	offset := int(c.SelectedDate.FirstDay().FirstDay().Weekday())
	date := c.SelectedDate.FirstDay().FirstDay().AddDate(0, 0, -offset)

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
