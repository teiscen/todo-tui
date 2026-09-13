package calender

/*
*								  Sun = 0, Sat = 6
*	number of rows = ceil((days + start days)/7)
 */
import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

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
	invalidDay := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	validDay := lipgloss.NewStyle().Foreground(lipgloss.Color("#c0c0c0"))

	var days string

	count := 0
	for i := 0; i < int(getFirstDay()); i++ {
		days += invalidDay.Render("●") + " "
		count++
	}
	for day := 1; day <= getNumOfDays(); day++ {
		days += validDay.Render("●") + " "
		count++
		if count%7 == 0 {
			days += "\n"
		}
	}
	for i := 0; i < int(getLastDay()); i++ {
		days += invalidDay.Render("●") + " "
	}
	days += "\n"

	return tea.NewView(days)
}

func TestCalender() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
