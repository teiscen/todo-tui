package calendar

import (
	"strings"
	"time"

	backend "todo-tui/backend_old"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Grid struct {
	Date     backend.Date
	IsValid  bool
	IsMarked bool
}

type Calendar struct {
	SelectedDate backend.Date
	Style        CalendarStyle
	Grid         [6][7]Grid
}

func NewCalendar() Calendar {
	return Calendar{
		SelectedDate: backend.TimeToDate(time.Now()),
		Style:        GetDefaultStyle(),
	}
}

func (c Calendar) Init() tea.Cmd {
	return nil
}

func (c Calendar) View() tea.View {
	var out string

	out = c.Style.RenderHeader() + strings.Repeat("\n", c.Style.RowSpacing+1)

	for row := range c.Grid {
		for col := range c.Grid[row] {
			cell := c.Grid[row][col]

			var style lipgloss.Style
			if cell.IsMarked {
				style = c.Style.AccentStyle(cell.IsValid)
			} else if col == 0 || col == 6 {
				style = c.Style.WeekendStyle(cell.IsValid)
			} else {
				style = c.Style.WeekdayStyle(cell.IsValid)
			}

			var icon string
			if cell.Date == c.SelectedDate {
				icon = c.Style.CharSelected
			} else {
				icon = c.Style.CharIcon
			}

			out += style.Render(icon)
		}
		// Need the default 1, so added +1 spacing
		out += strings.Repeat("\n", c.Style.RowSpacing+1)
	}

	return tea.NewView(lipgloss.NewStyle().Padding(1, 1).Render(out))
}

// hjkl - movements
// x    - mark/unmark
func (c Calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "h":
			c.Movement(Left)
		case "j":
			c.Movement(Down)
		case "k":
			c.Movement(Up)
		case "l":
			c.Movement(Right)
		case "x":
			c.ToggleMark()
		case "q", "ctrl+c":
			return c, tea.Quit
		}
	}
	return c, nil
}
