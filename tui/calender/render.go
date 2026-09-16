package calender

type cell struct {
	renderF func(ThemeColor) string
	cDate   CalenderDate
}
type Grid []cell

func CreateGrid(c Calender, k MonthKey) (grid Grid) {
	// cDate := CalenderDate{k.Year, k.Month, 1}

	// func Foo(startDate CalenderDate, itr int) grid

	// lastWeekday := cDate.lastDay().weekday()
	return Grid{}
}
