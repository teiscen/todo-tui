package backend

import (
	"testing"
	"time"
)

func date(y int, m time.Month, d int) Date {
	return NewDate(time.Date(y, m, d, 0, 0, 0, 0, time.UTC))
}

func TestDateAccessors(t *testing.T) {
	d := date(2024, time.January, 1) // a Monday
	if d.Year() != 2024 {
		t.Errorf("Year() = %d, want 2024", d.Year())
	}
	if d.Month() != 1 {
		t.Errorf("Month() = %d, want 1", d.Month())
	}
	if d.Day() != 1 {
		t.Errorf("Day() = %d, want 1", d.Day())
	}
	if d.Weekday() != int(time.Monday) {
		t.Errorf("Weekday() = %d, want %d", d.Weekday(), int(time.Monday))
	}
}

func TestIsWeekend(t *testing.T) {
	cases := []struct {
		d    Date
		want bool
	}{
		{date(2024, time.January, 6), true},  // Saturday
		{date(2024, time.January, 7), true},  // Sunday
		{date(2024, time.January, 1), false}, // Monday
		{date(2024, time.January, 5), false}, // Friday
	}
	for _, c := range cases {
		if got := c.d.IsWeekend(); got != c.want {
			t.Errorf("%s IsWeekend() = %v, want %v", c.d.FormatSQL(), got, c.want)
		}
	}
}

func TestNumDays(t *testing.T) {
	cases := []struct {
		d    Date
		want int
	}{
		{date(2024, time.February, 15), 29}, // leap year
		{date(2023, time.February, 15), 28},
		{date(2024, time.January, 15), 31},
		{date(2024, time.April, 15), 30},
		{date(2024, time.December, 15), 31},
	}
	for _, c := range cases {
		if got := c.d.NumDays(); got != c.want {
			t.Errorf("%s NumDays() = %d, want %d", c.d.FormatSQL(), got, c.want)
		}
	}
}

func TestFormat(t *testing.T) {
	if got := date(2024, time.January, 2).Format(); got != "January 2" {
		t.Errorf("Format() = %q, want %q", got, "January 2")
	}
}

func TestFormatSQLAndKey(t *testing.T) {
	d := date(2024, time.January, 2)
	if got := d.FormatSQL(); got != "2024-01-02" {
		t.Errorf("FormatSQL() = %q, want %q", got, "2024-01-02")
	}
	if got := d.Key(); got != "2024-01-02" {
		t.Errorf("Key() = %q, want %q", got, "2024-01-02")
	}
}

func TestEqual(t *testing.T) {
	a := date(2024, time.January, 2)
	if !a.Equal(date(2024, time.January, 2)) {
		t.Error("Equal() = false for identical dates")
	}
	if a.Equal(date(2024, time.January, 3)) {
		t.Error("Equal() = true for different days")
	}
	if a.Equal(date(2024, time.February, 2)) {
		t.Error("Equal() = true for different months")
	}
}

func TestCompare(t *testing.T) {
	ref := date(2024, time.June, 15)
	cases := []struct {
		other Date
		want  int
	}{
		{date(2024, time.June, 14), 1}, // after
		{date(2024, time.June, 15), 0}, // equal
		{date(2024, time.June, 16), -1},
		{date(2024, time.May, 15), 1},
		{date(2024, time.July, 15), -1},
		{date(2023, time.June, 15), 1},
		{date(2025, time.June, 15), -1},
	}
	for _, c := range cases {
		if got := ref.Compare(c.other); got != c.want {
			t.Errorf("Compare(%s) = %d, want %d", c.other.FormatSQL(), got, c.want)
		}
	}
}

func TestBeforeAfter(t *testing.T) {
	a := date(2024, time.June, 15)
	b := date(2024, time.June, 16)
	if !a.Before(b) {
		t.Error("Before() = false for earlier date")
	}
	if a.After(b) {
		t.Error("After() = true for earlier date")
	}
	if b.Before(a) || !b.After(a) {
		t.Error("Before/After inverted for later date")
	}
}

func TestDaysUntil(t *testing.T) {
	a := date(2024, time.January, 1)
	cases := []struct {
		other Date
		want  int
	}{
		{date(2024, time.January, 10), 9},
		{date(2024, time.January, 1), 0},
		{date(2023, time.December, 25), -7},
	}
	for _, c := range cases {
		if got := a.DaysUntil(c.other); got != c.want {
			t.Errorf("DaysUntil(%s) = %d, want %d", c.other.FormatSQL(), got, c.want)
		}
	}
}

