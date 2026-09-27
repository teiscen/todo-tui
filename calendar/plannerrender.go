package calendar

import (
	"todo-tui/backend"

	lipgloss "charm.land/lipgloss/v2"
)

//      Thick
// ┏━━━━━━━━━━━━━━┓
// ┃  Bubble Tea  ┃
// ┗━━━━━━━━━━━━━━┛

func (p Planner) genBorderCalendar(
	mainBody string,
) string {
	header := p.pRI.Header
	footer := p.pRI.Footer
	borderColor := p.pRI.BorderColor

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

	paddedHeader := " " + header.Str + " "
	paddedFooter := " " + footer.Str + " "

	headerStyled := header.Hex.ToStyle().Bold(headerIsBold).Render(paddedHeader)
	leftPadding, rightPadding := padding(paddedHeader)

	topBorderLeft := "╭" + leftPadding
	topBorderRight := rightPadding + "╮"

	footerStyled := footer.Hex.ToStyle().Bold(footerIsBold).Render(paddedFooter)
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

func (p Planner) genBorderNotes(note string) string {
	borderColor := p.pRI.BorderColor
	return borderColor.ToStyle().
		Border(lipgloss.RoundedBorder(), false, true, true, true).
		Padding(0, 1).
		Render(note)
}

func (p Planner) RenderBorder(cal string, note string) string {
	cal = p.genBorderCalendar(cal)
	note = p.genBorderNotes(note)
	combined := lipgloss.JoinVertical(lipgloss.Center, cal, note)
	borderColor := backend.HexCode("#C0C0C0")
	return borderColor.ToStyle().
		Border(lipgloss.RoundedBorder()).
		// Padding(1, 2).
		Render(combined)
}
