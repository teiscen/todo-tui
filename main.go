package main

import (
	"todo-tui/calendar"
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
	calendar.CalendarTestUpdate()
}
