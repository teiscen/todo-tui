package backend

import (
	"fmt"
	"time"
)

type (
	Year  int
	Month int
	Day   int
	Date  struct {
		Year  Year
		Month Month
		Date  Day
	}
)

func (c Date) Weekday() time.Weekday {
	tDate := calendarDateToTime(c)
	return tDate.Weekday()
}

func (c Date) NumDays() int {
	numDays := int(c.LastDay().Date)
	return numDays
}

func (c Date) LastDay() Date {
	y := int(c.Year)
	m := time.Month(c.Month + 1)
	d := 0
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return TimeToDate(inTime)
}

func (c Date) FirstDay() Date {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := 1
	locn := time.Now().Location()
	inTime := time.Date(y, m, d, 0, 0, 0, 0, locn)
	return TimeToDate(inTime)
}

func (c Date) AddDate(y int, m int, d int) Date {
	t := calendarDateToTime(c)
	t = t.AddDate(y, m, d)
	return TimeToDate(t)
}

func TimeToDate(t time.Time) Date {
	y, mTemp, d := t.Date()
	m := int(mTemp)

	return Date{
		Year:  Year(y),
		Month: Month(m),
		Date:  Day(d),
	}
}

func calendarDateToTime(c Date) time.Time {
	y := int(c.Year)
	m := time.Month(c.Month)
	d := int(c.Date)
	locn := time.Now().Location()
	return time.Date(y, m, d, 0, 0, 0, 0, locn)
}

func (d Date) Format() string {
	return calendarDateToTime(d).Format("January 2")
}

// If Date is the struct {Year, Month, Date int} from earlier:
func (d Date) FormatSQL() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Date)
}
