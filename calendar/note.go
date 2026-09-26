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
	ta textarea.Model
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

	ti.SetValue("DEFAULT VALUE .................................................")
	return Note{ti}
}

func (n Note) Render() string {
	return "\n" + n.ta.View()
}
