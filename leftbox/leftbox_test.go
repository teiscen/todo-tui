package leftbox

import (
	"strings"
	"testing"

	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func testLabels() backend.Labels {
	return backend.Labels{
		Selected: "work",
		Order:    []backend.LabelID{"uni", "work"},
		Info: map[backend.LabelID]backend.LabelInfo{
			"uni":  {Name: "University", Color: "#81c8be"},
			"work": {Name: "Work", Color: "#ca9ee6"},
		},
	}
}

// testLeftbox returns a small deterministic list: four rows across three
// groups, with row 2 expanded so a detail line is interleaved in the body.
func testLeftbox() Leftbox {
	l := New(backend.DefaultTheme(), testLabels(), 50, 20)
	l.Rows = []Row{
		{ID: 1, Tab: TabSmall, Group: GroupOverdue, Title: "one"},
		{ID: 2, Tab: TabSmall, Group: GroupToday, Title: "two"},
		{ID: 3, Tab: TabSmall, Group: GroupToday, Title: "three"},
		{ID: 4, Tab: TabSmall, Group: GroupLater, Title: "four"},
	}
	l.Expanded = map[int]bool{2: true}
	return l
}

func manyRows(n int) []Row {
	rows := make([]Row, n)
	for i := range rows {
		rows[i] = Row{ID: i + 1, Tab: TabSmall, Group: GroupToday, Title: "task"}
	}
	return rows
}

func keyPress(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Text: text, Code: []rune(text)[0]}
}

// TestPaneSizeInvariant pins the panel to its configured width/height no matter
// how much (or how little) content it holds, so the dashboard columns align.
func TestPaneSizeInvariant(t *testing.T) {
	rowsets := [][]Row{
		nil,
		{{ID: 1, Tab: TabSmall, Group: GroupToday, Title: "one"}},
		testLeftbox().Rows,
		manyRows(40), // scrolls the body
	}
	for _, rows := range rowsets {
		l := New(backend.DefaultTheme(), testLabels(), 50, 20)
		l.Rows = rows
		view := l.View()
		if w, h := lipgloss.Width(view), lipgloss.Height(view); w != l.Style.Width || h != l.Style.Height {
			t.Errorf("View() = %dx%d, want %dx%d (rows=%d)", w, h, l.Style.Width, l.Style.Height, len(rows))
		}
	}
}

// TestHoverMatchesLayout verifies that the line a task is drawn on is exactly
// the line HoverAt maps back to that task, and that non-item lines (headings,
// blanks, expanded detail) never select anything.
func TestHoverMatchesLayout(t *testing.T) {
	l := testLeftbox()
	contentW := l.Style.Width - 4
	header := l.header(contentW)
	bodyH := l.Style.Height - 2 - len(header)

	body, cursorLine := l.layout(contentW)
	start := windowStart(len(body), cursorLine, bodyH)
	end := min(len(body), start+bodyH)

	for i := start; i < end; i++ {
		y := 1 + len(header) + (i - start)
		probe := l
		probe.Cursor = -1
		if body[i].kind == itemLine {
			if !probe.HoverAt(y) {
				t.Errorf("HoverAt(%d) = false, want true for item index %d", y, body[i].index)
			} else if probe.Cursor != body[i].index {
				t.Errorf("HoverAt(%d) selected %d, want %d", y, probe.Cursor, body[i].index)
			}
		} else if probe.HoverAt(y) {
			t.Errorf("HoverAt(%d) = true on non-item line (kind %d), want false", y, body[i].kind)
		}
	}
}

// TestCollapseKeepsHeadingReachable verifies a collapsed group still draws its
// heading (with the collapsed marker) and that a selection stays valid.
func TestCollapseKeepsHeadingReachable(t *testing.T) {
	l := testLeftbox()
	l.Cursor = 1 // the first TODAY row (ID 2)
	l.Collapsed[GroupToday] = true
	l.clampCursor()

	contentW := l.Style.Width - 4
	body, _ := l.layout(contentW)

	var headings, items int
	collapsedHeading := false
	for _, line := range body {
		switch line.kind {
		case headingLine:
			headings++
			if strings.Contains(line.text, "TODAY") && strings.Contains(line.text, "▸") {
				collapsedHeading = true
			}
		case itemLine:
			items++
		}
	}
	if headings != 3 {
		t.Errorf("headings = %d, want 3 (collapsed group keeps its heading)", headings)
	}
	if items != 2 {
		t.Errorf("items = %d, want 2 after collapsing TODAY", items)
	}
	if !collapsedHeading {
		t.Error("collapsed TODAY heading should show the ▸ marker")
	}
	if _, ok := l.Selected(); !ok {
		t.Error("Selected() = false after collapse, want a reachable row")
	}
}