func TestAddDateValueSemantics(t *testing.T) {
	orig := date(2024, time.January, 31)
	next := orig.AddDate(0, 0, 1)
	prev := orig.AddDate(0, 0, -1)

	if !next.Equal(date(2024, time.February, 1)) {
		t.Errorf("AddDate(+1 day) = %s, want 2024-02-01", next.FormatSQL())
	}
	if !prev.Equal(date(2024, time.January, 30)) {
		t.Errorf("AddDate(-1 day) = %s, want 2024-01-30", prev.FormatSQL())
	}
	// The receiver must not be mutated by a value-method AddDate.
	if !orig.Equal(date(2024, time.January, 31)) {
		t.Errorf("AddDate mutated receiver: got %s", orig.FormatSQL())
	}
}

func TestFirstDayLastDay(t *testing.T) {
	cases := []struct {
		d        Date
		first    Date
		last     Date
		numDays  int
	}{
		{date(2024, time.February, 15), date(2024, time.February, 1), date(2024, time.February, 29), 29},
		{date(2023, time.February, 15), date(2023, time.February, 1), date(2023, time.February, 28), 28},
		{date(2024, time.January, 31), date(2024, time.January, 1), date(2024, time.January, 31), 31},
		{date(2024, time.April, 30), date(2024, time.April, 1), date(2024, time.April, 30), 30},
	}
	for _, c := range cases {
		orig := c.d
		if got := c.d.FirstDay(); !got.Equal(c.first) {
			t.Errorf("%s FirstDay() = %s, want %s", c.d.FormatSQL(), got.FormatSQL(), c.first.FormatSQL())
		}
		if got := c.d.LastDay(); !got.Equal(c.last) {
			t.Errorf("%s LastDay() = %s, want %s", c.d.FormatSQL(), got.FormatSQL(), c.last.FormatSQL())
		}
		if got := c.d.NumDays(); got != c.numDays {
			t.Errorf("%s NumDays() = %d, want %d", c.d.FormatSQL(), got, c.numDays)
		}
		// FirstDay/LastDay must not mutate the receiver.
		if !c.d.Equal(orig) {
			t.Errorf("FirstDay/LastDay mutated receiver: got %s, want %s", c.d.FormatSQL(), orig.FormatSQL())
		}
	}
}

func TestAddMonths(t *testing.T) {
	cases := []struct {
		d    Date
		n    int
		want Date
	}{
		// Day clamping across month lengths.
		{date(2024, time.January, 31), 1, date(2024, time.February, 29)},  // leap year clamp
		{date(2023, time.January, 31), 1, date(2023, time.February, 28)},
		{date(2024, time.March, 31), -1, date(2024, time.February, 29)},   // clamp backwards
		{date(2024, time.May, 31), 1, date(2024, time.June, 30)},
		// Day preserved when valid.
		{date(2024, time.January, 15), 1, date(2024, time.February, 15)},
		{date(2024, time.January, 31), 2, date(2024, time.March, 31)},     // 31 is valid in March
		// Year boundary.
		{date(2024, time.December, 15), 1, date(2025, time.January, 15)},
		{date(2024, time.January, 15), -1, date(2023, time.December, 15)},
		// Zero move is identity.
		{date(2024, time.January, 31), 0, date(2024, time.January, 31)},
	}
	for _, c := range cases {
		orig := c.d
		got := c.d.AddMonths(c.n)
		if !got.Equal(c.want) {
			t.Errorf("%s AddMonths(%d) = %s, want %s", c.d.FormatSQL(), c.n, got.FormatSQL(), c.want.FormatSQL())
		}
		// AddMonths must not mutate the receiver.
		if !c.d.Equal(orig) {
			t.Errorf("AddMonths(%d) mutated receiver: got %s, want %s", c.n, c.d.FormatSQL(), orig.FormatSQL())
		}
	}
}

func TestParseSQL(t *testing.T) {
	d, err := ParseSQL("2024-01-02")
	if err != nil {
		t.Fatalf("ParseSQL returned error: %v", err)
	}
	if !d.Equal(date(2024, time.January, 2)) {
		t.Errorf("ParseSQL = %s, want 2024-01-02", d.FormatSQL())
	}

	if _, err := ParseSQL("not-a-date"); err == nil {
		t.Error("ParseSQL did not return an error for invalid input")
	}
}
