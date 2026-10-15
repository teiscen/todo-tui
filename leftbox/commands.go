package leftbox

import tea "charm.land/bubbletea/v2"

// ToggleTodoMsg: space was pressed. The root writes it to the database
// (done_at for one-time todos, a completions row for recurring ones) and reloads.
type ToggleTodoMsg struct {
	ID        int
	Recurring bool
	Done      bool // the state it is in now, before the toggle
}

func (l Leftbox) toggleCmd() tea.Cmd {
	r, ok := l.Selected()
	if !ok {
		return nil
	}
	return func() tea.Msg {
		return ToggleTodoMsg{ID: r.ID, Recurring: r.Tab == TabRecurring, Done: r.Done}
	}
}
