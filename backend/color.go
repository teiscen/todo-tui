package backend

import lipgloss "charm.land/lipgloss/v2"

type (
	ColorID string
	HexCode string
	Colors  map[ColorID]HexCode

	ColoredString struct {
		str   string
		cInfo HexCode
	}
)

// Temporary
func (c HexCode) ToStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(string(c)))
}
