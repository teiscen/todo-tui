package database

import (
	"database/sql"
	"testing"

	"todo-tui/backend"
)

func TestSaveTaskInsertAndUpdate(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	id, err := SaveTask(db, Task{Title: "draft", Horizon: "small", Notes: "n1"})
	if err != nil {
		t.Fatalf("SaveTask insert: %v", err)
	}
	if id == 0 {
		t.Fatal("SaveTask insert returned id 0")
	}

	tasks, err := FetchTasks(db)
	if err != nil {
		t.Fatalf("FetchTasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("FetchTasks returned %d tasks, want 1", len(tasks))
	}
	if tasks[0].Title != "draft" || tasks[0].Notes != "n1" || tasks[0].Horizon != "small" {
		t.Errorf("inserted task = %+v, want draft/n1/small", tasks[0])
	}
	if tasks[0].Recurring {
		t.Error("non-recurring task reported recurring")
	}

	if _, err := SaveTask(db, Task{ID: tasks[0].ID, Title: "final", Horizon: "long", Notes: "n2"}); err != nil {
		t.Fatalf("SaveTask update: %v", err)
	}
	tasks, err = FetchTasks(db)
	if err != nil {
		t.Fatalf("FetchTasks after update: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Title != "final" || tasks[0].Horizon != "long" || tasks[0].Notes != "n2" {
		t.Errorf("updated task = %+v, want final/long/n2", tasks[0])
	}
}

func TestSaveTaskRecurrence(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	id, err := SaveTask(db, Task{
		Title:     "standup",
		Horizon:   "small",
		Recurring: true,
		Kind:      backend.RecurWeekdays,
		Weekdays:  backend.WeekdaysMonFri,
	})
	if err != nil {
		t.Fatalf("SaveTask recurring: %v", err)
	}

	tasks, err := FetchTasks(db)
	if err != nil {
		t.Fatalf("FetchTasks: %v", err)
	}
	if len(tasks) != 1 || !tasks[0].Recurring || tasks[0].Kind != backend.RecurWeekdays || tasks[0].Weekdays != backend.WeekdaysMonFri {
		t.Errorf("recurring task = %+v, want weekdays/%d", tasks[0], backend.WeekdaysMonFri)
	}

	// Removing recurrence on update drops the recurrence row.
	if _, err := SaveTask(db, Task{ID: int(id), Title: "standup", Horizon: "small"}); err != nil {
		t.Fatalf("SaveTask remove recurrence: %v", err)
	}
	tasks, err = FetchTasks(db)
	if err != nil {
		t.Fatalf("FetchTasks after removing recurrence: %v", err)
	}
	if tasks[0].Recurring {
		t.Errorf("task still recurring after removing recurrence: %+v", tasks[0])
	}
}

func TestToggleTaskOneTime(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	id, err := SaveTask(db, Task{Title: "one", Horizon: "small"})
	if err != nil {
		t.Fatalf("SaveTask: %v", err)
	}

	// done=false means "currently not done", so this marks it done.
	if err := ToggleTask(db, int(id), false, false); err != nil {
		t.Fatalf("ToggleTask mark: %v", err)
	}
	tasks, _ := FetchTasks(db)
	if !tasks[0].DoneAt.Valid {
		t.Error("expected done_at set after marking done")
	}

	// done=true means "currently done", so this unmarks it.
	if err := ToggleTask(db, int(id), false, true); err != nil {
		t.Fatalf("ToggleTask unmark: %v", err)
	}
	tasks, _ = FetchTasks(db)
	if tasks[0].DoneAt.Valid {
		t.Error("expected done_at cleared after unmarking")
	}
}

func TestToggleTaskRecurring(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	id, err := SaveTask(db, Task{Title: "habit", Horizon: "small", Recurring: true, Kind: backend.RecurDaily})
	if err != nil {
		t.Fatalf("SaveTask: %v", err)
	}

	done, err := IsTaskDoneToday(db, int(id), true, sql.NullString{})
	if err != nil {
		t.Fatalf("IsTaskDoneToday: %v", err)
	}
	if done {
		t.Fatal("expected not done today initially")
	}

	if err := ToggleTask(db, int(id), true, false); err != nil {
		t.Fatalf("ToggleTask mark: %v", err)
	}
	done, _ = IsTaskDoneToday(db, int(id), true, sql.NullString{})
	if !done {
		t.Error("expected done today after marking")
	}

	if err := ToggleTask(db, int(id), true, true); err != nil {
		t.Fatalf("ToggleTask unmark: %v", err)
	}
	done, _ = IsTaskDoneToday(db, int(id), true, sql.NullString{})
	if done {
		t.Error("expected not done today after unmarking")
	}
}

func TestIsTaskDoneTodayOneTime(t *testing.T) {
	db := openTestDB(t)

	done, err := IsTaskDoneToday(db, 1, false, sql.NullString{String: "2026-01-01T00:00:00Z", Valid: true})
	if err != nil {
		t.Fatalf("IsTaskDoneToday: %v", err)
	}
	if !done {
		t.Error("one-time task with done_at should report done")
	}

	done, err = IsTaskDoneToday(db, 1, false, sql.NullString{})
	if err != nil {
		t.Fatalf("IsTaskDoneToday: %v", err)
	}
	if done {
		t.Error("one-time task without done_at should report not done")
	}
}
