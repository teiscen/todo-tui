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
		color = c.Style.ColorAccent
	} else if infoAddr.Date.Day() == 0 ||
		infoAddr.Date.Day() == 6 {
		color = c.Style.ColorWeekend
	} else {
		color = c.Style.ColorWeekday
	}

	if infoAddr.Date.Month() != c.SelectedDate.Month() {
		color, _ = color.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return color
}

func (c *Calendar) SetAccent(color backend.HexCode) {
	c.Style.ColorAccent = color
}

// func (c *Calendar) UpdateGrid(hasEntry func(backend.Date) bool) {
// 	checkValid := func(target backend.Date) bool {
// 		return target.Month == c.SelectedDate.Month
// 	}
//
// 	checkMarked := func(target backend.Date) bool {
// 		return hasEntry(target)
// 	}
//
// 	date := c.SelectedDate.GridStart()
//
// 	for row := range c.State {
// 		for col := range c.State[row] {
// 			cell := &c.State[row][col]
//
// 			*cell = Info{
// 				Date:     date,
// 				IsValid:  checkValid(date),
// 				IsMarked: checkMarked(date),
// 			}
//
// 			date = date.AddDate(0, 0, 1)
// 		}
// 	}
// }
