package calender

import "charm.land/lipgloss/v2"

type calenderStyle interface {
	getColumnWidth() int
	getRowSpacing() int
	weekdayRenderer() string
	invalidRenderer(c ThemeColor) string
	cellRenderer(c ThemeColor, isBold bool) string
	selectedRenderer(c ThemeColor) string
}

type (
	singleCalenderStyle struct{}
	doubleCalenderStyle struct{}
)

func (doubleCalenderStyle) getColumnWidth() int {
	return 4
}

func (doubleCalenderStyle) getRowSpacing() int {
	return 1
}

func (s doubleCalenderStyle) weekdayRenderer() string {
	var out string
	temp := func(c ThemeColor, str string) string {
		out := c.toStyle().
			Width(s.getColumnWidth()).Align(lipgloss.Center).
			Bold(true).Render(str)
		return out
	}
	out += "\n"
	out += temp(ColorRed, "Su")
	out += temp(ColorValid, "Mo")
	out += temp(ColorValid, "Tu")
	out += temp(ColorValid, "We")
	out += temp(ColorValid, "Th")
	out += temp(ColorValid, "Fr")
	out += temp(ColorRed, "Sa")
	out += "\n"
	return out
}

func (s doubleCalenderStyle) invalidRenderer(c ThemeColor) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Render(largeHexRing)
}

func (s doubleCalenderStyle) cellRenderer(c ThemeColor, isBold bool) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).Bold(isBold).
		Render(largeCircle)
}

func (s doubleCalenderStyle) selectedRenderer(c ThemeColor) string {
	content := ">" + largeCircle + "<"
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Bold(true).Render(content)
}

// Single Calender Styles
func (singleCalenderStyle) getColumnWidth() int {
	return 3
}

func (singleCalenderStyle) getRowSpacing() int {
	return 0
}

func (s singleCalenderStyle) weekdayRenderer() string {
	var out string
	temp := func(c ThemeColor, str string) string {
		out := c.toStyle().
			Width(s.getColumnWidth()).Align(lipgloss.Center).
			Bold(false).Render(str)
		return out
	}
	out += "\n"
	out += temp(ColorRed, "S")
	out += temp(ColorValid, "M")
	out += temp(ColorValid, "T")
	out += temp(ColorValid, "W")
	out += temp(ColorValid, "T")
	out += temp(ColorValid, "F")
	out += temp(ColorRed, "S")
	out += "\n"
	return out
}

func (s singleCalenderStyle) invalidRenderer(c ThemeColor) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Render(smallRing)
}

func (s singleCalenderStyle) cellRenderer(c ThemeColor, isBold bool) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).Bold(isBold).
		Render(smallCircle)
}

func (s singleCalenderStyle) selectedRenderer(c ThemeColor) string {
	// an use either: ◉ or >●
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Bold(true).Render(selectedCircle)
	// Bold(True).Render(">" + smallCircle)
}
