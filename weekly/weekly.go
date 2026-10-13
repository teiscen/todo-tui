package weekly

import (
	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type task struct {
	date    backend.Date
	text    string
	labelID backend.LabelID
}
type Foo struct {
	date  backend.Date
	tasks []task
}

type Weekly struct {
	Style   WeeklyStyle
	Labels  backend.Labels
	Tasks   []Foo
	Overdue []task
}

func NewWeekly() Weekly {
	return Weekly{
		Style: GetDefaultStyle(),
	}
}

func (l Weekly) Init() tea.Cmd {
	return nil
}

func (l Weekly) View() tea.View {
	top := backend.GenBorder(
		"Next 7 Days",
		l.Style.ColorBorder.ToStyle(),
		backend.Top,
		50,
	)

	box := lipgloss.NewStyle().
		Width(50).
		Height(17).
		Border(
			lipgloss.RoundedBorder(),
			false, true, true, true,
		)

	content := box.Render("Weekly")
	combined := lipgloss.JoinVertical(lipgloss.Left, top, content)

	return tea.NewView(combined)
}

func (l Weekly) Update(msg tea.Msg) (Weekly, tea.Cmd) {
	var cmd tea.Cmd
	return l, cmd
}
