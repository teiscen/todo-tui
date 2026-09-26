package main

import (
	"fmt"
	"os"

	"todo-tui/backend"
	"todo-tui/calendar"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	p calendar.Planner
}

func (m *Model) Initialize(s backend.State) {
	m.p = calendar.NewPlanner(s)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
}

func (m Model) View() tea.View {
	str := m.p.Render()
	return tea.NewView(str)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "shift-c":
			return m, tea.Quit

		case "tab":
			// m.focus = m.focus.toggleFocus()
			// if m.focus == FocusCal {
			// 	m.notes.Blur()
			// } else {
			// 	cmd = m.notes.Focus()
			// 	cmds = append(cmds, cmd)
			// }
		}
	}

	// if m.focus == FocusNote {
	// 	m.notes, cmd = m.notes.Update(msg)
	// 	cmds = append(cmds, cmd)
	// }

	// m = m.reconcile()

	return m, tea.Batch(cmds...)
}

func TestRendering() {
	state := backend.MockState()

	model := Model{}
	model.Initialize(state)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
