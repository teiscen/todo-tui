package todo

import (
	"fmt"
	"strings"
	"testing"

	"todo-tui/backend"
	"todo-tui/config"
	"todo-tui/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func testLabels() backend.Labels {
	return backend.Labels{
		Selected: "uni",
		Order:    []backend.LabelID{"uni", "work"},
		Info: map[backend.LabelID]backend.LabelInfo{
			"uni":  {Name: "University", Color: "#81c8be"},
			"work": {Name: "Work", Color: "#ca9ee6"},
		},
	}
}

func testConfig() config.Config {
	return config.Config{
		Layout: config.LayoutConfig{
			Tasks:   config.PanelSize{Width: 50, Height: 43},
			Preview: config.PanelSize{Width: 50, Height: 43},
			Editor:  config.PanelSize{Width: 50, Height: 43},
		},
		Labels: testLabels(),
	}
}

func newTestTodo(t *testing.T) *Todo {
	t.Helper()
	db, err := database.ConnectDB(":memory:")
	if err != nil {
		t.Fatalf("ConnectDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := database.CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}
	for i := 1; i <= 3; i++ {
		if _, err := database.SaveTask(db, database.Task{Title: fmt.Sprintf("task %d", i), Horizon: "small"}); err != nil {
			t.Fatalf("SaveTask %d: %v", i, err)
		}
	}
	return NewTodo(testConfig(), db)
}

func keyPress(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Text: text, Code: []rune(text)[0]}
}

// TestKeysReachOnlyFocusedPane verifies "j" affects the tasks list when the
// task pane is focused and the calendar when the planner is focused — never both.
func TestKeysReachOnlyFocusedPane(t *testing.T) {
	td := newTestTodo(t)

	td.focus = PaneTasks
	dateBefore := *td.plannerModel.SelectedDate()
	td.Update(keyPress("j"))
	if td.leftboxModel.Cursor != 1 {
		t.Errorf("leftbox cursor = %d, want 1 after j", td.leftboxModel.Cursor)
	}
	if after := *td.plannerModel.SelectedDate(); after != dateBefore {
		t.Errorf("planner date moved %v -> %v while tasks focused", dateBefore, after)
	}

	td.focus = PanePlanner
	cursorBefore := td.leftboxModel.Cursor
	td.Update(keyPress("j"))
	if td.leftboxModel.Cursor != cursorBefore {
		t.Errorf("leftbox cursor = %d, want %d while planner focused", td.leftboxModel.Cursor, cursorBefore)
	}
	if after := *td.plannerModel.SelectedDate(); after == dateBefore {
		t.Error("planner date did not move after j while planner focused")
	}
}

// TestPreviewFollowsCursor verifies the preview pane tracks the selected task
// as the left panel cursor moves.
func TestPreviewFollowsCursor(t *testing.T) {
	td := newTestTodo(t)

	td.leftboxModel.Cursor = 0
	first, ok1 := td.leftboxModel.Selected()
	td.leftboxModel.Cursor = 1
	second, ok2 := td.leftboxModel.Selected()
	td.leftboxModel.Cursor = 0
	if !ok1 || !ok2 {
		t.Fatal("need at least 2 selectable rows")
	}
	if td.previewModel.Selected.ID != first.ID {
		t.Fatalf("initial preview ID = %d, want %d", td.previewModel.Selected.ID, first.ID)
	}

	td.focus = PaneTasks
	td.Update(keyPress("j"))
	if td.previewModel.Selected.ID != second.ID {
		t.Errorf("preview ID after j = %d, want %d", td.previewModel.Selected.ID, second.ID)
	}
}

// TestPaneHeightsMatch verifies the tasks, preview, and planner panes all
// render to the same configured height so their borders stay aligned.
func TestPaneHeightsMatch(t *testing.T) {
	td := newTestTodo(t)
	tasksH := lipgloss.Height(td.leftboxModel.View())
	previewH := lipgloss.Height(td.previewModel.View())
	plannerH := lipgloss.Height(td.plannerModel.View())
	if tasksH != previewH || tasksH != plannerH {
		t.Errorf("pane heights differ: tasks=%d preview=%d planner=%d", tasksH, previewH, plannerH)
	}
	if tasksH != td.leftboxModel.Style.Height {
		t.Errorf("tasks height %d != configured %d", tasksH, td.leftboxModel.Style.Height)
	}
}

// TestStaleLoadsDropped verifies an async planner result for a month/label the
// user has left is dropped, while a current one is applied.
func TestStaleLoadsDropped(t *testing.T) {
	td := newTestTodo(t)
	label := *td.plannerModel.SelectedLabel()
	date := *td.plannerModel.SelectedDate()

	td.Update(plannerLoadedMsg{
		LabelID: label,
		NewDate: date.AddMonths(1),
		sDate:   date.AddMonths(1),
		entries: map[string]string{"2026-10-15": "stale"},
	})
	if len(td.plannerModel.Entries) != 0 {
		t.Errorf("stale load populated entries: %v", td.plannerModel.Entries)
	}

	td.Update(plannerLoadedMsg{
		LabelID: label,
		NewDate: date,
		sDate:   date,
		entries: map[string]string{"2026-10-15": "fresh"},
	})
	if td.plannerModel.Entries["2026-10-15"] != "fresh" {
		t.Errorf("fresh load entries = %v, want the fresh entry", td.plannerModel.Entries)
	}
}

// TestTypingGuardReservesKeysForEditor verifies navigation keys do not change
// pane focus while the task editor is open.
func TestTypingGuardReservesKeysForEditor(t *testing.T) {
	td := newTestTodo(t)
	td.focus = PaneTasks
	td.Update(keyPress("a"))
	if !td.formActive {
		t.Fatal("expected the editor to open after 'a'")
	}

	focusBefore := td.focus
	td.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if td.focus != focusBefore {
		t.Errorf("tab changed focus to %d while editing, want %d", td.focus, focusBefore)
	}
}

// TestWithHelpWrapsToPaneWidth verifies the keyboard hints below a panel wrap
// (instead of widening the column or being truncated).
func TestWithHelpWrapsToPaneWidth(t *testing.T) {
	content := strings.Repeat("x", 30)
	movement := "a one · b two · c three · d four"
	specific := "e five · f six · g seven"
	out := withHelp(content, false, movement, specific)

	if w := lipgloss.Width(out); w != 30 {
		t.Errorf("withHelp width = %d, want 30", w)
	}
	for _, want := range []string{"a one", "b two", "c three", "d four", "e five", "f six", "g seven"} {
		if !strings.Contains(out, want) {
			t.Errorf("withHelp output missing %q (truncated?)", want)
		}
	}
}
