package calender

import (
	"math/rand/v2"
	"time"
)

type randomizerParameter struct {
	monthOffset     int
	entryCountRange struct {
		lower int
		upper int
	}
}

var r = randomizerParameter{
	monthOffset: 2,
	entryCountRange: struct {
		lower int
		upper int
	}{
		lower: 5,
		upper: 15,
	},
}

func mockCalender() Calender {
	// min and max are inclusive
	randUniqueNums := func(min int, max int, count int) []int {
		results := make([]int, count)
		for i := range results {
			results[i] = rand.IntN(max+1) + min
		}
		return results
	}

	c := timeToCalenderDate(time.Now())

	// Loop through all the MONTHS
	out := make(map[CalenderDate]map[Label]Entry)
	for x := -1 * r.monthOffset; x < r.monthOffset; x++ {
		m := c.addDate(0, x, 0)

		// For each LABEL
		for y := range labelInfo {

			// generate which Days of entries
			numEntries := rand.IntN(r.entryCountRange.upper-r.entryCountRange.lower+1) + r.entryCountRange.lower
			days := randUniqueNums(1, m.numDays(), numEntries)

			// For each of those DAYS
			for _, d := range days {
				cDate := CalenderDate{Year: m.Year, Month: m.Month, Date: Date(d)}
				if out[cDate] == nil {
					out[cDate] = make(map[Label]Entry)
				}
				out[cDate][y] = Entry{
					Status: 1,
					Msg:    "mock entry",
				}
			}
		}
	}

	return Calender{out}
}
