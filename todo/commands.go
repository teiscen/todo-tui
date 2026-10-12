package todo

import (
	"todo-tui/backend"
	"todo-tui/database"
)

func (t *Todo) ReloadPlannerState(lID backend.LabelID, date backend.Date) {
	// Start Date is the first Sunday from the first of the month
	// End Date is 41 days after it.
	sDate := date
	sDate.FirstDay()
	sDate.ChangeDate(0, 0, -sDate.Weekday())
	eDate := sDate.AddDate(0, 0, 41)

	entries, _ := database.FetchGridEntries(
		t.db, sDate, eDate, lID,
	)
	t.plannerModel.ReloadState(sDate, entries)
}

func (t *Todo) AddEntry(lID backend.LabelID, date backend.Date, msg string) {
	database.AddEntry(t.db, date.FormatSQL(), string(lID), msg)
}

func (t *Todo) RemoveEntry(lID backend.LabelID, date backend.Date) {
	database.RemoveEntry(t.db, date.FormatSQL(), string(lID))
}

// loadEntries fetches the visible month for the selected label and rebuilds
// the calendar grid. On error the grid is rebuilt with no marks.
// func (t *Todo) loadEntries() {
// 	cal := &t.plannerModel.CalendarModel
//
// 	entries, err := database.FetchGrid(t.db, cal.SelectedDate, t.plannerModel.Labels.Selected)
// 	t.err = err
//
// 	t.plannerModel.Entries = entries
// 	cal.UpdateGrid(t.plannerModel.HasEntry)
// }

// // reloadGrid is loadEntries plus showing the selected day's entry in Notes.
// func (t *Todo) reloadGrid() {
// 	t.loadEntries()
// 	t.plannerModel.SyncNotes()
// }
//
// // saveEntry keeps Notes as it is, so the cursor doesn't jump while typing.
// func (t *Todo) saveEntry(m planner.SaveEntryMsg) {
// 	err := database.AddEntry(t.db, m.Date.FormatSQL(), string(m.Label), m.Msg)
// 	if t.err = err; err != nil {
// 		return
// 	}
// 	t.loadEntries()
// }
//
// func (t *Todo) removeEntry(m planner.RemoveEntryMsg) {
// 	err := database.RemoveEntry(t.db, m.Date.FormatSQL(), string(m.Label))
// 	if t.err = err; err != nil {
// 		return
// 	}
// 	t.reloadGrid()
// }
