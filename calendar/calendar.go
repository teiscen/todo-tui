package calendar

import (
	"strings"

	"todo-tui/backend"

	lipgloss "charm.land/lipgloss/v2"
)

// TODO: need to restructure some of the render logic
// to make better use of custom themes/color
type Calendar struct {
	cri      backend.CalendarRenderInfo
	calstyle CalendarStyle
}

func (c Calendar) Render() string {
	header := c.calstyle.weekdayRenderer()

	var cells []string
	for row := range c.cri.Grid {
		for col := range c.cri.Grid[row] {
			cell := c.cri.Grid[row][col]
			color := cell.Color

			var str string
			if cell.IsValid {
				str = c.calstyle.validRenderer(color)
			} else {
				str = c.calstyle.invalidRenderer(color)
			}
			if cell.Selected {
				str = c.calstyle.selectedRenderer(color)
			}
			cells = append(cells, str)
		}
	}

	var weeks []string
	for i := 0; i < len(cells); i += 7 {
		end := i + 7
		if end > len(cells) {
			end = len(cells)
		}
		weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Top, cells[i:end]...))
	}

	sep := "\n" + strings.Repeat("\n", c.calstyle.rowSpacing)
	weeksBlock := strings.Join(weeks, sep)

	str := header + "\n" + weeksBlock
	return lipgloss.NewStyle().Padding(1, 1).Render(str)
}
