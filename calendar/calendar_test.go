package calendar

import (
	"testing"
	"time"

	backend "todo-tui/backend"
)

func d(y int, m time.Month, day int) backend.Date {
	return backend.NewDate(time.Date(y, m, day, 0, 0, 0, 0, time.UTC))
}

// setGrid fills the 6x7 state with consecutive days starting at start.
func setGrid(c *Calendar, start backend.Date) {
	for row := range c.State {
		for col := range c.State[row] {
			c.State[row][col] = Info{Date: start}
			start = start.AddDate(0, 0, 1)
		}
	}
}

func TestMovement(t *testing.T) {
	c := New(backend.DefaultTheme(), d(2026, time.January, 15))

	// Right stays in the same month.
	if cmd := c.Movement(Right); c.SelectedDate != d(2026, time.January, 16) {
		t.Errorf("Right selected = %v, want Jan 16", c.SelectedDate)
	} else if _, ok := cmd().(ChangeDateMsg); !ok {
		t.Errorf("Right cmd = %T, want ChangeDateMsg", cmd())
	}

	// Down moves a week.
	c.SelectedDate = d(2026, time.January, 15)
	c.Movement(Down)
	if c.SelectedDate != d(2026, time.January, 22) {
		t.Errorf("Down selected = %v, want Jan 22", c.SelectedDate)
	}

	// Left across a month boundary emits ChangeMonthMsg.
	c.SelectedDate = d(2026, time.January, 1)
	cmd := c.Movement(Left)
	if c.SelectedDate != d(2025, time.December, 31) {
		t.Errorf("Left selected = %v, want Dec 31 2025", c.SelectedDate)
	}
	if msg, ok := cmd().(ChangeMonthMsg); !ok {
		t.Errorf("Left cmd = %T, want ChangeMonthMsg", cmd())
	} else if msg.NewDate != c.SelectedDate {
		t.Errorf("ChangeMonthMsg.NewDate = %v, want %v", msg.NewDate, c.SelectedDate)
	}

	// Up moves back a week.
	c.SelectedDate = d(2026, time.January, 8)
	c.Movement(Up)
	if c.SelectedDate != d(2026, time.January, 1) {
		t.Errorf("Up selected = %v, want Jan 1", c.SelectedDate)
	}
}

func TestGetSelectedColor(t *testing.T) {
	theme := backend.DefaultTheme()
	c := New(theme, d(2026, time.January, 15))
	setGrid(&c, d(2026, time.January, 1))

	// Jan 15 2026 is a Thursday (weekday).
	c.SelectedDate = d(2026, time.January, 15)
	if got := c.GetSelectedColor(); got != theme.Weekday {
		t.Errorf("weekday color = %q, want %q", got, theme.Weekday)
	}

	// Jan 3 2026 is a Saturday (weekend).
	c.SelectedDate = d(2026, time.January, 3)
	if got := c.GetSelectedColor(); got != theme.Weekend {
		t.Errorf("weekend color = %q, want %q", got, theme.Weekend)
	}

	// Jan 4 2026 is a Sunday (weekend).
	c.SelectedDate = d(2026, time.January, 4)
	if got := c.GetSelectedColor(); got != theme.Weekend {
		t.Errorf("sunday color = %q, want %q", got, theme.Weekend)
	}

	// A marked day uses the accent color regardless of weekday/weekend.
	c.SelectedDate = d(2026, time.January, 3)
	if info := c.GetInfo(c.SelectedDate); info != nil {
		info.IsMarked = true
	}
	if got := c.GetSelectedColor(); got != theme.Accent {
		t.Errorf("marked color = %q, want %q", got, theme.Accent)
	}
}

func TestToggleMark(t *testing.T) {
	c := New(backend.DefaultTheme(), d(2026, time.January, 15))
	setGrid(&c, d(2026, time.January, 1))

	info := c.GetInfo(c.SelectedDate)
	if info == nil {
		t.Fatal("selected date not found in grid")
	}
	if info.IsMarked {
		t.Fatal("selected date should start unmarked")
	}

	if cmd := c.ToggleMark(); !info.IsMarked {
		t.Error("selected date should be marked after first toggle")
	} else if msg, ok := cmd().(ToggleMarkMsg); !ok || !msg.Status {
		t.Errorf("ToggleMark cmd = %T %+v, want Status=true", msg, msg)
	}

	c.ToggleMark()
	if info.IsMarked {
		t.Error("selected date should be unmarked after second toggle")
	}
}
