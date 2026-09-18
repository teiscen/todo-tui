package calender

import (
	lipgloss "charm.land/lipgloss/v2"
)

// CharacterAlias
const (
	largeCircle    = "⬤"
	largeHexRing   = "⬡"
	smallCircle    = "●"
	smallRing      = "◯"
	selectedCircle = "◉"
)

// Colors
type ThemeColor int

const (
	ColorInvalid ThemeColor = iota
	ColorValid
	ColorRed
	ColorGreen
	ColorBlue
	ColorMagenta
	ColorOrange
	ColorTeal
)

func (c ThemeColor) toStyle() lipgloss.Style {
	base := lipgloss.NewStyle()

	switch c {
	case ColorInvalid:
		return base.Foreground(lipgloss.Color("#808080")) // Fixed hex length (6 digits)
	case ColorValid:
		return base.Foreground(lipgloss.Color("#C0C0C0"))
	case ColorBlue:
		return base.Foreground(lipgloss.Color("#6C93C7"))
	case ColorGreen:
		return base.Foreground(lipgloss.Color("#7FB380"))
	case ColorRed:
		return base.Foreground(lipgloss.Color("#C77373"))
	case ColorMagenta:
		return base.Foreground(lipgloss.Color("#B080B0"))
	case ColorOrange:
		return base.Foreground(lipgloss.Color("#D8A15C"))
	case ColorTeal:
		return base.Foreground(lipgloss.Color("#6BB3A3"))
	default:
		return base
	}
}

type Label int

const (
	LabelTask Label = iota
	LabelSleep
	LabelWorkout
	LabelReading
	LabelMisc1
	LabelMisc2
)

type LabelInfo struct {
	name  string
	color ThemeColor
}

var labelInfo = map[Label]LabelInfo{
	LabelTask:    {"Tasks", ColorRed},
	LabelMisc1:   {"Misc1", ColorGreen},
	LabelReading: {"Reading", ColorBlue},
	LabelWorkout: {"Workout", ColorMagenta},
	LabelSleep:   {"Sleep", ColorTeal},
	LabelMisc2:   {"Misc2", ColorOrange},
}

func (l Label) info() LabelInfo {
	return labelInfo[l]
}

func (l Label) name() string {
	return l.info().name
}

func (l Label) color() ThemeColor {
	return l.info().color
}

// CalenderStyling
type calenderStyle interface {
	getColumnWidth() int
	getRowSpacing() int
	weekdayRenderer() string
	invalidRenderer(c ThemeColor) string
	cellRenderer(c ThemeColor) string
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

func (s doubleCalenderStyle) cellRenderer(c ThemeColor) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Render(largeCircle)
}

func (s doubleCalenderStyle) selectedRenderer(c ThemeColor) string {
	content := "[" + largeCircle + " ]"
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

func (s singleCalenderStyle) cellRenderer(c ThemeColor) string {
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Render(smallCircle)
}

func (s singleCalenderStyle) selectedRenderer(c ThemeColor) string {
	// an use either: ◉ or >●
	return c.toStyle().
		Width(s.getColumnWidth()).Align(lipgloss.Center).
		Bold(true).Render(selectedCircle)
	// Bold(True).Render(">" + smallCircle)
}

// Border
func genBorder(
	mainBody string, borderColor ThemeColor,
	header string, headerColor ThemeColor, headerIsBold bool,
	footer string, footerColor ThemeColor, footerIsBold bool,
) string {
	padding := func(str string) (left string, right string) {
		dif := lipgloss.Width(mainBody) - lipgloss.Width(str)

		for range dif / 2 {
			left += "─"
		}
		right = left
		if dif%2 != 0 {
			right += "─"
		}
		return left, right
	}

	paddedHeader := " " + header + " "
	leftPadding, rightPadding := padding(paddedHeader)
	topBorderLeft := "╭" + leftPadding
	topBorderRight := rightPadding + "╮"
	headerStyled := headerColor.toStyle().Bold(headerIsBold).Render(paddedHeader)
	topBorder := topBorderLeft + headerStyled + topBorderRight

	paddedFooter := " " + footer + " "
	leftPadding, rightPadding = padding(paddedFooter)
	botBorderLeft := "╰" + leftPadding
	botBorderRight := rightPadding + "╯"
	footerStyled := footerColor.toStyle().Bold(footerIsBold).Render(paddedFooter)
	botBorder := botBorderLeft + footerStyled + botBorderRight

	sideBorder := borderColor.toStyle().
		Border(lipgloss.NormalBorder(), false, true, false, true).
		Render(mainBody)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, sideBorder, botBorder)
}
