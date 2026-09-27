package main

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit

	case "tab":
		// return m, m.SwapPane()
	}

	switch m.focus {
	case PaneCalendar:
		return m.handleCalendarKey(msg)

	case PaneNote:
		return m.handleNoteKey(msg)
	}

	m.p.Update(m.s)
	return m, nil
}

func (m Model) handleCalendarKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "h":
		m.CalendarMove(Left)
	case "j":
		m.CalendarMove(Down)
	case "k":
		m.CalendarMove(Up)
	case "l":
		m.CalendarMove(Right)
	case "ctrl+h":
		m.CalendarMove(PrevLabel)
	case "ctrl+j":
		m.CalendarMove(NextMonth)
	case "ctrl+k":
		m.CalendarMove(PrevMonth)
	case "ctrl+l":
		m.CalendarMove(NextLabel)
	case "ctrl+w":
		m.s.NoteFocused = true
		m.focus = PaneNote
		m.p.Update(m.s)
		return m, nil
	}

	m.p.Update(m.s)
	m.p.Note.SetEntry()

	return m, nil
}

func (m Model) handleNoteKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+w":
		m.s.NoteFocused = false
		m.focus = PaneCalendar
		m.p.Update(m.s)
		return m, nil
	}

	var cmd tea.Cmd
	m.p.Note.Ta, cmd = m.p.Note.Ta.Update(msg)

	return m, cmd
}

type CalendarMovement int

const (
	Left CalendarMovement = iota
	Down
	Up
	Right

	PrevMonth
	NextMonth

	PrevLabel
	NextLabel
)

func (m *Model) CalendarMove(cm CalendarMovement) {
	date := m.s.SelectedDate
	switch cm {
	// Calendar General Movements
	case Left:
		m.s.SelectedDate = date.AddDate(0, 0, -1)
	case Up:
		m.s.SelectedDate = date.AddDate(0, 0, -7)
	case Down:
		m.s.SelectedDate = date.AddDate(0, 0, 7)
	case Right:
		m.s.SelectedDate = date.AddDate(0, 0, 1)

	// Monthly movement
	// Bug related to end of month monvement maybe
	case PrevMonth:
		m.s.SelectedDate = date.AddDate(0, -1, 0)
	case NextMonth:
		m.s.SelectedDate = date.AddDate(0, 1, 0)
	// Label Movement
	case PrevLabel:
		m.s.SelectedLabel = m.s.PrevLabel()
	case NextLabel:
		m.s.SelectedLabel = m.s.NextLabel()
	}
}
