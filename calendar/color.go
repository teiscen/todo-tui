package calendar

import (
	lipgloss "charm.land/lipgloss/v2"
)

type Color string

func (c Color) toStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(string(c)))
}
