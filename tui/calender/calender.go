package calender

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type calenderPage struct {
	color   ThemeColor
	entries []time.Time
}

type model struct {
	daySelected   int
	pageSelected  int
	style         calenderStyle
	calenderPages []calenderPage
}

// Need to populate the calenderPages from a config file
func initialModel() model {
	return model{
		daySelected:  0,
		pageSelected: 0,
		// style:        singleCalenderStyle{},
		style:         doubleCalenderStyle{},
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
		case "ctrl+c", "q", "shift-c":
			return m, tea.Quit
		case "h", "left":
			index := (m.pageSelected - 1 + len(m.calenderPages)) % len(m.calenderPages)
			m.pageSelected = index
		case "l", "right":
			index := (m.pageSelected + 1) % len(m.calenderPages)
			m.pageSelected = index
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var cells []string
	for i := 0; i < int(getFirstDay()); i++ {
		cells = append(cells, m.style.invalidRenderer(ColorInvalid))
	}
	for i := 0; i < int(getNumOfDays()); i++ {
		cells = append(cells, m.style.cellRenderer(ColorValid, false))
	}
	for i := 0; i < int(getLastDay()); i++ {
		cells = append(cells, m.style.invalidRenderer(ColorInvalid))
	}

	entryToIndex := func(entry time.Time) int {
		return int(getFirstDay()) + entry.Day()
	}
	for _, e := range m.calenderPages[m.pageSelected].entries {
		color := m.calenderPages[m.pageSelected].color
		cells[entryToIndex(e)] = m.style.cellRenderer(color, true)
	}

	var weeks []string
	for i := 0; i < len(cells); i += 7 {
		end := i + 7
		if end > len(cells) {
			end = len(cells)
		}
		// Join takes variadic args (...string), not a slice so the ... unpacks it into individual args
		weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Top, cells[i:end]...))
	}

	header := m.style.weekdayRenderer()
	sep := "\n" + strings.Repeat("\n", m.style.getRowSpacing())
	weeksBlock := strings.Join(weeks, sep)
	grid := header + "\n" + weeksBlock + "\n"

	complete := genBorder(grid, ColorValid,
		"Sept 10", ColorValid, true,
		"Label", m.calenderPages[m.pageSelected].color, true)

	return tea.NewView(complete)
}

func TestCalender() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
