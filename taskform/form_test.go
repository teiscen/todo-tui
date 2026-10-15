package taskform

import (
	"testing"

	"todo-tui/backend"
)

func testLabels() backend.Labels {
	return backend.Labels{
		Selected: "work",
		Order:    []backend.LabelID{"uni", "work", "gym"},
		Info: map[backend.LabelID]backend.LabelInfo{
			"uni":  {Name: "University", Color: "#81c8be"},
			"work": {Name: "Work", Color: "#ca9ee6"},
			"gym":  {Name: "Gym", Color: "#ef9f76"},
		},
	}
}

func TestNewDefaults(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	if f.Kind != backend.RecurDaily {
		t.Errorf("New Kind = %q, want %q", f.Kind, backend.RecurDaily)
	}
	if f.Interval != "3" {
		t.Errorf("New Interval = %q, want 3", f.Interval)
	}
	if f.Task.Horizon != "small" {
		t.Errorf("New Horizon = %q, want small", f.Task.Horizon)
	}
}

func TestBeginLoadsTaskAndDefaultsInterval(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Begin(Task{Title: "Ship it", Kind: backend.RecurInterval, Recurring: true}, true)
	if f.Title != "Ship it" {
		t.Errorf("Begin Title = %q, want Ship it", f.Title)
	}
	if f.Interval != "3" {
		t.Errorf("Begin interval fallback = %q, want 3", f.Interval)
	}

	f.Begin(Task{Title: "x", Kind: backend.RecurInterval, Recurring: true, IntervalDays: 7}, true)
	if f.Interval != "7" {
		t.Errorf("Begin interval = %q, want 7", f.Interval)
	}
}

func TestSaveValidation(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "   "
	if cmd := f.saveCommand(); cmd != nil {
		t.Error("empty title should not produce a save command")
	} else if f.Error != "Title is required" {
		t.Errorf("empty title Error = %q, want %q", f.Error, "Title is required")
	}

	f = New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "Valid"
	f.Due = "not-a-date"
	if cmd := f.saveCommand(); cmd != nil {
		t.Error("invalid due date should not produce a save command")
	} else if f.Error != "Due date must use YYYY-MM-DD" {
		t.Errorf("invalid due date Error = %q", f.Error)
	}

	f = New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "Valid"
	f.Recurring = true
	f.Kind = backend.RecurInterval
	f.Interval = "0"
	if cmd := f.saveCommand(); cmd != nil {
		t.Error("invalid interval should not produce a save command")
	} else if f.Error != "Interval must be a positive number" {
		t.Errorf("invalid interval Error = %q", f.Error)
	}
}

func TestSaveWeekdaysDefault(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "Standup"
	f.Recurring = true
	f.Kind = backend.RecurWeekdays
	msg := runSave(t, &f)
	if msg.Task.Weekdays != backend.WeekdaysMonFri {
		t.Errorf("weekdays mask = %d, want %d", msg.Task.Weekdays, backend.WeekdaysMonFri)
	}
	if msg.Task.Kind != backend.RecurWeekdays {
		t.Errorf("kind = %q, want %q", msg.Task.Kind, backend.RecurWeekdays)
	}
}

func TestSaveInterval(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "Water plants"
	f.Recurring = true
	f.Kind = backend.RecurInterval
	f.Interval = "5"
	msg := runSave(t, &f)
	if msg.Task.IntervalDays != 5 {
		t.Errorf("IntervalDays = %d, want 5", msg.Task.IntervalDays)
	}
}

func TestSaveDueDateAndHorizon(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Title = "Report"
	f.Due = "2026-10-20"
	f.Task.Horizon = "long"
	msg := runSave(t, &f)
	if msg.Task.DueDate != "2026-10-20" {
		t.Errorf("DueDate = %q, want 2026-10-20", msg.Task.DueDate)
	}
	if msg.Task.Horizon != "long" {
		t.Errorf("Horizon = %q, want long", msg.Task.Horizon)
	}
}

func runSave(t *testing.T, f *Form) SavedMsg {
	t.Helper()
	cmd := f.saveCommand()
	if cmd == nil {
		t.Fatalf("saveCommand returned nil (error %q)", f.Error)
	}
	msg := cmd()
	saved, ok := msg.(SavedMsg)
	if !ok {
		t.Fatalf("saveCommand returned %T, want SavedMsg", msg)
	}
	return saved
}

func TestChangeChoiceLabel(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Focus = fieldLabel
	f.Task.LabelID = "work"
	f.changeChoice(true) // work -> gym
	if f.Task.LabelID != "gym" {
		t.Errorf("forward LabelID = %q, want gym", f.Task.LabelID)
	}
	f.changeChoice(false) // gym -> work
	if f.Task.LabelID != "work" {
		t.Errorf("backward LabelID = %q, want work", f.Task.LabelID)
	}
}

func TestChangeChoiceHorizon(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Focus = fieldHorizon
	if f.Task.Horizon != "small" {
		t.Fatalf("initial Horizon = %q, want small", f.Task.Horizon)
	}
	f.changeChoice(true)
	if f.Task.Horizon != "long" {
		t.Errorf("forward Horizon = %q, want long", f.Task.Horizon)
	}
	f.changeChoice(true)
	if f.Task.Horizon != "small" {
		t.Errorf("wrap Horizon = %q, want small", f.Task.Horizon)
	}
}

func TestChangeChoiceSchedule(t *testing.T) {
	f := New(backend.DefaultTheme(), testLabels(), 40, 20)
	f.Focus = fieldSchedule
	if f.Kind != backend.RecurDaily {
		t.Fatalf("initial Kind = %q, want %q", f.Kind, backend.RecurDaily)
	}
	f.changeChoice(true)
	if f.Kind != backend.RecurWeekdays {
		t.Errorf("forward Kind = %q, want %q", f.Kind, backend.RecurWeekdays)
	}
	f.changeChoice(true)
	if f.Kind != backend.RecurInterval {
		t.Errorf("forward Kind = %q, want %q", f.Kind, backend.RecurInterval)
	}
	f.changeChoice(true) // wraps around
	if f.Kind != backend.RecurDaily {
		t.Errorf("wrap Kind = %q, want %q", f.Kind, backend.RecurDaily)
	}
}
