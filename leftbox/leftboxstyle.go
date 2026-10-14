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

	ColorBorder backend.HexCode
	ColorText   backend.HexCode
	ColorDim    backend.HexCode
	ColorAccent backend.HexCode
	ColorLate   backend.HexCode
	ColorSoon   backend.HexCode
}

func GetDefaultStyle() LeftboxStyle {
	return LeftboxStyle{
		Width:  50,
		Height: 43,

		TagWidth:  4,
		ProgWidth: 10,
		DueWidth:  7,

		CursorChar: "▸",
		CheckOpen:  "[ ]",
		CheckDone:  "[x]",
		BarFull:    "▰",
		BarEmpty:   "▱",
		BarSegs:    4,

		ColorBorder: backend.HexCode("#C6D0F5"),
		ColorText:   backend.HexCode("#C6D0F5"),
		ColorDim:    backend.HexCode("#939AB7"),
		ColorAccent: backend.HexCode("#CA9EE6"),
		ColorLate:   backend.HexCode("#E78284"),
		ColorSoon:   backend.HexCode("#E5C890"),
	}
}

func (s LeftboxStyle) ToneColor(t Tone) backend.HexCode {
	switch t {
	case ToneLate:
		return s.ColorLate
	case ToneSoon:
		return s.ColorSoon
	}
	return s.ColorDim
}
