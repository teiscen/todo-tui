package planner

import (
	"fmt"
	"strings"
	"testing"
	"time"

	backend "todo-tui/backend"

	"charm.land/lipgloss/v2"
)

func testPlanner() Planner {
	return New(backend.DefaultTheme(), backend.Labels{
		Selected: "work",
		Order:    []backend.LabelID{"uni", "work", "gym"},
		Info: map[backend.LabelID]backend.LabelInfo{
			"uni":  {Name: "University", Color: "#81c8be"},
			"work": {Name: "Work", Color: "#ca9ee6"},
			"gym":  {Name: "Gym", Color: "#ef9f76"},
		},
	}, backend.NewDate(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), 43)
}

func d(y int, m time.Month, day int) backend.Date {
	return backend.NewDate(time.Date(y, m, day, 0, 0, 0, 0, time.UTC))
}

// borderPrefix is the ANSI color sequence immediately before a vertical border
// glyph when it is correctly colored.
func borderPrefix(c backend.HexCode) string {
	r, g, b, err := c.ParseRGB()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm│", r, g, b)
}

func TestViewColorsBorders(t *testing.T) {
	p := testPlanner()
	if prefix := borderPrefix(p.Style.Theme.Border); !strings.Contains(p.View(), prefix) {
		t.Errorf("calendar/notes border not colored with theme border color %q", prefix)
	}
}

func TestViewMutesBordersWhenDeselected(t *testing.T) {
	p := testPlanner()
	p.Style.Theme.Border = backend.HexCode("#737994")
	if prefix := borderPrefix(p.Style.Theme.Border); !strings.Contains(p.View(), prefix) {
		t.Errorf("calendar/notes border not muted with %q when deselected", prefix)
	}
}

func TestViewHeightMatchesRequest(t *testing.T) {
	p := testPlanner()
	if h := lipgloss.Height(p.View()); h != 43 {
		t.Errorf("planner height = %d, want 43", h)
	}
}

func TestToggleFocus(t *testing.T) {
	p := testPlanner()
	if p.Focus != FocusCalendar {
		t.Fatalf("initial focus = %d, want FocusCalendar", p.Focus)
	}
	p.ToggleFocus()
	if p.Focus != FocusNotes {
		t.Errorf("focus after toggle = %d, want FocusNotes", p.Focus)
	}
	if !p.NotesModel.TextArea.Focused() {
		t.Error("notes textarea should be focused after toggle")
	}
	p.ToggleFocus()
	if p.Focus != FocusCalendar {
		t.Errorf("focus after second toggle = %d, want FocusCalendar", p.Focus)
	}
	if p.NotesModel.TextArea.Focused() {
		t.Error("notes textarea should be blurred after second toggle")
	}
}

func TestUpdateMonth(t *testing.T) {
	p := testPlanner()
	before := *p.SelectedDate()
	cmd := p.UpdateMonth(true)
	if after := *p.SelectedDate(); after != before.AddMonths(1) {
		t.Errorf("SelectedDate after next = %v, want %v", after, before.AddMonths(1))
	}
	msg := cmd()
	cm, ok := msg.(ChangeCalendarMsg)
	if !ok || cm.NewDate != *p.SelectedDate() || cm.LabelID != *p.SelectedLabel() {
		t.Errorf("next cmd = %#v, want ChangeCalendarMsg for %v/%q", msg, *p.SelectedDate(), *p.SelectedLabel())
	}

	p.UpdateMonth(false)
	if after := *p.SelectedDate(); after != before {
		t.Errorf("SelectedDate after prev = %v, want %v", after, before)
	}
}

func TestUpdateLabel(t *testing.T) {
	p := testPlanner() // Selected "work", order [uni work gym]
	cmd := p.UpdateLabel(true)
	if got := *p.SelectedLabel(); got != "gym" {
		t.Errorf("SelectedLabel after next = %q, want gym", got)
	}
	if got := p.CalendarModel.Style.Theme.Accent; got != backend.HexCode("#ef9f76") {
		t.Errorf("calendar accent = %q, want gym color", got)
	}
	if cm, ok := cmd().(ChangeCalendarMsg); !ok || cm.LabelID != "gym" {
		t.Errorf("next cmd = %#v, want ChangeCalendarMsg for gym", cmd())
	}

	p.UpdateLabel(false)
	if got := *p.SelectedLabel(); got != "work" {
		t.Errorf("SelectedLabel after prev = %q, want work", got)
	}
}

func TestReloadState(t *testing.T) {
	p := testPlanner()
	sel := *p.SelectedDate()
	entries := map[string]string{sel.Key(): "note text"}

	p.ReloadState(sel, entries)

	if !p.HasEntry(sel) {
		t.Error("expected an entry for the selected date after reload")
	}
	if info := p.CalendarModel.GetInfo(sel); info == nil || !info.IsMarked {
		t.Error("expected selected date to be marked in the grid")
	}
	if got := p.NotesModel.TextArea.Value(); got != "note text" {
		t.Errorf("notes value = %q, want %q", got, "note text")
	}
}

func TestUpdateCalendarState(t *testing.T) {
	p := testPlanner()
	start := d(2026, time.January, 1)
	p.Entries = map[string]string{d(2026, time.January, 3).Key(): "x"}
	p.UpdateCalendarState(start)

	if p.CalendarModel.State[0][0].Date != start {
		t.Errorf("grid[0][0] = %v, want %v", p.CalendarModel.State[0][0].Date, start)
	}
	if !p.CalendarModel.State[0][2].IsMarked {
		t.Error("grid[0][2] (Jan 3) should be marked")
	}
	if p.CalendarModel.State[0][1].IsMarked {
		t.Error("grid[0][1] (Jan 2) should not be marked")
	}
}

func TestSyncNotes(t *testing.T) {
	p := testPlanner()
	sel := *p.SelectedDate()

	p.Entries = map[string]string{sel.Key(): "hello"}
	p.SyncNotes()
	if got := p.NotesModel.TextArea.Value(); got != "hello" {
		t.Errorf("notes value = %q, want %q", got, "hello")
	}

	p.Entries = map[string]string{}
	p.SyncNotes()
	if got := p.NotesModel.TextArea.Value(); got != "" {
		t.Errorf("notes value = %q, want empty after clearing entry", got)
	}
}
