package leftbox

import "todo-tui/backend"

// MockRows is stand-in data to look at the view before the database is wired in.
// Delete once the loader exists.
func MockRows() []Row {
	col := map[backend.LabelID]backend.HexCode{
		"uni": "#81c8be", "work": "#ca9ee6", "gym": "#ef9f76", "app": "#babbf1", "vol": "#f2d5cf",
	}
	r := func(id int, tab Tab, g Group, title string, lab backend.LabelID, done bool, sd, st int, due string, tone Tone) Row {
		return Row{
			ID: id, Tab: tab, Group: g, Title: title, LabelID: lab, Color: col[lab],
			Done: done, StepsDone: sd, StepsTotal: st, Due: due, DueTone: tone,
		}
	}
	rec := func(id int, g Group, title string, lab backend.LabelID, done bool, sched, due string) Row {
		return Row{
			ID: id, Tab: TabRecurring, Group: g, Title: title, LabelID: lab, Color: col[lab],
			Done: done, Schedule: sched, Due: due,
		}
	}
	return []Row{
		r(1, TabSmall, GroupOverdue, "Write lab report", "uni", false, 3, 5, "2d late", ToneLate),
		r(3, TabSmall, GroupToday, "Email supervisor", "work", false, 0, 0, "today", ToneSoon),
		r(9, TabSmall, GroupNoDate, "Laundry", "", false, 0, 0, "", ToneNormal),
		r(4, TabSmall, GroupWeek, "Order gym shoes", "gym", false, 0, 0, "in 6d", ToneNormal),
		r(5, TabSmall, GroupLater, "Book dentist", "", false, 0, 0, "in 12d", ToneNormal),
		r(6, TabSmall, GroupLater, "Pay tuition", "uni", false, 0, 0, "in 20d", ToneNormal),
		r(10, TabSmall, GroupDone, "Renew bus pass", "", true, 0, 0, "", ToneNormal),
		r(11, TabSmall, GroupDone, "Return library books", "", true, 0, 0, "", ToneNormal),

		r(2, TabLong, GroupWeek, "Apply to co-op", "app", false, 1, 5, "in 3d", ToneNormal),
		r(7, TabLong, GroupLater, "Midterm prep", "uni", false, 2, 6, "in 8d", ToneNormal),
		r(8, TabLong, GroupLater, "Plan volunteer shifts", "vol", false, 1, 4, "in 21d", ToneNormal),

		rec(12, GroupToday, "Gym", "gym", false, "M/W/F", "today"),
		rec(13, GroupToday, "Read 20 minutes", "", false, "daily", "today"),
		rec(15, GroupLater, "Water plants", "", false, "every 3d", "in 2d"),
		rec(14, GroupLater, "Review notes", "uni", false, "Su", "in 5d"),
	}
}
