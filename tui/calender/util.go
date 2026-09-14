package calender

import (
	"math"
	"time"

	"charm.land/lipgloss/v2"
)

// Style Related
func getStyleForDay(day int, cursor int, hasEntry bool, color ThemeColor) lipgloss.Style {
	var style lipgloss.Style
	if hasEntry {
		style = color.toStyle()
	} else {
		style = ColorValid.toStyle()
	}
	return style.Underline(day == cursor)
}

// Time related
func getFirstDay() time.Weekday {
	firstDay := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	return firstDay.Weekday()
}

func getLastDay() time.Weekday {
	lastDay := time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.Now().Location())
	return lastDay.Weekday()
}

func getNumOfDays() int {
	return time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.Now().Location()).Day()
}

func getNumOfRows() int {
	totalDays := float64(getNumOfDays() + int(getFirstDay()))
	return int(math.Ceil(totalDays / 7.0))
}

/*
func (m model) ViewOld() tea.View {
	var days string
	// Days of Week
	days += ColorRed.toStyle().Render("S ")
	days += ColorValid.toStyle().Render("M T W T F ")
	days += ColorRed.toStyle().Render("S")
	days += lipgloss.NewStyle().Render("\n")
	// Days from the previous months in the week
	count := 0
	for i := 0; i < int(getFirstDay()); i++ {
		days += ColorInvalid.toStyle().Render(smallCircle) + " "
		count++
	}
	// Days of the month
	for day := 1; day <= getNumOfDays(); day++ {
		// days += getStyleForDay(day, cursor int, )
		days += ColorValid.toStyle().Render(smallCircle) + " "
		count++
		if count%7 == 0 {
			days += "\n"
		}
	}
	// Days from the next month in the week
	for i := 0; i < int(getLastDay()); i++ {
		days += ColorInvalid.toStyle().Render(smallCircle) + " "
	}
	days += "\n"

	return tea.NewView(days)
}
*/
