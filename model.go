package main

import (
	"todo-tui/backend"
	"todo-tui/calendar"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Pane int

const (
	PaneCalendar Pane = iota
	PaneNote
)

type Model struct {
	p     calendar.Planner
	focus Pane
	s     backend.State
}

func (m *Model) Initialize(s backend.State) {
	m.p = calendar.NewPlanner(s)
	m.focus = PaneCalendar
	m.s = s
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
}

func (m Model) View() tea.View {
	str := m.p.Render()
	return tea.NewView(str)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	if m.focus == PaneNote {
		var cmd tea.Cmd
		// m.p.Note.Ta, cmd = m.p.Note.Ta.Update(msg)
		return m, cmd
	}

	return m, nil
}
