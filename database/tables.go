package database

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite" // registers the "sqlite" driver
)

// Open enables foreign keys on every pooled connection (SQLite has them off by default).
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", "file:"+path+"?_pragma=foreign_keys(1)")
}

// CreateTables runs every table creator in dependency order, then the indexes.
func CreateTables(db *sql.DB) error {
	steps := []struct {
		name string
		fn   func(*sql.DB) (sql.Result, error)
	}{
		{"labels", CreateTableLabels},
		{"todos", CreateTableTodos},
		{"steps", CreateTableSteps},
		{"step_notes", CreateTableStepNotes},
		{"recurrence", CreateTableRecurrence},
		{"completions", CreateTableCompletions},
		{"time_log", CreateTableTimeLog},
		{"links", CreateTableLinks},
		{"entries", CreateTableEntries},
		{"indexes", CreateIndexes},
	}
	for _, s := range steps {
		if _, err := s.fn(db); err != nil {
			return fmt.Errorf("create %s: %w", s.name, err)
		}
	}
	return nil
}

func CreateTableLabels(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS labels (
		name  TEXT PRIMARY KEY,
		color TEXT NOT NULL
	);`)
}

func CreateTableTodos(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		title      TEXT NOT NULL,
		label_id   TEXT REFERENCES labels(name) ON UPDATE CASCADE ON DELETE SET NULL,
		horizon    TEXT NOT NULL DEFAULT 'small' CHECK (horizon IN ('small', 'long')),
		due_date   TEXT,
		notes      TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
		done_at    TEXT
	);`)
}

func CreateTableSteps(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS steps (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id    INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
		pos        INTEGER NOT NULL,
		text       TEXT NOT NULL,
		status     TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'started', 'done')),
		started_at TEXT,
		done_at    TEXT,
		due_date   TEXT,
		UNIQUE (todo_id, pos)
	);`)
}

func CreateTableStepNotes(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS step_notes (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		step_id    INTEGER NOT NULL REFERENCES steps(id) ON DELETE CASCADE,
		body       TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
	);`)
}

func CreateTableRecurrence(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS recurrence (
		todo_id       INTEGER PRIMARY KEY REFERENCES todos(id) ON DELETE CASCADE,
		kind          TEXT NOT NULL CHECK (kind IN ('daily', 'weekdays', 'interval')),
		weekdays      INTEGER NOT NULL DEFAULT 0 CHECK (weekdays BETWEEN 0 AND 127),
		interval_days INTEGER CHECK (interval_days IS NULL OR interval_days > 0),
		start_date    TEXT NOT NULL,
		end_date      TEXT,
		CHECK (kind <> 'weekdays' OR weekdays > 0),
		CHECK (kind <> 'interval' OR interval_days IS NOT NULL)
	);`)
}

func CreateTableCompletions(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS completions (
		todo_id INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
		date    TEXT NOT NULL,
		PRIMARY KEY (todo_id, date)
	) WITHOUT ROWID;`)
}

func CreateTableTimeLog(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS time_log (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
		step_id INTEGER REFERENCES steps(id) ON DELETE SET NULL,
		date    TEXT NOT NULL,
		minutes INTEGER NOT NULL CHECK (minutes > 0)
	);`)
}

func CreateTableLinks(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS links (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
		label   TEXT NOT NULL,
		url     TEXT NOT NULL
	);`)
}

func CreateTableEntries(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE TABLE IF NOT EXISTS entries (
		date       TEXT NOT NULL,
		label_name TEXT NOT NULL REFERENCES labels(name) ON UPDATE CASCADE ON DELETE CASCADE,
		msg        TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (date, label_name)
	);`)
}

func CreateIndexes(db *sql.DB) (sql.Result, error) {
	return db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_todos_open_due ON todos(due_date) WHERE done_at IS NULL;
	CREATE INDEX IF NOT EXISTS idx_todos_label    ON todos(label_id);
	CREATE INDEX IF NOT EXISTS idx_steps_todo     ON steps(todo_id, pos);
	CREATE INDEX IF NOT EXISTS idx_steps_due      ON steps(due_date) WHERE status <> 'done';
	CREATE INDEX IF NOT EXISTS idx_time_todo_date ON time_log(todo_id, date);`)
}
