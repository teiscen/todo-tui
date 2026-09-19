package calendar

import (
	"time"
)

type (
	Status int
	Entry  struct {
		Status Status
		Msg    string
	}
)

type (
	Year         int
	Month        int
	Date         int
	CalendarDate struct {
		Year  Year
		Month Month
		Date  Date
	}
)

func (c CalendarDate) weekday() time.Weekday {
	tDate := calendarDateToTime(c)
	return tDate.Weekday()
}

func (c CalendarDate) numDays() int {
	numDays := int(c.lastDay().Date)
	return numDays
}

func (c CalendarDate) lastDay() CalendarDate {
	y := int(c.Year)
	m := time.Month(c.Month + 1)
	d := 0
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return timeToCalendarDate(inTime)
}

func (c CalendarDate) firstDay() CalendarDate {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := 1
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return timeToCalendarDate(inTime)
}

func (c CalendarDate) addDate(y int, m int, d int) CalendarDate {
	t := calendarDateToTime(c)
	t = t.AddDate(y, m, d)
	return timeToCalendarDate(t)
}

func timeToCalendarDate(t time.Time) CalendarDate {
	y, mTemp, d := t.Date()
	m := int(mTemp)

	return CalendarDate{
		Year:  Year(y),
		Month: Month(m),
		Date:  Date(d),
	}
}

func calendarDateToTime(c CalendarDate) time.Time {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := int(c.Date)
	locn := time.Now().Location()
	return time.Date(y, m, d, 0, 0, 0, 0, locn)
}

type Calendar struct {
	Months map[CalendarDate]map[Label]Entry
}

func (c Calendar) getEntry(d CalendarDate, l Label) (val Entry, ok bool) {
	val, ok = c.Months[d][l]
	return
}
