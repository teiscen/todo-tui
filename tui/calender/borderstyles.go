package calender

import (
	lipgloss "charm.land/lipgloss/v2"
)

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
