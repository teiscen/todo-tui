package calendar

import (
	backend "todo-tui/backend"

	tea "charm.land/bubbletea/v2"
)

// Mark/Unmark
// Movement
type direction int

const (
	Left direction = iota
	Down
	Up
	Right
)

type (
	ChangeMonthMsg struct{ NewDate backend.Date }
	ChangeDateMsg  struct{ NewDate backend.Date }
	ToggleMarkMsg  struct {
		Date   backend.Date
		Status bool
	}
)

func (c *Calendar) Movement(dir direction) tea.Cmd {
	oldDate := c.SelectedDate
	switch dir {
	case Left:
		c.SelectedDate = oldDate.AddDate(0, 0, -1)
	case Down:
		c.SelectedDate = oldDate.AddDate(0, 0, 7)
	case Up:
		c.SelectedDate = oldDate.AddDate(0, 0, -7)
	case Right:
		c.SelectedDate = oldDate.AddDate(0, 0, 1)
	}

	if oldDate.Month() != c.SelectedDate.Month() {
		return func() tea.Msg {
			return ChangeMonthMsg{NewDate: c.SelectedDate}
		}
	} else {
		return func() tea.Msg {
			return ChangeDateMsg{NewDate: c.SelectedDate}
		}
	}
}

func (c *Calendar) ToggleMark() tea.Cmd {
	infoAddr := c.GetInfo(c.SelectedDate)
	if infoAddr == nil {
		return nil
	}
	infoAddr.IsMarked = !infoAddr.IsMarked
	return func() tea.Msg {
		return ToggleMarkMsg{
			Date:   c.SelectedDate,
			Status: infoAddr.IsMarked,
		}
	}
}
