package main

import (
	"fmt"
	"os"

	"todo-tui/planner"

	tea "charm.land/bubbletea/v2"
)

// func TestRendering() {
// 	state := backend.MockState()
//
// 	model := Model{}
// 	model.Initialize(state)
// 	model.p.Note.SetEntry()
//
// 	p := tea.NewProgram(model)
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }

func main() {
	// testJSON()
	// tui.TestTUI()
	// testing.TestingMain()
	// TestRendering()

	// calendar.CalendarTestPrint()
	// calendar.CalendarTestUpdate()

	// cal := calendar.TestCalendar()
	// p := tea.NewProgram(cal)

	// note := notes.TestNotes()
	// p := tea.NewProgram(note)

	planner := planner.TestPlanner()
	p := tea.NewProgram(planner)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// notes.NotesTestPrint()
	// notes.NotesTestToggle()
}
