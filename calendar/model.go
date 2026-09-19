package calendar

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	selDate  CalendarDate
	selLabel int
	labels   []Label
	style    CalendarStyle
	calendar Calendar
	rInfo    RenderInfo
	config   Config
}

// Need to populate the calenderPages from a config file
func initialModel() Model {
	config := defaultConfig()

	labels := make([]Label, 0, len(config.Labels))
	for label := range config.Labels {
		labels = append(labels, label)
	}

	m := Model{
		timeToCalendarDate(time.Now()),
		0,
		labels,
		doubleStyle,
		mockCalendar(config.Labels),
		RenderInfo{},
		config,
	}

	m.rInfo = CreateGrid(m)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Quitting
		case "ctrl+c", "q", "shift-c":
			return m, tea.Quit
		// Changing Calendar Pages
		case "p", "left":
			m.movePage(Left)
		case "n", "right":
			m.movePage(Right)
		// Movement
		case "k":
			m.moveCursor(Up)
			m.rInfo = CreateGrid(m)
			return m, tea.ClearScreen
		case "h":
			m.moveCursor(Left)
			m.rInfo = CreateGrid(m)
			return m, tea.ClearScreen
		case "l":
			m.moveCursor(Right)
			m.rInfo = CreateGrid(m)
			return m, tea.ClearScreen
		case "j":
			m.moveCursor(Down)
			m.rInfo = CreateGrid(m)
			return m, tea.ClearScreen
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	label := m.config.Labels[m.labels[m.selLabel]]
	labelColor := m.config.Colors[label.ColorName]

	complete := genBorder(
		m.Render(),
		m.cursorColor(),
		calendarDateToTime(m.selDate).Format("January 02"),
		m.cursorColor(),
		label.Name,
		labelColor,
	)

	return tea.NewView(complete)
}

func TestCalendar() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
