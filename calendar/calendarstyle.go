package calendar

import (
	backend "todo-tui/backend"
	"todo-tui/ui"

	"charm.land/lipgloss/v2"
)

type CalendarStyle struct {
	ColumnWidth int
	RowSpacing  int

	DayStr []string

	CharIcon     string
	CharSelected string

	Theme backend.Theme
}

// NewStyle builds a calendar style from the shared theme.
func NewStyle(theme backend.Theme) CalendarStyle {
	return CalendarStyle{
		ColumnWidth: 4,
		RowSpacing:  1,

		DayStr: []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"},

		CharIcon:     "⬤",
		CharSelected: "|⬤ |",

		Theme: theme,
	}
}

func (c CalendarStyle) WeekdayStyle(isValid bool) lipgloss.Style {
	color := c.Theme.Weekday
	if !isValid {
		color, _ = c.Theme.Weekday.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return ui.Fg(color).Width(c.ColumnWidth).Align(lipgloss.Center)
}

func (c CalendarStyle) WeekendStyle(isValid bool) lipgloss.Style {
	color := c.Theme.Weekend
	if !isValid {
		color, _ = c.Theme.Weekend.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return ui.Fg(color).Width(c.ColumnWidth).Align(lipgloss.Center)
}

func (c CalendarStyle) AccentStyle(isValid bool) lipgloss.Style {
	color := c.Theme.Accent
	if !isValid {
		color, _ = c.Theme.Accent.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	}
	return ui.Fg(color).Width(c.ColumnWidth).Align(lipgloss.Center)
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
