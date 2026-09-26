package calendar

import (
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

func getNote() Note {
	ti := textarea.New()
	ti.ShowLineNumbers = false
	ti.Prompt = ""
	ti.SetWidth(MaxWidth)
	ti.SetHeight(MaxHeight)

	ti.KeyMap = textarea.DefaultKeyMap()

	styles := textarea.DefaultDarkStyles()
	styles.Cursor.Shape = tea.CursorUnderline

	// // Remove gutter/prompt padding completely
	// styles.Focused.Prompt = styles.Focused.Prompt.Width(0).Margin(0).Padding(0)
	// styles.Blurred.Prompt = styles.Blurred.Prompt.Width(0).Margin(0).Padding(0)
	//
	// // Flush line numbers to the hard left edge
	// styles.Focused.LineNumber = styles.Focused.LineNumber.Margin(0).Padding(0).Align(lipgloss.Left)
	// styles.Blurred.LineNumber = styles.Blurred.LineNumber.Margin(0).Padding(0).Align(lipgloss.Left)
	//
	// // Remove base textarea outer padding if present
	// styles.Focused.Base = styles.Focused.Base.Margin(0).Padding(0)
	// styles.Blurred.Base = styles.Blurred.Base.Margin(0).Padding(0)

	ti.SetStyles(styles)

	ti.SetValue("DEFAULT VALUE .................................................")
	return Note{ti}
}

func Init() {}

func (n Note) Render() string {
	return "\n" + n.ta.View()
}
