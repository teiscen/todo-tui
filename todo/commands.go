package todo

import (
	"sort"

	"todo-tui/backend"
	"todo-tui/database"
	"todo-tui/weekly"
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

func (t *Todo) LoadWeeklyData(refDate backend.Date) {
	today := refDate
	weekLater := today.AddDate(0, 0, 6)  // 7 days total (including today)
	pastDays := today.AddDate(0, 0, -99) // 99 days ago for overdue

	// Fetch overdue entries (99 days ago up to today)
	overdueEntries, err := database.FetchAllEntriesInRange(t.db, pastDays, today)
	if err != nil {
		overdueEntries = make(map[string]map[string]string)
	}

	// Fetch upcoming entries (today to 7 days later)
	upcomingEntries, err := database.FetchAllEntriesInRange(t.db, today, weekLater)
	if err != nil {
		upcomingEntries = make(map[string]map[string]string)
	}

	// Sort overdue dates (newest first for overdue list)
	overdueDates := make([]string, 0, len(overdueEntries))
	for d := range overdueEntries {
		overdueDates = append(overdueDates, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(overdueDates))) // newest first

	// Sort upcoming dates (oldest first for chronological display)
	upcomingDates := make([]string, 0, len(upcomingEntries))
	for d := range upcomingEntries {
		upcomingDates = append(upcomingDates, d)
	}
	sort.Strings(upcomingDates)

	var tasks []weekly.Day
	var overdue []weekly.Task

	// Process overdue (newest first)
	for _, dateKey := range overdueDates {
		byLabel := overdueEntries[dateKey]
		date, _ := backend.ParseSQL(dateKey)

		for labelID, msg := range byLabel {
			overdue = append(overdue, weekly.Task{
				Date:    date,
				Text:    msg,
				LabelID: backend.LabelID(labelID),
			})
		}
	}

	// Process upcoming (chronological)
	for _, dateKey := range upcomingDates {
		byLabel := upcomingEntries[dateKey]
		date, _ := backend.ParseSQL(dateKey)

		var dayTasks []weekly.Task
		for labelID, msg := range byLabel {
			dayTasks = append(dayTasks, weekly.Task{
				Date:    date,
				Text:    msg,
				LabelID: backend.LabelID(labelID),
			})
		}

		if len(dayTasks) > 0 {
			tasks = append(tasks, weekly.Day{
				Date:  date,
				Tasks: dayTasks,
			})
		}
	}

	t.weeklyModel.Tasks = tasks
	t.weeklyModel.Overdue = overdue
	t.weeklyModel.TopLine = 5 // Start scrolled down to show overdue
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
