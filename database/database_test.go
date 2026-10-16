package database

import (
	"database/sql"
	"testing"

	"todo-tui/backend"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := ConnectDB(":memory:")
	if err != nil {
		t.Fatalf("ConnectDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestCreateTablesIdempotent(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("first CreateTables: %v", err)
	}
	// Safe to run again on an existing schema.
	if err := CreateTables(db); err != nil {
		t.Fatalf("second CreateTables: %v", err)
	}
}

func TestSyncLabels(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	labels := backend.Labels{
		Order: []backend.LabelID{"uni", "work"},
		Info: map[backend.LabelID]backend.LabelInfo{
			"uni":  {Color: "#81c8be"},
			"work": {Color: "#ca9ee6"},
		},
	}
	if err := SyncLabels(db, labels); err != nil {
		t.Fatalf("SyncLabels: %v", err)
	}

	// Upsert: change a color, drop one label from the config, and add another.
	labels.Order = []backend.LabelID{"uni", "gym"}
	labels.Info["uni"] = backend.LabelInfo{Color: "#ffffff"}
	labels.Info["gym"] = backend.LabelInfo{Color: "#ef9f76"}
	if err := SyncLabels(db, labels); err != nil {
		t.Fatalf("second SyncLabels: %v", err)
	}

	got := map[string]string{}
	rows, err := db.Query("SELECT name, color FROM labels ORDER BY name")
	if err != nil {
		t.Fatalf("query labels: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, color string
		if err := rows.Scan(&name, &color); err != nil {
			t.Fatalf("scan label: %v", err)
		}
		got[name] = color
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate labels: %v", err)
	}

	want := map[string]string{
		"uni":  "#ffffff", // updated in place
		"work": "#ca9ee6", // no longer in config, kept
		"gym":  "#ef9f76", // added
	}
	if len(got) != len(want) {
		t.Fatalf("label count = %d, want %d: %v", len(got), len(want), got)
	}
	for name, color := range want {
		if got[name] != color {
			t.Errorf("label %q color = %q, want %q", name, got[name], color)
		}
	}
}

func TestEntriesRoundTrip(t *testing.T) {
	db := openTestDB(t)
	if err := CreateTables(db); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}
	labels := backend.Labels{
		Order: []backend.LabelID{"uni"},
		Info:  map[backend.LabelID]backend.LabelInfo{"uni": {Color: "#81c8be"}},
	}
	if err := SyncLabels(db, labels); err != nil {
		t.Fatalf("SyncLabels: %v", err)
	}

	sDate, err := backend.ParseSQL("2026-10-01")
	if err != nil {
		t.Fatalf("ParseSQL start: %v", err)
	}
	eDate, err := backend.ParseSQL("2026-10-31")
	if err != nil {
		t.Fatalf("ParseSQL end: %v", err)
	}

	if err := AddEntry(db, "2026-10-15", "uni", "study"); err != nil {
		t.Fatalf("AddEntry: %v", err)
	}
	// Same day + label upserts the message instead of creating a second row.
	if err := AddEntry(db, "2026-10-15", "uni", "revised"); err != nil {
		t.Fatalf("AddEntry upsert: %v", err)
	}
	if err := AddEntry(db, "2026-10-16", "uni", "review"); err != nil {
		t.Fatalf("AddEntry second: %v", err)
	}

	msgs, err := FetchGridEntries(db, sDate, eDate, "uni")
	if err != nil {
		t.Fatalf("FetchGridEntries: %v", err)
	}
	if msgs["2026-10-15"] != "revised" {
		t.Errorf("2026-10-15 = %q, want %q", msgs["2026-10-15"], "revised")
	}
	if msgs["2026-10-16"] != "review" {
		t.Errorf("2026-10-16 = %q, want %q", msgs["2026-10-16"], "review")
	}
	if _, ok := msgs["2026-10-17"]; ok {
		t.Errorf("unexpected entry for 2026-10-17: %v", msgs)
	}

	if err := RemoveEntry(db, "2026-10-15", "uni"); err != nil {
		t.Fatalf("RemoveEntry: %v", err)
	}
	msgs, err = FetchGridEntries(db, sDate, eDate, "uni")
	if err != nil {
		t.Fatalf("FetchGridEntries after remove: %v", err)
	}
	if _, ok := msgs["2026-10-15"]; ok {
		t.Errorf("2026-10-15 still present after remove: %v", msgs)
	}
	if msgs["2026-10-16"] != "review" {
		t.Errorf("2026-10-16 = %q, want %q after remove", msgs["2026-10-16"], "review")
	}
}
