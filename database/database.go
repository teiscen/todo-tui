package database

import (
	"database/sql"
	"fmt"

	"todo-tui/backend"

	_ "github.com/glebarez/go-sqlite"
)

func FetchGridEntries(
	db *sql.DB, sDate, eDate backend.Date, labelID backend.LabelID,
) (map[string]string, error) {
	query := `
		SELECT date, msg FROM entries
		WHERE label_name = ? AND date BETWEEN ? AND ?
	`
	rows, err := db.Query(
		query,
		labelID, sDate.FormatSQL(), eDate.FormatSQL(),
	)
	if err != nil {
		return nil, fmt.Errorf("fetch grid: %w", err)
	}
	defer rows.Close()

	msgs := make(map[string]string)
	for rows.Next() {
		var d, msg string
		if err := rows.Scan(&d, &msg); err != nil {
			return nil, fmt.Errorf("fetch grid: %w", err)
		}
		date, err := backend.ParseSQL(d)
		if err != nil {
			return nil, fmt.Errorf("fetch grid %w", err)
		}
		msgs[date.Key()] = msg
	}
	return msgs, nil
}

// ConnectDB opens (or creates) the sqlite file at dbPath.
// Use ":memory:" for an in-memory database.
func ConnectDB(dbPath string) (*sql.DB, error) {
	// foreign_keys must be set per connection or ON DELETE CASCADE does nothing
	dsn := dbPath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

// Basic IO
func AddEntry(db *sql.DB, date, name, msg string) error {
	query := `INSERT INTO entries (date, label_name, msg) VALUES (?, ?, ?)
	          ON CONFLICT(date, label_name) DO UPDATE SET msg = excluded.msg;`
	_, err := db.Exec(query, date, name, msg)
	if err != nil {
		return fmt.Errorf("failed to add entry: %w", err)
	}
	return nil
}

func RemoveEntry(db *sql.DB, date, name string) error {
	query := `DELETE FROM entries WHERE date = ? AND label_name = ?;`
	_, err := db.Exec(query, date, name)
	if err != nil {
		return fmt.Errorf("failed to remove entry: %w", err)
	}
	return nil
}
