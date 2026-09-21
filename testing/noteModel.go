package testing

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type Model struct {
	textarea textarea.Model
}

// Need to populate the calenderPages from a config file
func initialModel() Model {
	ti := textarea.New()
	ti.Prompt = ""
	ti.ShowLineNumbers = true
	ti.EndOfBufferCharacter = '~'
	ti.MaxWidth = 30
	ti.SetHeight(30) // <- this is what actually controls the fixed height
	ti.KeyMap = textarea.DefaultKeyMap()

	ti.SetValue("This is a test of random values to see \n\n how it works \n\t\t wsdfsdfsdaf")

	styles := textarea.DefaultDarkStyles()
	styles.Cursor.Shape = tea.CursorUnderline
	styles.Cursor.Color = lipgloss.Color("205")
	ti.SetStyles(styles)
	ti.Focus()

	return Model{ti}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// Update styling now that we know the background color.
		// m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case "ctrl+c":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	const (
		footer = "\n(ctrl+c to quit)\n"
	)

	var c *tea.Cursor
	if !m.textarea.VirtualCursor() {
		c = m.textarea.Cursor()

		if c != nil {
			// Set the y offset of the cursor based on the position of the textarea
			// in the application.
			offset := lipgloss.Height("")
			c.Y += offset
		}
	}

	f := strings.Join([]string{
		"",
		m.textarea.View(),
		footer,
	}, "\n")

	v := tea.NewView(f)
	v.Cursor = c
	return v
}

func TestNotes() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
