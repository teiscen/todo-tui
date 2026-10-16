package leftbox

import "todo-tui/backend"

type LeftboxStyle struct {
	// Outer size including the border, to line up with the other columns.
	Width  int
	Height int

	// Column widths on a row. The title gets whatever is left.
	TagWidth  int
	ProgWidth int
	DueWidth  int

	CursorChar string
	CheckOpen  string
	CheckDone  string
	BarFull    string
	BarEmpty   string
	BarSegs    int

	Theme backend.Theme
}

// NewStyle builds a leftbox style from the shared theme.
func NewStyle(theme backend.Theme) LeftboxStyle {
	return LeftboxStyle{
		Width:  50,
		Height: 43,

		TagWidth:  4,
		ProgWidth: 10,
		DueWidth:  7,

		CursorChar: " ",
		CheckOpen:  "[ ]",
		CheckDone:  "[x]",
		BarFull:    "▰",
		BarEmpty:   "▱",
		BarSegs:    4,

		Theme: theme,
	}
}
