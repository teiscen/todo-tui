package calendar

import (
	"strings"

	backend "todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Info struct {
	Date     backend.Date
	IsMarked bool
}

type Calendar struct {
	SelectedDate backend.Date
	Style        CalendarStyle
	State        [6][7]Info
}

func NewCalendar(theme backend.Theme, date backend.Date) Calendar {
	return Calendar{
		SelectedDate: date,
		Style:        NewStyle(theme),
	}
}

func (c Calendar) Init() tea.Cmd {
	return nil
}

func (c Calendar) View() string {
	// Append the header
	out := c.Style.RenderHeader()
	// Loop through the state grid and begin rendering
	for row := range c.State {
		// Add the spacing between the rows
		out += strings.Repeat("\n", c.Style.RowSpacing+1)
		for col := range c.State[row] {
			// Determine if its the current month (Lowers saturation)
			cell := c.State[row][col]
			isValid := cell.Date.Month() == c.SelectedDate.Month()

			var style lipgloss.Style
			if cell.IsMarked {
				style = c.Style.AccentStyle(isValid)
			} else if col == 0 || col == 6 {
				style = c.Style.WeekendStyle(isValid)
			} else {
				style = c.Style.WeekdayStyle(isValid)
			}

			var icon string
			if cell.Date.Equal(c.SelectedDate) {
				icon = c.Style.CharSelected
			} else {
				icon = c.Style.CharIcon
			}

			out += style.Render(icon)
		}
	}

	cal := lipgloss.NewStyle().Padding(1, 0).Render(out)
	return cal
}

// hjkl - movements
func (c Calendar) Update(msg tea.Msg) (Calendar, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "h":
			cmd = c.Movement(Left)
		case "j":
			cmd = c.Movement(Down)
		case "k":
			cmd = c.Movement(Up)
		case "l":
			cmd = c.Movement(Right)
		case "x":
			cmd = c.ToggleMark()
		case "q", "ctrl+c":
			return c, tea.Quit
		}
	}

	return c, cmd
}
