package calendar

import (
	"math/rand/v2"
	"time"
)

type MockRange struct {
	monthOffset     int
	entryCountRange struct {
		lower int
		upper int
	}
}

var r = MockRange{
	monthOffset: 2,
	entryCountRange: struct {
		lower int
		upper int
	}{
		lower: 5,
		upper: 15,
	},
}

func mockCalendar(labelInfo map[Label]LabelInfo) Calendar {
	randUniqueNums := func(min, max, count int) []int {
		available := make([]int, max-min+1)

		for i := range available {
			available[i] = min + i
		}

		rand.Shuffle(len(available), func(i, j int) {
			available[i], available[j] = available[j], available[i]
		})

		return available[:count]
	}

	now := timeToCalendarDate(time.Now())
	entries := make(map[CalendarDate]map[Label]Entry)

	for monthOffset := -r.monthOffset; monthOffset <= r.monthOffset; monthOffset++ {
		month := now.addDate(0, monthOffset, 0)

		for label := range labelInfo {
			numEntries := rand.IntN(
				r.entryCountRange.upper-r.entryCountRange.lower+1,
			) + r.entryCountRange.lower

			days := randUniqueNums(1, month.numDays(), numEntries)

			for _, day := range days {
				date := CalendarDate{
					Year:  month.Year,
					Month: month.Month,
					Date:  Date(day),
				}

				if entries[date] == nil {
					entries[date] = make(map[Label]Entry)
				}

				entries[date][label] = Entry{
					Status: 1,
					Msg:    "mock entry",
				}
			}
		}
	}

	return Calendar{Months: entries}
}
