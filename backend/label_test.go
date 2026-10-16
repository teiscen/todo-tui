package backend

import "testing"

func TestLabelsStepForward(t *testing.T) {
	l := Labels{
		Selected: "work",
		Info: map[LabelID]LabelInfo{
			"uni":  {Name: "University", Color: "#81c8be"},
			"work": {Name: "Work", Color: "#ca9ee6"},
			"gym":  {Name: "Gym", Color: "#ef9f76"},
		},
		Order: []LabelID{"uni", "work", "gym"},
	}

	l.Step(true)
	if l.Selected != "gym" {
		t.Errorf("Step(true) = %q, want %q", l.Selected, "gym")
	}

	l.Step(true)
	if l.Selected != "uni" {
		t.Errorf("Step(true) wrap = %q, want %q", l.Selected, "uni")
	}

	l.Step(false)
	if l.Selected != "gym" {
		t.Errorf("Step(false) wrap = %q, want %q", l.Selected, "gym")
	}

	l.Step(false)
	if l.Selected != "work" {
		t.Errorf("Step(false) = %q, want %q", l.Selected, "work")
	}
}

func TestLabelsStepEmpty(t *testing.T) {
	var l Labels
	l.Step(true) // must not panic on an empty order
	if l.Selected != "" {
		t.Errorf("Step on empty Labels set Selected to %q, want empty", l.Selected)
	}
}

func TestLabelsGetColorAndName(t *testing.T) {
	l := Labels{
		Selected: "work",
		Info:     map[LabelID]LabelInfo{"work": {Name: "Work", Color: "#ca9ee6"}},
	}
	if got := l.GetColor(); got != HexCode("#ca9ee6") {
		t.Errorf("GetColor() = %q, want %q", got, "#ca9ee6")
	}
	if got := l.GetName(); got != "Work" {
		t.Errorf("GetName() = %q, want %q", got, "Work")
	}
}
