package calendar

import backend "todo-tui/backend_old"

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
