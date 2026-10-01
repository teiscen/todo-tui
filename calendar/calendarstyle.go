package calendar

import (
	backend "todo-tui/backend_old"

	"charm.land/lipgloss/v2"
)

type CalendarStyle struct {
	ColumnWidth int
	RowSpacing  int

	DayStr []string

	CharIcon     string
	CharSelected string

	ColorWeekday backend.HexCode
	ColorWeekend backend.HexCode
	ColorAccent  backend.HexCode
}

func GetDefaultStyle() CalendarStyle {
	return CalendarStyle{
		ColumnWidth: 4,
		RowSpacing:  1,

		DayStr: []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"},

		CharIcon:     "⬤",
		CharSelected: "|⬤ |",

		ColorWeekday: backend.HexCode("#C6D0F5"),
		ColorWeekend: backend.HexCode("#939AB7"),
		ColorAccent:  backend.HexCode("#CA9EE6"),
	}
}

func (c CalendarStyle) WeekdayStyle(isValid bool) lipgloss.Style {
	color := c.ColorWeekday
	if !isValid {
		color, _ = c.ColorWeekday.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return color.ToStyle().Width(c.ColumnWidth).Align(lipgloss.Center)
}

func (c CalendarStyle) WeekendStyle(isValid bool) lipgloss.Style {
	color := c.ColorWeekend
	if !isValid {
		color, _ = c.ColorWeekend.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return color.ToStyle().Width(c.ColumnWidth).Align(lipgloss.Center)
}

func (c CalendarStyle) AccentStyle(isValid bool) lipgloss.Style {
	color := c.ColorAccent
	if !isValid {
		color, _ = c.ColorAccent.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return color.ToStyle().Width(c.ColumnWidth).Align(lipgloss.Center)
}

func (c CalendarStyle) RenderCell(style lipgloss.Style, isSelected bool) string {
	str := c.CharIcon
	if isSelected {
		str = c.CharSelected
	}
	return style.Render(str)
}

func (c CalendarStyle) RenderHeader() string {
	var out string

	for i, day := range c.DayStr {
		style := c.WeekdayStyle(true)

		if i == 0 || i == 6 {
			style = c.WeekendStyle(true)
		}

		out += style.Bold(true).Render(day)
	}

	return out
}
