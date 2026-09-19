package calendar

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

// can decide to shift some computation over to
// RenderInfo methods (such as no InBounds)
// Calculating some thins as we go
type Cell struct {
	cDate    CalendarDate
	inBound  bool
	labelMap map[Label]Entry
}

type RenderInfo struct {
	cMonth CalendarDate
	grid   []Cell
}

func (m Model) cellColor(c Cell) Color {
	label := m.labels[m.selLabel]

	if _, ok := c.labelMap[label]; ok {
		info := m.config.Labels[label]
		return m.config.Colors[info.ColorName]
	}

	if c.inBound {
		return m.config.Colors["valid"]
	}

	return m.config.Colors["invalid"]
}

func (m Model) cursorColor() Color {
	label := m.labels[m.selLabel]

	if entries := m.calendar.Months[m.selDate]; entries != nil {
		if _, ok := entries[label]; ok {
			info := m.config.Labels[label]
			return m.config.Colors[info.ColorName]
		}
	}

	if m.selDate.Month == m.rInfo.cMonth.Month &&
		m.selDate.Year == m.rInfo.cMonth.Year {
		return m.config.Colors["valid"]
	}

	return m.config.Colors["invalid"]
}

func CreateGrid(m Model) RenderInfo {
	var out []Cell

	// Find the first day
	d := m.selDate.firstDay()
	firstOffset := int(d.weekday())
	startDate := d.addDate(0, 0, -1*firstOffset)

	d = m.selDate.lastDay()
	trailOffset := 6 - int(d.weekday())

	totalCount := firstOffset + m.selDate.numDays() + trailOffset

	for i := 0; i < totalCount; i++ {
		cDate := startDate.addDate(0, 0, i)
		inBounds := true
		if i < firstOffset || i >= firstOffset+m.selDate.numDays() {
			inBounds = false
		}
		out = append(out, Cell{
			cDate,
			inBounds,
			m.calendar.Months[cDate],
		})
	}

	return RenderInfo{m.selDate.firstDay(), out}
}

func (m Model) Render() string {
	cells := make([]string, len(m.rInfo.grid))

	for i, c := range m.rInfo.grid {
		color := m.cellColor(c)

		if c.cDate == m.selDate {
			cells[i] = m.style.selectedRenderer(color)
		} else if !c.inBound {
			cells[i] = m.style.invalidRenderer(color)
		} else {
			cells[i] = m.style.cellRenderer(color)
		}
	}

	var weeks []string
	for i := 0; i < len(cells); i += 7 {
		end := i + 7
		if end > len(cells) {
			end = len(cells)
		}

		weeks = append(
			weeks,
			lipgloss.JoinHorizontal(lipgloss.Top, cells[i:end]...),
		)
	}

	header := m.style.weekdayRenderer(
		m.config.Colors["valid"],
		m.config.Colors["red"],
	)

	sep := "\n" + strings.Repeat("\n", m.style.rowSpacing)
	weeksBlock := strings.Join(weeks, sep)
	return header + "\n" + weeksBlock + "\n"
}
