package database

import (
	"database/sql"
	"fmt"

	"todo-tui/backend"

	_ "github.com/glebarez/go-sqlite"
)

// t.db.FetchRange(start ackend.Date, lID backend.LabelID)
func FetchGrid(db *sql.DB, sDate backend.Date, lID backend.LabelID) (map[backend.Date]backend.Entry, error) {
	// the 6x7 grid starts on the Sunday on or before the 1st of the month
	start := sDate.FirstDay()
	start = start.AddDate(0, 0, -int(start.Weekday()))
	end := start.AddDate(0, 0, 41)

	rows, err := db.Query(`
		SELECT date, msg FROM entries
		WHERE label_name = ? AND date BETWEEN ? AND ?`,
		lID, start.FormatSQL(), end.FormatSQL())
	if err != nil {
		return nil, fmt.Errorf("fetch grid: %w", err)
	}
	defer rows.Close()

	msgs := make(map[string]string)
	for rows.Next() {
		var d, msg string
		if err := rows.Scan(&d, &msg); err != nil {
			return nil, fmt.Errorf("fetch grid scan: %w", err)
		}
		msgs[d] = msg
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetch grid rows: %w", err)
	}

	// key by backend.Date: walk the 42 cells and look each one up
	out := make(map[backend.Date]backend.Entry, len(msgs))
	d := start
	for i := 0; i < 42; i++ {
		if msg, ok := msgs[d.FormatSQL()]; ok {
			out[d] = backend.Entry{Status: backend.Full, Msg: msg}
		}
		d = d.AddDate(0, 0, 1)
	}
	return out, nil
}

func ConnectDB(dbPath string) (*sql.DB, error) {
	// create a connection to the sqlite db, or create one if it doesnt exist
	// to connect to SQlit db in memory replace dbPath with :memory:
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("issue connected to db")
	}

	// fmt.Println("Connected to db succesfully")

	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("issue connected to db")
	}
	return db, nil
}

func CreateTableLabel(db *sql.DB) (sql.Result, error) {
	query := `CREATE TABLE IF NOT EXISTS labels (
		name TEXT PRIMARY KEY,
		color TEXT NOT NULL
	);`
	return db.Exec(query)
}

func CreateTableEntries(db *sql.DB) (sql.Result, error) {
	query := `CREATE TABLE IF NOT EXISTS entries (
		date TEXT NOT NULL,
		label_name TEXT NOT NULL REFERENCES labels(name) ON DELETE CASCADE,
		msg TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (date, label_name)
	);`
	return db.Exec(query)
}

// TODO:
// Phase out the backend types (entry and label)
// to use these instead

type Label struct {
	Name  string
	Color backend.HexCode
}

type Entry struct {
	Date backend.Date
	Name string // Matches label name
	Msg  string
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

func GetEntry(db *sql.DB, date, name string) (*Entry, error) {
	query := `SELECT date, label_name, msg FROM entries WHERE date = ? AND label_name = ?;`
	var e Entry
	err := db.QueryRow(query, date, name).Scan(&e.Date, &e.Name, &e.Msg)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil if no entry found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}
	return &e, nil
}

func UpdateEntry(db *sql.DB, date, name, msg string) error {
	query := `
		UPDATE entries
		SET msg = ?
		WHERE date = ? AND label_name = ?;
	`

	result, err := db.Exec(query, msg, date, name)
	if err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no entry found for date %s and label %s", date, name)
	}

	return nil
}

func AddLabel(db *sql.DB, name, color string) error {
	query := `INSERT INTO labels (name, color) VALUES (?, ?)
	          ON CONFLICT(name) DO UPDATE SET color = excluded.color;`
	_, err := db.Exec(query, name, color)
	if err != nil {
		return fmt.Errorf("failed to add label: %w", err)
	}
	return nil
}

func RemoveLabel(db *sql.DB, name string) error {
	// ON DELETE CASCADE will automatically clear referencing entries
	query := `DELETE FROM labels WHERE name = ?;`
	_, err := db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to remove label: %w", err)
	}
	return nil
}

// MORE complicated things
type Something map[string]string

func GetMonth(db *sql.DB, name, startDate, endDate string) (Something, error) {
	query := `
		SELECT date, msg
		FROM entries
		WHERE label_name = ? AND date >= ? AND date <= ?;
	`
	rows, err := db.Query(query, name, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query date range: %w", err)
	}
	defer rows.Close()

	result := make(Something)
	for rows.Next() {
		var msg, date string
		if err := rows.Scan(&date, &msg); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		result[date] = msg
	}

	return result, nil
}
