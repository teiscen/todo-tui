package planner

import (
	"fmt"
	"strings"
	"testing"
	"time"

	backend "todo-tui/backend"
)

func testPlanner() Planner {
	return New(backend.DefaultTheme(), backend.Labels{
		Selected: "work",
		Order:    []backend.LabelID{"work"},
		Info:     map[backend.LabelID]backend.LabelInfo{"work": {Name: "Work", Color: "#ca9ee6"}},
	}, backend.NewDate(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
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
