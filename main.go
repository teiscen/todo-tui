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

	if err := database.CreateTables(db); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	_, err = tea.NewProgram(todo.NewTodo(cfg, db)).Run()
	return err
}
