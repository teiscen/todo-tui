package backend

//
// import (
// 	"math/rand/v2"
// 	"time"
// )
//
// type MockRange struct {
// 	monthOffset     int
// 	entryCountRange struct {
// 		lower int
// 		upper int
// 	}
// }
//
// var r = MockRange{
// 	monthOffset: 2,
// 	entryCountRange: struct {
// 		lower int
// 		upper int
// 	}{
// 		lower: 5,
// 		upper: 15,
// 	},
// }
//
// func randUniqueNums(min int, max int, count int) []int {
// 	available := make([]int, max-min+1)
//
// 	for i := range available {
// 		available[i] = min + i
// 	}
//
// 	rand.Shuffle(len(available), func(i, j int) {
// 		available[i], available[j] = available[j], available[i]
// 	})
//
// 	return available[:count]
// }
//
// // Changes mock Calendar to return labels alongside Calendar
// func MockCalendar() Calendar {
// 	now := TimeToDate(time.Now())
// 	entries := make(map[Date]map[LabelID]Entry)
//
// 	var labels Labels
// 	for monthOffset := -r.monthOffset; monthOffset <= r.monthOffset; monthOffset++ {
// 		month := now.AddDate(0, monthOffset, 0)
//
// 		labels, _ = MockConfig()
//
// 		for _, label := range labels.Order {
// 			numEntries := rand.IntN(
// 				r.entryCountRange.upper-r.entryCountRange.lower+1,
// 			) + r.entryCountRange.lower
//
// 			days := randUniqueNums(1, month.NumDays(), numEntries)
//
// 			for _, day := range days {
// 				date := Date{
// 					Year:  month.Year,
// 					Month: month.Month,
// 					Date:  Day(day),
// 				}
//
// 				if entries[date] == nil {
// 					entries[date] = make(map[LabelID]Entry)
// 				}
//
// 				entries[date][label] = Entry{
// 					Status: 1,
// 					Msg:    "mock entry",
// 				}
// 			}
// 		}
// 	}
// 	return Calendar{Months: entries, Labels: labels}
// }
