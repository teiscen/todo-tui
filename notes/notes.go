package notes

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

var (
	HEIGHT = 25
	WIDTH  = 26
)

// Notes wraps the textarea and its presentation style for the planner panel.
type Notes struct {
	TextArea textarea.Model
	Style    NotesStyle
}

// NewNotes creates a textarea configured for calendar-entry notes.
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

// View applies the configured textarea styles and returns its rendered content.
func (n Notes) View() tea.View {
	// May need to check the focus status
	n.TextArea.SetStyles(n.Style.Styles)
	return tea.NewView(n.TextArea.View())
}

// Update handles notes-specific shortcuts and delegates editing to the textarea.
func (n Notes) Update(msg tea.Msg) (Notes, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return n, func() tea.Msg {
				return WriteTextMsg{
					Text: n.TextArea.Value(),
				}
			}
		case "ctrl+c":
			return n, tea.Quit
		}
	}

	// Delegate all other messages (keystrokes, cursor blink ticks) to the textarea
	n.TextArea, cmd = n.TextArea.Update(msg)
	return n, cmd
}
