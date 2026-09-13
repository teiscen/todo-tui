package calender

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type calenderPage struct {
	color   ThemeColor
	entries []time.Time
}

type model struct {
	daySelected   int
	pageSelected  int
	calenderPages []calenderPage
}

// Need to populate the calenderPages from a config file
func initialModel() model {
	return model{
		daySelected:   0,
		pageSelected:  0,
		calenderPages: mockCalenderPages,
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
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var days string
	// Days of Week
	days += ColorRed.toStyle().Render("S ")
	days += ColorValid.toStyle().Render("M T W T F ")
	days += ColorRed.toStyle().Render("S")
	days += lipgloss.NewStyle().Render("\n")
	// Days from the previous months in the week
	count := 0
	for i := 0; i < int(getFirstDay()); i++ {
		days += ColorInvalid.toStyle().Render(smallCircle) + " "
		count++
	}
	// Days of the month
	for day := 1; day <= getNumOfDays(); day++ {
		// days += getStyleForDay(day, cursor int, )
		days += ColorValid.toStyle().Render(smallCircle) + " "
		count++
		if count%7 == 0 {
			days += "\n"
		}
	}
	// Days from the next month in the week
	for i := 0; i < int(getLastDay()); i++ {
		days += ColorInvalid.toStyle().Render(smallCircle) + " "
	}
	days += "\n"

	return tea.NewView(days)
}

func TestCalender() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
