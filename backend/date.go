package backend

import (
	"time"
)

type Date struct {
	t time.Time
}

func (d Date) GetTime() time.Time {
	return d.t
}

func NewDate(t time.Time) Date {
	return Date{t: t}
}

func (d Date) IsWeekend() bool {
	if d.Weekday() == 0 || d.Weekday() == 6 {
		return true
	}
	return false
}

func (c Date) Year() int {
	return c.t.Year()
}

func (c Date) Month() int {
	return int(c.t.Month())
}

func (c Date) Day() int {
	return c.t.Day()
}

func (c Date) Location() *time.Location {
	return c.t.Location()
}

func (c Date) Weekday() int {
	return int(c.t.Weekday())
}

func (c Date) NumDays() int {
	c.LastDay()
	return c.t.Day()
}

func (c *Date) LastDay() {
	c.t = time.Date(
		c.Year(), time.Month(c.Month()+1), 0,
		0, 0, 0, 0, c.Location(),
	)
}

func (c *Date) FirstDay() {
	c.t = time.Date(
		c.Year(), time.Month(c.Month()), 1,
		0, 0, 0, 0, c.Location(),
	)
}

func (c *Date) ChangeDate(y, m, d int) {
	c.t = c.t.AddDate(y, m, d)
}

func (c Date) AddDate(y, m, d int) Date {
	c.ChangeDate(y, m, d)
	return c
}

func (c Date) Format() string {
	return c.t.Format("January 2")
}

// FormatSQL returns a zero-padded YYYY-MM-DD string, so text comparisons
// (BETWEEN, ORDER BY) in SQLite sort chronologically.
func (c Date) FormatSQL() string {
	return c.t.Format("2006-01-02")
}

// Key returns a timezone-independent key for map lookups (YYYY-MM-DD).
func (c Date) Key() string {
	return c.t.Format("2006-01-02")
}

// Equal reports whether c and other represent the same calendar date (ignores timezone).
func (c Date) Equal(other Date) bool {
	return c.Year() == other.Year() && c.Month() == other.Month() && c.Day() == other.Day()
}

// Compare returns -1 if c < other, 0 if equal, 1 if c > other (date-only).
func (c Date) Compare(other Date) int {
	if c.Year() != other.Year() {
		if c.Year() < other.Year() {
			return -1
		}
		return 1
	}
	if c.Month() != other.Month() {
		if c.Month() < other.Month() {
			return -1
		}
		return 1
	}
	if c.Day() != other.Day() {
		if c.Day() < other.Day() {
			return -1
		}
		return 1
	}
	return 0
}

// Before reports whether c is before other (date-only).
func (c Date) Before(other Date) bool {
	return c.Compare(other) < 0
}

// After reports whether c is after other (date-only).
func (c Date) After(other Date) bool {
	return c.Compare(other) > 0
}

// DaysUntil returns the number of days from c to other (positive if other is in the future).
func (c Date) DaysUntil(other Date) int {
	cMid := time.Date(c.Year(), time.Month(c.Month()), c.Day(), 0, 0, 0, 0, time.UTC)
	oMid := time.Date(other.Year(), time.Month(other.Month()), other.Day(), 0, 0, 0, 0, time.UTC)
	return int(oMid.Sub(cMid).Hours() / 24)
}

func ParseSQL(str string) (Date, error) {
	t, err := time.Parse("2006-01-02", str)
	return Date{t: t}, err
}