// TestReloadKeepsSelectionByID verifies reload follows the selected task even
// when it moves, and clamps into range when it disappears.
func TestReloadKeepsSelectionByID(t *testing.T) {
	l := testLeftbox()
	l.Cursor = 1 // select ID 2

	l.ReloadState([]Row{
		{ID: 2, Tab: TabSmall, Group: GroupLater, Title: "two"}, // moved to the end
		{ID: 5, Tab: TabSmall, Group: GroupToday, Title: "five"},
		{ID: 3, Tab: TabSmall, Group: GroupToday, Title: "three"},
		{ID: 1, Tab: TabSmall, Group: GroupOverdue, Title: "one"},
	})
	if sel, ok := l.Selected(); !ok || sel.ID != 2 {
		t.Errorf("ReloadState selected ID %d (ok=%v), want 2 after the task moves", sel.ID, ok)
	}

	l.ReloadState([]Row{{ID: 9, Tab: TabSmall, Group: GroupToday, Title: "nine"}})
	if sel, ok := l.Selected(); !ok || sel.ID != 9 {
		t.Errorf("ReloadState selected ID %d (ok=%v), want 9 after the task disappears", sel.ID, ok)
	}
}

func TestTabCyclesAndResetsCursor(t *testing.T) {
	l := testLeftbox()
	l.Cursor = 2

	out, cmd := l.Update(keyPress("]"))
	if cmd != nil {
		t.Error("tab switch should not emit a command")
	}
	if out.Tab != TabLong {
		t.Errorf("after ] Tab = %d, want TabLong", out.Tab)
	}
	if out.Cursor != 0 {
		t.Errorf("after ] Cursor = %d, want 0", out.Cursor)
	}

	out, _ = out.Update(keyPress("["))
	if out.Tab != TabSmall {
		t.Errorf("after [ Tab = %d, want TabSmall", out.Tab)
	}
}

func TestFilterCycles(t *testing.T) {
	l := testLeftbox()
	if l.Filter != "" {
		t.Fatalf("initial filter = %q, want empty", l.Filter)
	}
	for i, want := range []backend.LabelID{"uni", "work", ""} {
		out, cmd := l.Update(keyPress("f"))
		if cmd != nil {
			t.Error("filter step should not emit a command")
		}
		if out.Filter != want {
			t.Errorf("step %d filter = %q, want %q", i, out.Filter, want)
		}
		if out.Cursor != 0 {
			t.Errorf("step %d Cursor = %d, want 0", i, out.Cursor)
		}
		l = out
	}
}

func TestSpaceEmitsToggle(t *testing.T) {
	l := testLeftbox()
	l.Cursor = 1 // ID 2, one-time, not done
	_, cmd := l.Update(keyPress(" "))
	if cmd == nil {
		t.Fatal("space should emit a toggle command")
	}
	toggle, ok := cmd().(ToggleTodoMsg)
	if !ok {
		t.Fatalf("space emitted %T, want ToggleTodoMsg", cmd())
	}
	if toggle.ID != 2 || toggle.Recurring || toggle.Done {
		t.Errorf("toggle = %+v, want ID 2, not recurring, not done", toggle)
	}
}

func TestSpaceEmitsToggleRecurring(t *testing.T) {
	l := New(backend.DefaultTheme(), testLabels(), 50, 20)
	l.Tab = TabRecurring
	l.Rows = []Row{{ID: 7, Tab: TabRecurring, Group: GroupToday, Title: "habit", Done: true}}
	_, cmd := l.Update(keyPress(" "))
	if cmd == nil {
		t.Fatal("space should emit a toggle command")
	}
	toggle, ok := cmd().(ToggleTodoMsg)
	if !ok {
		t.Fatalf("space emitted %T, want ToggleTodoMsg", cmd())
	}
	if toggle.ID != 7 || !toggle.Recurring || !toggle.Done {
		t.Errorf("toggle = %+v, want ID 7, recurring, done", toggle)
	}
}
