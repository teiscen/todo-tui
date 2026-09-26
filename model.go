package main

import (
	"fmt"
	"os"

	"todo-tui/backend"
	"todo-tui/calendar"

	tea "charm.land/bubbletea/v2"
)

func TestRendering() {
	state := backend.MockState()
	plan := calendar.InitialPlanner(
		state.GetCalendarRenderInfo(),
		state.GetNoteRenderInfo(),
	)

	p := tea.NewProgram(plan)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
