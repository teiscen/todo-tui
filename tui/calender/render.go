package calender

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

type (
	RenderCell struct {
		isValid bool
		cDate   CalenderDate
		str     string
	}
	RenderGrid []RenderCell
)

func (m model) CreateGrid() (g RenderGrid) {
	cal := m.calender
	cDate := m.cursor
	label := m.page
	labelColor := label.info().color

	firstOfMonth := cDate.firstDay()
	firstDayOffset := int(firstOfMonth.weekday())
	startDate := firstOfMonth.addDate(0, 0, -1*firstDayOffset)
	for i := 0; i < firstDayOffset; i++ {
		d := startDate.addDate(0, 0, i)
		str := m.style.invalidRenderer(ColorInvalid)
		if _, ok := cal.Months[d][label]; ok {
			str = m.style.invalidRenderer(labelColor)
		}
		g = append(g, RenderCell{false, d, str})
	}

	for i := 0; i < cDate.numDays(); i++ {
		d := firstOfMonth.addDate(0, 0, i)
		str := m.style.cellRenderer(ColorValid)
		if _, ok := cal.Months[d][label]; ok {
			str = m.style.cellRenderer(labelColor)
		}
		g = append(g, RenderCell{true, d, str})
	}

	lastOfMonth := cDate.lastDay()
	lastDayOffset := int(lastOfMonth.weekday())
	trailingDays := 6 - lastDayOffset
	for i := 0; i < trailingDays; i++ {
		d := lastOfMonth.addDate(0, 0, i+1)
		str := m.style.invalidRenderer(ColorInvalid)
		if _, ok := cal.Months[d][label]; ok {
			str = m.style.invalidRenderer(labelColor)
		}
		g = append(g, RenderCell{false, d, str})
	}

	return
}

func (m model) UpdateSelected(cal Calender, cDate CalenderDate, grid RenderGrid, l Label) (g RenderGrid) {
	for i, c := range grid {
		if cDate == c.cDate {
			var color ThemeColor
			if _, ok := cal.Months[cDate][l]; ok {
				color = l.info().color
			} else if c.isValid {
				color = ColorValid
			} else {
				color = ColorInvalid
			}

			grid[i] = RenderCell{
				c.isValid,
				c.cDate,
				m.style.selectedRenderer(color),
			}
		}
	}
	return grid
}

func (m model) Render() string {
	cells := make([]string, len(m.grid))
	for i, c := range m.grid {
		cells[i] = c.str
	}

	var weeks []string
	for i := 0; i < len(cells); i += 7 {
		end := i + 7
		if end > len(cells) {
			end = len(cells)
		}
		weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Top, cells[i:end]...))
	}

	header := m.style.weekdayRenderer()
	sep := "\n" + strings.Repeat("\n", m.style.getRowSpacing())
	weeksBlock := strings.Join(weeks, sep)

	return header + "\n" + weeksBlock + "\n"
}
