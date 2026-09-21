package calendar

import lipgloss "charm.land/lipgloss/v2"

//      Thick
// ┏━━━━━━━━━━━━━━┓
// ┃  Bubble Tea  ┃
// ┗━━━━━━━━━━━━━━┛

func genBorderCalendar(
	mainBody string, borderColor Color,
	header string, headerColor Color,
	footer string, footerColor Color,
) string {
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

	headerStyled := headerColor.toStyle().Bold(headerIsBold).Render(paddedHeader)
	leftPadding, rightPadding := padding(paddedHeader)

	topBorderLeft := "╭" + leftPadding
	topBorderRight := rightPadding + "╮"

	footerStyled := footerColor.toStyle().Bold(footerIsBold).Render(paddedFooter)
	leftPadding, rightPadding = padding(paddedFooter)

	botBorderLeft := "├" + leftPadding
	botBorderRight := rightPadding + "┤"

	topBorder := topBorderLeft + headerStyled + topBorderRight
	botBorder := botBorderLeft + footerStyled + botBorderRight

	sideBorder := borderColor.toStyle().
		Border(lipgloss.NormalBorder(), false, true, false, true).
		Render(mainBody)

	str := lipgloss.JoinVertical(lipgloss.Left, topBorder, sideBorder, botBorder)
	return lipgloss.NewStyle().Padding(0, 1).Render(str)
}

func genBorderNotes(note string, borderColor Color) string {
	return borderColor.toStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Padding(0, 1).
		Render(note)
}

func genBorder(cal string, note string, borderColor Color) string {
	combined := lipgloss.JoinVertical(lipgloss.Center, cal, note)
	return borderColor.toStyle().
		Border(lipgloss.RoundedBorder()).
		// Padding(1, 2).
		Render(combined)
}
