package main

import (
	"fmt"
	"os"

	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
)

func TestRendering() {
	state := backend.MockState()

	model := Model{}
	model.Initialize(state)
	model.p.Note.SetEntry()

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func main() {
	// testJSON()
	// tui.TestTUI()
	// testing.TestingMain()
	TestRendering()
}
