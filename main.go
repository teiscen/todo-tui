package main

import (
	"fmt"
	"os"

	"todo-tui/config"
	"todo-tui/database"
	"todo-tui/todo"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig("./artifacts/config.yaml")
	if err != nil {
		return err
	}

	db, err := database.ConnectDB("./artifacts/db.sqlite3")
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := database.CreateTableLabel(db); err != nil {
		return fmt.Errorf("create labels table: %w", err)
	}
	if _, err := database.CreateTableEntries(db); err != nil {
		return fmt.Errorf("create entries table: %w", err)
	}

	// Seed with the mock calendar until entries can be created in the app.
	n, err := database.LabelCount(db)
	if err != nil {
		return err
	}
	if n == 0 {
		if err := database.PopulateFromCalendar(db, cfg.Calendar); err != nil {
			return err
		}
	}

	_, err = tea.NewProgram(todo.NewTodo(cfg, db)).Run()
	return err
}
