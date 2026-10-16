package database

import (
	"database/sql"
	"fmt"
	"time"

	"todo-tui/backend"
)

// Task is the database representation used by the task list and editor.
type Task struct {
	ID           int
	Title        string
	LabelID      sql.NullString
	Horizon      string
	DueDate      sql.NullString
	Notes        string
	DoneAt       sql.NullString
	Recurring    bool
	Kind         string
	Weekdays     int
	IntervalDays sql.NullInt64
	StartDate    string
}

// FetchTasks loads tasks together with their optional recurrence rule.
func FetchTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query(`
		SELECT t.id, t.title, t.label_id, t.horizon, t.due_date, t.notes, t.done_at,
		       r.kind, r.weekdays, r.interval_days, r.start_date
		FROM todos AS t
		LEFT JOIN recurrence AS r ON r.todo_id = t.id
		ORDER BY CASE WHEN t.done_at IS NULL THEN 0 ELSE 1 END, t.id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("fetch tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	return tasks, nil
}

// scanTask converts one joined database row into a Task, including nullable recurrence fields.
func scanTask(rows *sql.Rows) (Task, error) {
	var task Task
	var kind, startDate sql.NullString
	var weekdays sql.NullInt64

	err := rows.Scan(
		&task.ID, &task.Title, &task.LabelID, &task.Horizon,
		&task.DueDate, &task.Notes, &task.DoneAt,
		&kind, &weekdays, &task.IntervalDays, &startDate,
	)
	if err != nil {
		return Task{}, fmt.Errorf("scan task: %w", err)
	}

	task.Recurring = kind.Valid
	task.Kind = kind.String
	task.Weekdays = int(weekdays.Int64)
	task.StartDate = startDate.String
	return task, nil
}

// SaveTask inserts or updates a task and its recurrence rule in one transaction.
// Removing recurrence from the form also removes the existing recurrence row.
func SaveTask(db *sql.DB, task Task) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin task save: %w", err)
	}
	defer tx.Rollback()

	id, err := saveTaskRow(tx, task)
	if err != nil {
		return 0, err
	}
	if err := saveRecurrence(tx, id, task); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit task save: %w", err)
	}
	return id, nil
}

func saveTaskRow(tx *sql.Tx, task Task) (int64, error) {
	var labelID, dueDate any
	if task.LabelID.Valid && task.LabelID.String != "" {
		labelID = task.LabelID.String
	}
	if task.DueDate.Valid && task.DueDate.String != "" {
		dueDate = task.DueDate.String
	}

	if task.ID == 0 {
		result, err := tx.Exec(
			"INSERT INTO todos (title, label_id, horizon, due_date, notes) VALUES (?, ?, ?, ?, ?)",
			task.Title, labelID, task.Horizon, dueDate, task.Notes,
		)
		if err != nil {
			return 0, fmt.Errorf("insert task: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("get inserted task ID: %w", err)
		}
		return id, nil
	}

	result, err := tx.Exec(
		"UPDATE todos SET title = ?, label_id = ?, horizon = ?, due_date = ?, notes = ? WHERE id = ?",
		task.Title, labelID, task.Horizon, dueDate, task.Notes, task.ID,
	)
	if err != nil {
		return 0, fmt.Errorf("update task %d: %w", task.ID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("check updated task: %w", err)
	}
	if affected == 0 {
		return 0, fmt.Errorf("task %d not found", task.ID)
	}
	return int64(task.ID), nil
}

func saveRecurrence(tx *sql.Tx, id int64, task Task) error {
	if !task.Recurring {
		if _, err := tx.Exec("DELETE FROM recurrence WHERE todo_id = ?", id); err != nil {
			return fmt.Errorf("remove task recurrence: %w", err)
		}
		return nil
	}

	kind := task.Kind
	if kind != backend.RecurDaily && kind != backend.RecurWeekdays && kind != backend.RecurInterval {
		kind = backend.RecurDaily
	}
	weekdays := task.Weekdays
	if kind == backend.RecurWeekdays && weekdays == 0 {
		weekdays = backend.WeekdaysMonFri
	}

	var interval any
	if kind == backend.RecurInterval {
		days := int64(1)
		if task.IntervalDays.Valid && task.IntervalDays.Int64 > 0 {
			days = task.IntervalDays.Int64
		}
		interval = days
	}
	startDate := task.StartDate
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}

	_, err := tx.Exec(`
		INSERT INTO recurrence (todo_id, kind, weekdays, interval_days, start_date)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(todo_id) DO UPDATE SET
			kind = excluded.kind,
			weekdays = excluded.weekdays,
			interval_days = excluded.interval_days,
			start_date = excluded.start_date
	`, id, kind, weekdays, interval, startDate)
	if err != nil {
		return fmt.Errorf("save task recurrence: %w", err)
	}
	return nil
}

// ToggleTask records completion for today for recurring tasks, or toggles done_at
// for one-time tasks. The Done argument is the row's current state.
func ToggleTask(db *sql.DB, id int, recurring, done bool) error {
	if recurring {
		today := time.Now().Format("2006-01-02")
		if done {
			_, err := db.Exec("DELETE FROM completions WHERE todo_id = ? AND date = ?", id, today)
			return err
		}
		_, err := db.Exec("INSERT OR IGNORE INTO completions (todo_id, date) VALUES (?, ?)", id, today)
		return err
	}

	if done {
		_, err := db.Exec("UPDATE todos SET done_at = NULL WHERE id = ?", id)
		return err
	}
	_, err := db.Exec("UPDATE todos SET done_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now') WHERE id = ?", id)
	return err
}

// IsTaskDoneToday checks today's completion record for recurring tasks and done_at
// for one-time tasks.
func IsTaskDoneToday(db *sql.DB, id int, recurring bool, doneAt sql.NullString) (bool, error) {
	if !recurring {
		return doneAt.Valid, nil
	}

	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM completions WHERE todo_id = ? AND date = ?",
		id, time.Now().Format("2006-01-02"),
	).Scan(&count)
	return count > 0, err
}
