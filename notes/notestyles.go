package notes

import (
	backend "todo-tui/backend"
	"todo-tui/ui"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type NotesStyle struct {
	Styles textarea.Styles

	Theme backend.Theme
}

func (n NotesStyle) GenStyle() textarea.Styles {
	cText := n.Theme.Text
	cAccent := n.Theme.Accent

	cTextMute, err := cText.MuteColorTarget(backend.HexCode("#303446"), 0.4)
	if err != nil {
		cTextMute = cText
	}
	cAccentMute, err := cAccent.MuteColorTarget(backend.HexCode("#303446"), 0.8)
	if err != nil {
		cAccentMute = cAccent
	}

	focused := textarea.StyleState{
		Base:             lipgloss.NewStyle(), // Leave default and then later match the overall color scheme
		Text:             ui.Fg(cText),
		LineNumber:       lipgloss.NewStyle(), // Not Used
		CursorLineNumber: lipgloss.NewStyle(), // Not Used
		CursorLine:       lipgloss.NewStyle().Background(lipgloss.Color(string(cAccentMute))),
		EndOfBuffer:      ui.Fg(cAccent),
		Placeholder:      ui.Fg(cAccent),
		Prompt:           lipgloss.NewStyle(), // Not Used
	}

	blurred := textarea.StyleState{
		Base:             lipgloss.NewStyle(), // Leave default and then later match the overall color scheme
		Text:             ui.Fg(cTextMute),
		LineNumber:       lipgloss.NewStyle(), // Not Used
		CursorLineNumber: lipgloss.NewStyle(), // Not Used
		CursorLine:       lipgloss.NewStyle().Background(lipgloss.Color(string(cAccentMute))),
		EndOfBuffer:      ui.Fg(cAccentMute),
		Placeholder:      ui.Fg(cAccent),
		Prompt:           lipgloss.NewStyle(), // Not Used
	}
	cursor := textarea.CursorStyle{
		Color: lipgloss.Color(string(cAccent)),
		Shape: tea.CursorBar, // tea.CursorBlock, tea.CursorUnderline
		Blink: false,
	}

	return textarea.Styles{
		Focused: focused,
		Blurred: blurred,
		Cursor:  cursor,
	}
}

// NewStyle builds a notes style from the shared theme.
func NewStyle(theme backend.Theme) NotesStyle {
	ns := NotesStyle{
		Styles: textarea.Styles{},
		Theme:  theme,
	}

	ns.Styles = ns.GenStyle()
	return ns
}
