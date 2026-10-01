package notes

import (
	backend "todo-tui/backend_old"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type NotesStyle struct {
	Styles textarea.Styles

	ColorText   backend.HexCode
	ColorCursor backend.HexCode
	ColorElse   backend.HexCode // Meant for anything extraneous such as Line numbers or Buffer
}

func (n NotesStyle) GenStyle() textarea.Styles {
	cText := n.ColorText
	cElse := n.ColorElse
	cCursor := n.ColorCursor

	cTextMute, err := cText.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	if err != nil {
		cTextMute = cText
	}
	cElseMute, err := cElse.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	if err != nil {
		cElseMute = cElse
	}
	// cCursorMute, err := cCursor.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	// if err != nil {
	// 	cCursorMute = cCursor
	// }

	focused := textarea.StyleState{
		Base:             lipgloss.NewStyle(), // Leave default and then later match the overall color scheme
		Text:             cText.ToStyle(),
		LineNumber:       lipgloss.NewStyle(), // Not Used
		CursorLineNumber: lipgloss.NewStyle(), // Not Used
		CursorLine:       lipgloss.NewStyle().Background(lipgloss.Color(string(cElse))),
		EndOfBuffer:      lipgloss.NewStyle(), // Not Used
		Placeholder:      cElse.ToStyle(),
		Prompt:           lipgloss.NewStyle(), // Not Used
	}

	blurred := textarea.StyleState{
		Base:             lipgloss.NewStyle(), // Leave default and then later match the overall color scheme
		Text:             cTextMute.ToStyle(),
		LineNumber:       lipgloss.NewStyle(), // Not Used
		CursorLineNumber: lipgloss.NewStyle(), // Not Used
		CursorLine:       lipgloss.NewStyle().Background(lipgloss.Color(string(cElseMute))),
		EndOfBuffer:      lipgloss.NewStyle(), // Not Used
		Placeholder:      cElseMute.ToStyle(),
		Prompt:           lipgloss.NewStyle(), // Not Used
	}
	cursor := textarea.CursorStyle{
		Color: lipgloss.Color(string(cCursor)),
		Shape: tea.CursorBar, // tea.CursorBlock, tea.CursorUnderline
		Blink: false,
		// BlinkSpeed: 500 * time.Millisecond, // Only matters if Blink is true
	}

	return textarea.Styles{
		Focused: focused,
		Blurred: blurred,
		Cursor:  cursor,
	}
}

func GetDefaultStyle() NotesStyle {
	cText := backend.HexCode("#C6D0F5")
	cElse := backend.HexCode("#CA9EE6")
	cCursor := backend.HexCode("#939AB7")

	ns := NotesStyle{
		Styles: textarea.Styles{},

		ColorText:   cText,
		ColorCursor: cCursor,
		ColorElse:   cElse,
	}

	ns.Styles = ns.GenStyle()
	return ns
}
