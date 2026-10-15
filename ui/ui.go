// Package ui holds the small rendering helpers every pane shares.
package ui

import (
	"image/color"
	"strings"

	"todo-tui/backend"

	"charm.land/lipgloss/v2"
)

// Fg returns a style with only the foreground color set.
func Fg(c backend.HexCode) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Color(c))
}

// Color converts a backend.HexCode into a lipgloss color.
func Color(c backend.HexCode) color.Color {
	return lipgloss.Color(string(c))
}

type Posn int

const (
	Top Posn = iota
	Mid
	Bot
)

// GenBorder draws a horizontal border line with text centered in it:
//
//	╭──── text ────╮
//
// width is the inner width (without the two corner glyphs). Pass text already
// styled, and a borderStyle for the line itself.
func GenBorder(text string, borderStyle lipgloss.Style, posn Posn, width int) string {
	tw := lipgloss.Width(text)
	pad := max(0, (width-tw)/2)
	left := strings.Repeat("─", pad)
	right := strings.Repeat("─", pad+tw%2) // an odd-width text gets the spare dash on the right

	switch posn {
	case Top:
		left, right = "╭"+left, right+"╮"
	case Mid:
		left, right = "├"+left, right+"┤"
	case Bot:
		left, right = "╰"+left, right+"╯"
	}

	return borderStyle.Render(left) + text + borderStyle.Render(right)
}

// PadRight pads s with trailing spaces to at least w cells.
func PadRight(s string, w int) string {
	return s + strings.Repeat(" ", max(0, w-lipgloss.Width(s)))
}

// PadLeft pads s with leading spaces to at least w cells.
func PadLeft(s string, w int) string {
	return strings.Repeat(" ", max(0, w-lipgloss.Width(s))) + s
}

// BarFill is how many of segs segments to fill. It never shows an empty bar
// once something is done, or a full one before everything is.
func BarFill(done, total, segs int) int {
	if total <= 0 {
		return 0
	}
	n := (done*segs + total/2) / total
	if done > 0 && n == 0 {
		n = 1
	}
	if done < total && n == segs {
		n = segs - 1
	}
	return n
}

// WrapJoin joins parts with sep, breaking into lines of at most width cells.
func WrapJoin(parts []string, sep string, width int) []string {
	var lines []string
	line := ""
	for _, p := range parts {
		switch {
		case line == "":
			line = p
		case lipgloss.Width(line)+lipgloss.Width(sep)+lipgloss.Width(p) <= width:
			line += sep + p
		default:
			lines = append(lines, line)
			line = p
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

// WrapText breaks text on spaces into lines of at most width cells, indenting
// every line by indent spaces. Words longer than a line are left whole.
func WrapText(text string, width, indent int) []string {
	pad := strings.Repeat(" ", indent)
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		switch {
		case line == "":
			line = word
		case indent+lipgloss.Width(line)+1+lipgloss.Width(word) <= width:
			line += " " + word
		default:
			lines = append(lines, pad+line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, pad+line)
	}
	return lines
}
