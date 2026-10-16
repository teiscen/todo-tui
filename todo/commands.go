package todo

import (
	"todo-tui/backend"
	"todo-tui/database"

	tea "charm.land/bubbletea/v2"
)

// plannerLoadedMsg carries the grid entries for a requested month/label. It
// echoes the request so the root can drop stale results when the user has
// already moved on to another month or label.
type plannerLoadedMsg struct {
	LabelID backend.LabelID
	NewDate backend.Date
	sDate   backend.Date
	entries map[string]string
}

// errMsg reports a failed asynchronous load.
type errMsg struct{ err error }

// loadPlanner fetches the calendar grid for a label/month off the main loop.
func (t *Todo) loadPlanner(lID backend.LabelID, date backend.Date) tea.Cmd {
	return func() tea.Msg {
		// Start Date is the first Sunday from the first of the month.
		// End Date is 41 days after it.
		first := date.FirstDay()
		sDate := first.AddDate(0, 0, -first.Weekday())
		eDate := sDate.AddDate(0, 0, 41)

		entries, err := database.FetchGridEntries(t.db, sDate, eDate, lID)
		if err != nil {
			return errMsg{err: err}
		}
		return plannerLoadedMsg{LabelID: lID, NewDate: date, sDate: sDate, entries: entries}
	}
}

func (t *Todo) AddEntry(lID backend.LabelID, date backend.Date, msg string) {
	if err := database.AddEntry(t.db, date.FormatSQL(), string(lID), msg); err != nil {
		t.err = err
	}
}

func (t *Todo) RemoveEntry(lID backend.LabelID, date backend.Date) {
	if err := database.RemoveEntry(t.db, date.FormatSQL(), string(lID)); err != nil {
		t.err = err
	}
}
