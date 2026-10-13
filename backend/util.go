package backend

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

type Posn int

const (
	Top Posn = iota
	Mid
	Bot
)

func GenBorder(text string, borderStyle lipgloss.Style, posn Posn, width int) string {
	// text = style.Render(text)

	var paddingL, paddingR, left, right string
	pad := strings.Repeat("─", max(0, (width-lipgloss.Width(text))/2))
	if lipgloss.Width(text)%2 == 1 {
		paddingR = "─"
	}
	paddingL += pad
	paddingR += pad

	switch posn {
	case Top:
		left = "╭" + paddingL
		right = paddingR + "╮"
	case Mid:
		left = "├" + paddingL
		right = paddingR + "┤"
	case Bot:
		left = "╰" + paddingL
		right = paddingR + "╯"
	}

	return borderStyle.Render(left) + text + borderStyle.Render(right)
}
