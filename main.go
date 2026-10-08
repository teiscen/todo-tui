package main

//* TODO:
// Clean up code by identifying commonly re-used portion of code and making them methods
// Clean up the end portion of Planner's Update func

import (
	"fmt"
	"os"

	"todo-tui/todo"

	tea "charm.land/bubbletea/v2"
)

func main() {
	configPath := "./artifacts/config.yaml"
	// config, _ := config.LoadConfig(configPath)
	// calendar.CalendarTestPrint()
	// calendar.CalendarTestUpdate()

	// cal := calendar.TestCalendar()
	// p := tea.NewProgram(cal)

	// note := notes.TestNotes()
	// p := tea.NewProgram(note)

	// planner := planner.TestPlanner()
	// planner := config.LoadPlannerModel()
	// p := tea.NewProgram(planner)

	todo := todo.NewTodo(configPath, nil)
	p := tea.NewProgram(todo)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// notes.NotesTestPrint()
	// notes.NotesTestToggle()
}
