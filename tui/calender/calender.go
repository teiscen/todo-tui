package calender

import (
	"fmt"
	"os"
	"strconv"
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
		daySelected:  time.Now().Day() + int(getFirstDay()),
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
		// Quitting out
		case "ctrl+c", "q", "shift-c":
			return m, tea.Quit
		// Changing Calender Pages
		case "h", "left":
			index := (m.pageSelected - 1 + len(m.calenderPages)) % len(m.calenderPages)
			m.pageSelected = index
		case "l", "right":
			index := (m.pageSelected + 1) % len(m.calenderPages)
			m.pageSelected = index
		// Movement
		case "w":
			m.daySelected = move(m.daySelected, Up, getTotalNumOfDays())
		case "a":
			m.daySelected = move(m.daySelected, Left, getTotalNumOfDays())
		case "s":
			m.daySelected = move(m.daySelected, Down, getTotalNumOfDays())
		case "d":
			m.daySelected = move(m.daySelected, Right, getTotalNumOfDays())
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	// Just generate the whole basic grid first
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

	// Replace the cells that have an entry to display color
	indexToDate := func(idx int) int {
		switch {
		case m.daySelected < int(getFirstDay()):
			daysInPrevMonth := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.Local).Day()
			return daysInPrevMonth - m.daySelected
		case m.daySelected > int(getFirstDay())+getNumOfDays():
			return m.daySelected - (int(getFirstDay()) + getNumOfDays())
		default:
			return m.daySelected - int(getFirstDay())
		}
	}
	entryToIndex := func(entry time.Time) int {
		return int(getFirstDay()) + entry.Day()
	}
	for _, e := range m.calenderPages[m.pageSelected].entries {
		idx := entryToIndex(e)
		color := m.calenderPages[m.pageSelected].color
		cells[idx] = m.style.cellRenderer(color, true)
	}

	// Replace the selected cell with the correct render
	var selectedColor ThemeColor
	var month int // -1 prev , 0 curr, +1  next
	if m.daySelected < int(getFirstDay()) {
		month = -1
		selectedColor = ColorInvalid
	} else if m.daySelected > int(getFirstDay())+getNumOfDays() {
		month = 1
		selectedColor = ColorInvalid
	} else {
		month = 0
		selectedColor = ColorValid
		for _, e := range m.calenderPages[m.pageSelected].entries {
			idx := entryToIndex(e)
			if idx == m.daySelected {
				selectedColor = m.calenderPages[m.pageSelected].color
			}
		}
	}
	cells[m.daySelected] = m.style.selectedRenderer(selectedColor)

	var topBorderStr string
	switch month {
	case -1:
		topBorderStr = (time.Now().Month() - 1).String()
	case 0:
		topBorderStr = time.Now().Month().String()
	case 1:
		topBorderStr = (time.Now().Month() + 1).String()
	}

	// Put it all together
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
		topBorderStr+" "+strconv.Itoa(indexToDate(m.daySelected)), ColorValid, true,
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
