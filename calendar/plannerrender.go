package calendar

import (
	"todo-tui/backend"

	lipgloss "charm.land/lipgloss/v2"
)

//      Thick
// ┏━━━━━━━━━━━━━━┓
// ┃  Bubble Tea  ┃
// ┗━━━━━━━━━━━━━━┛

func genBorderCalendar(
	mainBody string, // borderColor Color,
	// header string, headerColor Color,
	// footer string, footerColor Color,
) string {
	header := "HEADER"
	footer := "FOOTER"
	borderColor := backend.HexCode("#C0C0C0")
	headerColor := backend.HexCode("#C0C0C0")
	footerColor := backend.HexCode("#C0C0C0")

	headerIsBold := true
	footerIsBold := true

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
	paddedFooter := " " + footer + " "

	headerStyled := headerColor.ToStyle().Bold(headerIsBold).Render(paddedHeader)
	leftPadding, rightPadding := padding(paddedHeader)

	topBorderLeft := "╭" + leftPadding
	topBorderRight := rightPadding + "╮"

	footerStyled := footerColor.ToStyle().Bold(footerIsBold).Render(paddedFooter)
	leftPadding, rightPadding = padding(paddedFooter)

	botBorderLeft := "├" + leftPadding
	botBorderRight := rightPadding + "┤"

	topBorder := topBorderLeft + headerStyled + topBorderRight
	botBorder := botBorderLeft + footerStyled + botBorderRight

	sideBorder := borderColor.ToStyle().
		Border(lipgloss.NormalBorder(), false, true, false, true).
		Render(mainBody)

	str := lipgloss.JoinVertical(lipgloss.Left, topBorder, sideBorder, botBorder)
	return lipgloss.NewStyle().Padding(0, 1).Render(str)
}

func genBorderNotes(note string) string {
	borderColor := backend.HexCode("#C0C0C0")
	return borderColor.ToStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Padding(0, 1).
		Render(note)
}

func (Planner) Render(cal string, note string) string {
	cal = genBorderCalendar(cal)
	note = genBorderNotes(note)
	combined := lipgloss.JoinVertical(lipgloss.Center, cal, note)
	borderColor := backend.HexCode("#C0C0C0")
	return borderColor.ToStyle().
		Border(lipgloss.RoundedBorder()).
		// Padding(1, 2).
		Render(combined)
}
