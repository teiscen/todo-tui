package todo

import (
	"todo-tui/backend"
	"todo-tui/database"
)

func (t *Todo) ReloadPlannerState(lID backend.LabelID, date backend.Date) {
	// Start Date is the first Sunday from the first of the month.
	// End Date is 41 days after it.
	sDate := date
	sDate.FirstDay()
	sDate.ChangeDate(0, 0, -sDate.Weekday())
	eDate := sDate.AddDate(0, 0, 41)

	entries, err := database.FetchGridEntries(t.db, sDate, eDate, lID)
	if err != nil {
		t.err = err
		return
	}
	t.err = nil
	t.plannerModel.ReloadState(sDate, entries)
	t.previewModel.SetEntries(entries)
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
