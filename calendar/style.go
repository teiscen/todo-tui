package calendar

import (
	lipgloss "charm.land/lipgloss/v2"
)

const (
	largeCircle    = "⬤"
	largeHexRing   = "⬡"
	smallCircle    = "●"
	smallRing      = "◯"
	selectedCircle = "◉"
)

type CalendarStyle struct {
	columnWidth int
	rowSpacing  int

	weekdays    []string
	weekdayBold bool

	cellChar     string
	invalidChar  string
	selectedChar string
}

var singleStyle = CalendarStyle{
	columnWidth: 3,
	rowSpacing:  0,

	weekdays:    []string{"S", "M", "T", "W", "T", "F", "S"},
	weekdayBold: false,

	cellChar:     smallCircle,
	invalidChar:  smallRing,
	selectedChar: selectedCircle,
}

var doubleStyle = CalendarStyle{
	columnWidth: 4,
	rowSpacing:  1,

	weekdays:    []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"},
	weekdayBold: true,

	cellChar:     largeCircle,
	invalidChar:  largeHexRing,
	selectedChar: "[" + largeCircle + " ]",
}

func (s CalendarStyle) weekdayRenderer(valid Color, weekend Color) string {
	var out string

	for i, day := range s.weekdays {
		color := valid

		if i == 0 || i == 6 {
			color = weekend
		}

		out += color.toStyle().
			Width(s.columnWidth).
			Align(lipgloss.Center).
			Bold(s.weekdayBold).
			Render(day)
	}

	return "\n" + out + "\n"
}

func (s CalendarStyle) invalidRenderer(c Color) string {
	return c.toStyle().
		Width(s.columnWidth).
		Align(lipgloss.Center).
		Render(s.invalidChar)
}

func (s CalendarStyle) cellRenderer(c Color) string {
	return c.toStyle().
		Width(s.columnWidth).
		Align(lipgloss.Center).
		Render(s.cellChar)
}

func (s CalendarStyle) selectedRenderer(c Color) string {
	return c.toStyle().
		Width(s.columnWidth).
		Align(lipgloss.Center).
		Bold(true).
		Render(s.selectedChar)
}
