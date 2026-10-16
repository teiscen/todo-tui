package calendar

import backend "todo-tui/backend"

func (c *Calendar) GetInfo(date backend.Date) *Info {
	// TODO:
	// Find a more efficient method
	for row := range c.State {
		for col := range c.State[row] {
			if date == c.State[row][col].Date {
				return &c.State[row][col]
			}
		}
	}
	return nil
}

func (c Calendar) GetSelectedColor() backend.HexCode {
	var color backend.HexCode

	infoAddr := c.GetInfo(c.SelectedDate)
	if infoAddr == nil {
		return backend.HexCode("#FFFFFF")
	}
	if infoAddr.IsMarked {
		color = c.Style.Theme.Accent
	} else if infoAddr.Date.IsWeekend() {
		color = c.Style.Theme.Weekend
	} else {
		color = c.Style.Theme.Weekday
	}

	if infoAddr.Date.Month() != c.SelectedDate.Month() {
		color, _ = color.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return color
}

func (c *Calendar) SetAccent(color backend.HexCode) {
	c.Style.Theme.Accent = color
}
