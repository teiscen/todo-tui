package calender

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cursor   CalenderDate
	style    calenderStyle
	calender Calender
}

// Need to populate the calenderPages from a config file
func initialModel() model {
	return model{
		cursor:   timeToCalenderDate(time.Now()),
		style:    doubleCalenderStyle{},
		calender: mockCalender(),
		// renderCache: nil{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Need to add moving around months
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Quitting out
		case "ctrl+c", "q", "shift-c":
			return m, tea.Quit
		// Changing Calender Pages
		case "p", "left":
			// update the grid
		case "n", "right":
			// update the grid
		// Movement
		case "k":
		case "h":
		case "l":
		case "j":
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var out string

	return tea.NewView(out)
}

func TestCalender() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
