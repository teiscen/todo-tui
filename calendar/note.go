package calendar

import (
	"todo-tui/backend"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

var (
	MaxWidth  = 28
	MaxHeight = 23
)

type Note struct {
	Ta  textarea.Model
	nRI backend.NoteRenderInfo
}

func NewNote(s backend.State) Note {
	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Prompt = ""
	ti.SetWidth(MaxWidth)
	ti.SetHeight(MaxHeight)
	ti.KeyMap = textarea.DefaultKeyMap()

	styles := textarea.DefaultDarkStyles()
	styles.Cursor.Shape = tea.CursorUnderline

	ti.SetStyles(styles)

	nri := s.GetNoteRenderInfo()
	ti.SetValue(nri.Msg)
	return Note{
		ti,
		nri,
	}
}

func (n *Note) SetEntry() {
	n.Ta.SetValue(n.nRI.Msg)
}

// func (n *Note) Update(nRI backend.NoteRenderInfo) {
func (n *Note) Update(s backend.State) {
	n.nRI = s.GetNoteRenderInfo()

	if n.nRI.IsFocused {
		n.Ta.Focus()
	} else {
		n.Ta.Blur()
	}
}

func (n Note) Render() string {
	if n.nRI.IsFocused {
		n.Ta.Focus()
	} else {
		n.Ta.Blur()
	}
	return "\n" + n.Ta.View()
}
