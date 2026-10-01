package notes

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

var (
	HEIGHT = 23
	WIDTH  = 28
)

type Notes struct {
	TextArea textarea.Model
	Style    NotesStyle
}

func NewNotes(initText string) Notes {
	ta := textarea.New()
	ta.SetValue(initText)
	ta.SetWidth(WIDTH)
	ta.SetHeight(HEIGHT)
	style := GetDefaultStyle()

	ta.Prompt = ""
	ta.Placeholder = "No Entry Found."
	ta.ShowLineNumbers = false
	ta.KeyMap = textarea.DefaultKeyMap()

	return Notes{
		TextArea: ta,
		Style:    style,
	}
}

func (n Notes) Render() string {
	// May need to check the focus status
	n.TextArea.SetStyles(n.Style.Styles)
	return n.TextArea.View()
}

func (n *Notes) Update(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "ctrl+w":
		(*n).ToggleFocus()
	}
	n.TextArea, _ = n.TextArea.Update(msg)
}

// MODELS require
// func (m Model) Init() tea.Cmd {
// func (m Model) View() tea.View {
// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//
// func (n *Notes) Update(msg tea.Msg) tea.Cmd {
// 	var cmd tea.Cmd
//
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+w":
// 			if n.TextArea.Focused() {
// 				n.TextArea.Blur()
// 				return nil
// 			}
// 			// Focus() returns tea.Blink command to start the cursor animation
// 			return n.TextArea.Focus()
// 		}
// 	}
//
// 	// Delegate all other messages (keystrokes, cursor blink ticks) to the textarea
// 	n.TextArea, cmd = n.TextArea.Update(msg)
// 	return cmd
// }
