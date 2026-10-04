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

	ta.Prompt = ""
	ta.Placeholder = "No Entry Found."
	ta.ShowLineNumbers = false
	ta.EndOfBufferCharacter = '~'
	ta.KeyMap = textarea.DefaultKeyMap()

	ta.SetValue(initText)
	ta.SetWidth(WIDTH)
	ta.SetHeight(HEIGHT)
	ta.Blur()
	style := GetDefaultStyle()

	return Notes{
		TextArea: ta,
		Style:    style,
	}
}

func (n Notes) Init() tea.Cmd {
	return nil
}

func (n Notes) View() tea.View {
	// May need to check the focus status
	n.TextArea.SetStyles(n.Style.Styles)
	return tea.NewView(n.TextArea.View())
}

func (n Notes) Update(msg tea.Msg) (Notes, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+w":
			n.ToggleFocus()
			return n, nil
		case "ctrl+c":
			return n, tea.Quit
		}
	}

	// Delegate all other messages (keystrokes, cursor blink ticks) to the textarea
	n.TextArea, cmd = n.TextArea.Update(msg)
	return n, cmd
}
