package planner

import (
	"charm.land/lipgloss/v2"
)

func (p Planner) genCalendarBorder() string {
	headerStyle := p.CalendarModel.GetSelectedColor().ToStyle().Padding(0, 1)
	footerStyle := p.Calendar.GetCurrentLabel().Color.ToStyle().Padding(0, 1)

	headerText := p.CalendarModel.SelectedDate.Format()
	footerText := p.Calendar.GetCurrentLabel().Name

	headerFormatted := headerStyle.Render(headerText)
	footerFormatted := footerStyle.Render(footerText)

	width := 30
	padding := func(str string) string {
		var out string
		for range (width - lipgloss.Width(str)) / 2 {
			out += "-"
		}
		return out
	}

	headerPadding := padding(headerFormatted)
	footerPadding := padding(headerFormatted)

	headerFull := "╭" + headerPadding + headerFormatted + headerPadding + "╮"
	footerFull := "├" + footerPadding + footerFormatted + footerPadding + "┤"

	notesView := p.NotesModel.View().Content
	calendarView := p.CalendarModel.View().Content

	calendarBorder := p.Style.BorderColor.ToStyle().
		Border(lipgloss.NormalBorder(), false, true).Render(calendarView)

	notesBorder := p.Style.BorderColor.ToStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).Render(notesView)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerFull,
		calendarBorder,
		footerFull,
		notesBorder,
	)
}
