package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID int
	Title string
	LabelID sql.NullString
	Horizon string
	DueDate sql.NullString
	Notes string
	DoneAt sql.NullString
	Recurring bool
	Kind string
	Weekdays int
	IntervalDays sql.NullInt64
	StartDate string
}

func FetchTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT t.id, t.title, t.label_id, t.horizon, t.due_date, t.notes, t.done_at, r.kind, r.weekdays, r.interval_days, r.start_date FROM todos t LEFT JOIN recurrence r ON r.todo_id=t.id ORDER BY CASE WHEN t.done_at IS NULL THEN 0 ELSE 1 END, t.id DESC")
	if err != nil { return nil, fmt.Errorf("fetch tasks: %w", err) }
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var t Task
		var kind, start sql.NullString
		var weekdays sql.NullInt64
		if err := rows.Scan(&t.ID, &t.Title, &t.LabelID, &t.Horizon, &t.DueDate, &t.Notes, &t.DoneAt, &kind, &weekdays, &t.IntervalDays, &start); err != nil { return nil, fmt.Errorf("scan task: %w", err) }
		t.Recurring, t.Kind, t.Weekdays, t.StartDate = kind.Valid, kind.String, int(weekdays.Int64), start.String
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil { return nil, fmt.Errorf("read tasks: %w", err) }
	return tasks, nil
}

func SaveTask(db *sql.DB, t Task) (int64, error) {
	tx, err := db.Begin()
	if err != nil { return 0, err }
	defer tx.Rollback()
	var label, due any
	if t.LabelID.Valid && t.LabelID.String != "" { label = t.LabelID.String }
	if t.DueDate.Valid && t.DueDate.String != "" { due = t.DueDate.String }
	var id int64
	if t.ID == 0 {
		res, err := tx.Exec("INSERT INTO todos (title,label_id,horizon,due_date,notes) VALUES (?,?,?,?,?)", t.Title, label, t.Horizon, due, t.Notes)
		if err != nil { return 0, fmt.Errorf("insert task: %w", err) }
		id, err = res.LastInsertId()
		if err != nil { return 0, err }
	} else {
		id = int64(t.ID)
		res, err := tx.Exec("UPDATE todos SET title=?,label_id=?,horizon=?,due_date=?,notes=? WHERE id=?", t.Title, label, t.Horizon, due, t.Notes, id)
		if err != nil { return 0, fmt.Errorf("update task: %w", err) }
		n, err := res.RowsAffected()
		if err != nil { return 0, err }
		if n == 0 { return 0, fmt.Errorf("task %d not found", id) }
	}
	if t.Recurring {
		kind := t.Kind
		if kind != "daily" && kind != "weekdays" && kind != "interval" { kind = "daily" }
		weekdays := t.Weekdays
		if kind == "weekdays" && weekdays == 0 { weekdays = 62 }
		var interval any
		if kind == "interval" {
			n := int64(1)
			if t.IntervalDays.Valid && t.IntervalDays.Int64 > 0 { n = t.IntervalDays.Int64 }
			interval = n
		}
		start := t.StartDate
		if start == "" { start = time.Now().Format("2006-01-02") }
		_, err := tx.Exec("INSERT INTO recurrence (todo_id,kind,weekdays,interval_days,start_date) VALUES (?,?,?,?,?) ON CONFLICT(todo_id) DO UPDATE SET kind=excluded.kind, weekdays=excluded.weekdays, interval_days=excluded.interval_days, start_date=excluded.start_date", id, kind, weekdays, interval, start)
		if err != nil { return 0, fmt.Errorf("save recurrence: %w", err) }
	} else {
		if _, err := tx.Exec("DELETE FROM recurrence WHERE todo_id=?", id); err != nil { return 0, fmt.Errorf("remove recurrence: %w", err) }
	}
	if err := tx.Commit(); err != nil { return 0, err }
	return id, nil
}


func ToggleTask(db *sql.DB, id int, recurring, done bool) error {
	if recurring {
		today := time.Now().Format("2006-01-02")
		if done {
			_, err := db.Exec("DELETE FROM completions WHERE todo_id=? AND date=?", id, today)
			return err
		}
		_, err := db.Exec("INSERT OR IGNORE INTO completions (todo_id,date) VALUES (?,?)", id, today)
		return err
	}
	if done {
		_, err := db.Exec("UPDATE todos SET done_at=NULL WHERE id=?", id)
		return err
	}
	_, err := db.Exec("UPDATE todos SET done_at=strftime('%Y-%m-%dT%H:%M:%SZ','now') WHERE id=?", id)
	return err
}

func IsTaskDoneToday(db *sql.DB, id int, recurring bool, doneAt sql.NullString) (bool, error) {
	if !recurring { return doneAt.Valid, nil }
	var n int
	err := db.QueryRow("SELECT COUNT(*) FROM completions WHERE todo_id=? AND date=?", id, time.Now().Format("2006-01-02")).Scan(&n)
	return n > 0, err
}
