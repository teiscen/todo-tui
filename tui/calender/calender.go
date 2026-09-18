package calender

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	cursor   CalenderDate
	page     Label
	style    calenderStyle
	calender Calender
	grid     RenderGrid
}

// Need to populate the calenderPages from a config file
func initialModel() model {
	m := model{
		cursor:   timeToCalenderDate(time.Now()),
		page:     0,
		style:    doubleCalenderStyle{},
		calender: mockCalender(),
	}
	rg := m.CreateGrid()
	m.grid = m.UpdateSelected(m.calender, m.cursor, rg, m.page)
	return m
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
	var color ThemeColor
	if _, ok := m.calender.Months[m.cursor][m.page]; ok {
		color = m.page.info().color
	} else {
		color = ColorValid
	}

	complete := genBorder(
		m.Render(), ColorValid,
		calenderDateToTime(m.cursor).Format("January 02"), color, true,
		m.page.info().name, m.page.info().color, true,
	)

	return tea.NewView(complete)
}

func TestCalender() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
