package calender

type cell struct {
	renderF func(ThemeColor) string
	cDate   CalenderDate
}
type Grid []cell

func CreateGrid(c Calender, k MonthKey) (grid Grid) {
	// cDate := CalenderDate{k.Year, k.Month, 1}

	// func Foo(startDate CalenderDate, endDate CalenderDate) grid
	// newDate := startDate
	// for i := startDate.day; i <=/ endDate.day; i++ {
	//
	// 		newDate.addDate(0 1 0)
	// }

	// lastWeekday := cDate.lastDay().weekday()
	return Grid{}
}
