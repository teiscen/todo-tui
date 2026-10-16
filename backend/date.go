package backend

import (
	"time"
)

// Date is a calendar day with no time of day. It is always stored as midnight
// UTC, so two Dates for the same day compare equal with == as well as Equal.
// Every method returns a new Date; none of them modify the receiver.
type Date struct {
	t time.Time
}

func NewDate(t time.Time) Date {
	return Date{t: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

func (d Date) Year() int {
	return d.t.Year()
}

func (d Date) Month() int {
	return int(d.t.Month())
}

func (d Date) Day() int {
	return d.t.Day()
}

// Weekday returns the day of the week (Sunday = 0, Saturday = 6).
func (d Date) Weekday() int {
	return int(d.t.Weekday())
}

func (d Date) IsWeekend() bool {
	return d.Weekday() == 0 || d.Weekday() == 6
}

// FirstDay returns the 1st of the month containing d.
func (d Date) FirstDay() Date {
	return Date{t: time.Date(d.Year(), time.Month(d.Month()), 1, 0, 0, 0, 0, time.UTC)}
}

// LastDay returns the last day of the month containing d.
func (d Date) LastDay() Date {
	return Date{t: time.Date(d.Year(), time.Month(d.Month()+1), 0, 0, 0, 0, 0, time.UTC)}
}

// NumDays is the number of days in d's month.
func (d Date) NumDays() int {
	return d.LastDay().Day()
}

// AddDate follows time.AddDate: a day that doesn't exist rolls into the next month.
func (d Date) AddDate(y, m, days int) Date {
	return Date{t: d.t.AddDate(y, m, days)}
}

// AddMonths moves by whole months and clamps the day instead of rolling over,
// so Jan 31 + 1 month is Feb 28, not Mar 3.
func (d Date) AddMonths(n int) Date {
	first := d.FirstDay().AddDate(0, n, 0)
	return first.AddDate(0, 0, min(d.Day(), first.NumDays())-1)
}

func (d Date) Format() string {
	return d.t.Format("January 2")
}

// FormatSQL returns a zero-padded YYYY-MM-DD string, so text comparisons
// (BETWEEN, ORDER BY) in SQLite sort chronologically.
func (d Date) FormatSQL() string {
	return d.t.Format("2006-01-02")
}

// Key returns a timezone-independent key for map lookups (YYYY-MM-DD).
func (d Date) Key() string {
	return d.t.Format("2006-01-02")
}

// Equal reports whether d and other represent the same calendar date.
func (d Date) Equal(other Date) bool {
	return d.Compare(other) == 0
}

// Compare returns -1 if d < other, 0 if equal, 1 if d > other (date-only).
func (d Date) Compare(other Date) int {
	if d.Year() != other.Year() {
		if d.Year() < other.Year() {
			return -1
		}
		return 1
	}
	if d.Month() != other.Month() {
		if d.Month() < other.Month() {
			return -1
		}
		return 1
	}
	if d.Day() != other.Day() {
		if d.Day() < other.Day() {
			return -1
		}
		return 1
	}
	return 0
}

// Before reports whether d is before other (date-only).
func (d Date) Before(other Date) bool {
	return d.Compare(other) < 0
}

// After reports whether d is after other (date-only).
func (d Date) After(other Date) bool {
	return d.Compare(other) > 0
}

// DaysUntil returns the number of days from d to other (positive if other is in the future).
func (d Date) DaysUntil(other Date) int {
	dMid := time.Date(d.Year(), time.Month(d.Month()), d.Day(), 0, 0, 0, 0, time.UTC)
	oMid := time.Date(other.Year(), time.Month(other.Month()), other.Day(), 0, 0, 0, 0, time.UTC)
	return int(oMid.Sub(dMid).Hours() / 24)
}

func ParseSQL(str string) (Date, error) {
	t, err := time.Parse("2006-01-02", str)
	return Date{t: t}, err
}
