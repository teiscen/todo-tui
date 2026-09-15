package calender

import (
	"math/rand/v2"
	"time"
)

type randomizerParameter struct {
	monthOffset     int
	labelCount      int
	entryCountRange struct {
		lower int
		upper int
	}
}

var r = randomizerParameter{
	monthOffset: 2,
	labelCount:  5,
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

	type key = struct {
		Year  Year
		Month Month
	}

	loopMonths := func(c CalenderDate,
		innerF func(c CalenderDate) map[Date][]Entry,
	) map[key]map[Date][]Entry {
		out := make(map[key]map[Date][]Entry)
		for i := -1 * r.monthOffset; i < r.monthOffset; i++ {
			m := c.addDate(0, i, 0)
			k := key{Year: m.Year, Month: m.Month}
			out[k] = innerF(m)
		}
		return out
	}

	loopDays := func(c CalenderDate) map[Date][]Entry {
		out := make(map[Date][]Entry)
		for label := 0; label < r.labelCount; label++ {
			upperRange := r.entryCountRange.upper
			lowerRange := r.entryCountRange.lower
			count := rand.IntN(upperRange-lowerRange+1) + lowerRange

			numDays := int(c.lastDay().Date)
			days := randUniqueNums(1, numDays, count)
			for _, d := range days {
				out[Date(d)] = append(out[Date(d)], Entry{
					Label:  Label(label),
					Status: 1,
					Msg:    "mock entry",
				})
			}
		}
		return out
	}

	cDate := timeToCalenderDate(time.Now())
	return Calender{Months: loopMonths(cDate, loopDays)}
}
