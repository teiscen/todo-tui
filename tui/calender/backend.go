package calender

import (
	"time"
)

type (
	Label  int
	Status int
	Entry  struct {
		Label  Label
		Status Status
		Msg    string
	}
)

type (
	Year         int
	Month        int
	Date         int
	CalenderDate struct {
		Year  Year
		Month Month
		Date  Date
	}
)

func (c CalenderDate) lastDay() CalenderDate {
	y := int(c.Year)
	m := time.Month(c.Month + 1)
	d := 0
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return timeToCalenderDate(inTime)
}

func (c CalenderDate) firstDay() CalenderDate {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := 1
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return timeToCalenderDate(inTime)
}

func (c CalenderDate) addDate(y int, m int, d int) CalenderDate {
	t := calenderDateToTime(c)
	t.AddDate(y, m, d)
	return timeToCalenderDate(t)
}

func timeToCalenderDate(t time.Time) CalenderDate {
	y, mTemp, d := t.Date()
	m := int(mTemp)

	return CalenderDate{
		Year:  Year(y),
		Month: Month(m),
		Date:  Date(d),
	}
}

func calenderDateToTime(c CalenderDate) time.Time {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := int(c.Date)
	locn := time.Now().Location()
	return time.Date(y, m, d, 0, 0, 0, 0, locn)
}

type Calender struct {
	Months map[struct {
		Year  Year
		Month Month
	}]map[Date][]Entry
}

func (c Calender) getEntryList(cDate CalenderDate) (entryList []Entry, ok bool) {
	key := struct {
		Year  Year
		Month Month
	}{
		Year:  cDate.Year,
		Month: cDate.Month,
	}

	entryList, ok = c.Months[key][cDate.Date]
	return
}

func (c Calender) getEntry(entryList []Entry, label Label) (entry Entry) {
	for i := range entryList {
		if entryList[i].Label == label {
			return entryList[i]
		}
	}
	return
}
