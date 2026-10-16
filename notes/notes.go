package notes

import (
	backend "todo-tui/backend"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// WIDTH is the notes textarea width in cells, matched to the planner border.
const WIDTH = 26

// Notes wraps the textarea and its presentation style for the planner panel.
type Notes struct {
	TextArea textarea.Model
	Style    NotesStyle
}

// New creates a textarea configured for calendar-entry notes, sized to height.
func New(theme backend.Theme, initText string, height int) Notes {
	ta := textarea.New()

	ta.Prompt = ""
	ta.Placeholder = "No Entry Found."
	ta.ShowLineNumbers = false
	ta.EndOfBufferCharacter = '~'
	ta.KeyMap = textarea.DefaultKeyMap()

	ta.SetValue(initText)
	ta.SetWidth(WIDTH)
	ta.SetHeight(height)
	ta.Blur()
	style := NewStyle(theme)

	return Notes{
		TextArea: ta,
		Style:    style,
	}
}

func (n Notes) Init() tea.Cmd {
	return nil
}

// View applies the configured textarea styles and returns its rendered content.
func (n Notes) View() string {
	// May need to check the focus status
	n.TextArea.SetStyles(n.Style.Styles)
	return n.TextArea.View()
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

func (n *Notes) ChangeValue(str string) {
	n.TextArea.SetValue(str)
}

// SetAccent changes the accent color and regenerates the textarea styles.
func (n *Notes) SetAccent(c backend.HexCode) {
	n.Style.Theme.Accent = c
	n.Style.Styles = n.Style.GenStyle()
}
