package calender

import (
	"time"

	lipgloss "charm.land/lipgloss/v2"
)

const (
	largeCircle    = "⬤"
	largeHexRing   = "⬡"
	smallCircle    = "●"
	smallRing      = "◯"
	selectedCircle = "◉"

	arrow = "➤"
)

type ThemeColor int

const (
	ColorInvalid ThemeColor = iota
	ColorValid
	ColorBlue
	ColorGreen
	ColorRed
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

// MOCKS
var sampleEntries = []time.Time{
	time.Date(2006, time.September, 1, 9, 30, 0, 0, time.UTC),   // Sep 1, 9:30 AM
	time.Date(2006, time.September, 5, 14, 15, 0, 0, time.UTC),  // Sep 5, 2:15 PM
	time.Date(2006, time.September, 12, 18, 0, 0, 0, time.UTC),  // Sep 12, 6:00 PM
	time.Date(2006, time.September, 18, 8, 45, 0, 0, time.UTC),  // Sep 18, 8:45 AM
	time.Date(2006, time.September, 22, 12, 0, 0, 0, time.UTC),  // Sep 22, 12:00 PM
	time.Date(2006, time.September, 28, 20, 30, 0, 0, time.UTC), // Sep 28, 8:30 PM
}

var mockCalenderPages = []calenderPage{
	{
		color: ColorOrange,
		entries: []time.Time{
			sampleEntries[0],
			sampleEntries[1],
		},
	},
	{
		color: ColorMagenta,
		entries: []time.Time{
			sampleEntries[2],
		},
	},
	{
		color: ColorTeal,
		entries: []time.Time{
			sampleEntries[3],
			sampleEntries[4],
			sampleEntries[5],
		},
	},
	{color: ColorRed, entries: []time.Time{}},
	{color: ColorGreen, entries: []time.Time{}},
	{color: ColorBlue, entries: []time.Time{}},
}
