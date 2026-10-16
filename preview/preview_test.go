package preview

import (
	"strings"
	"testing"

	"todo-tui/backend"

	"charm.land/lipgloss/v2"
)

func TestPreviewEmptyState(t *testing.T) {
	p := New(backend.DefaultTheme(), 50, 20)
	p.SetTask(Detail{}, false)
	out := p.View()
	if !strings.Contains(out, "Select a task") {
		t.Errorf("empty state missing prompt: %q", out)
	}
}

func TestPreviewFieldRendering(t *testing.T) {
	p := New(backend.DefaultTheme(), 50, 43)
	p.SetTask(Detail{
		ID:         42,
		Title:      "My Task",
		LabelID:    "work",
		Color:      "#ca9ee6",
		Schedule:   "M/W/F",
		Due:        "today",
		DueTone:    backend.ToneSoon,
		StepsDone:  1,
		StepsTotal: 3,
	}, true)

	out := p.View()
	for _, want := range []string{"My Task", "work", "M/W/F", "today", "1 of 3", "#42"} {
		if !strings.Contains(out, want) {
			t.Errorf("field output missing %q", want)
		}
	}
}

func TestPreviewSizeInvariant(t *testing.T) {
	for _, has := range []bool{false, true} {
		p := New(backend.DefaultTheme(), 50, 20)
		p.SetTask(Detail{ID: 1, Title: "x", StepsDone: 1, StepsTotal: 3}, has)
		if w, h := lipgloss.Width(p.View()), lipgloss.Height(p.View()); w != 50 || h != 20 {
			t.Errorf("size = %dx%d, want 50x20 (hasTask=%v)", w, h, has)
		}
	}
}
