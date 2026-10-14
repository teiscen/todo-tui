package weekly

import (
	"strings"

	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Task struct {
	Date    backend.Date
	Text    string
	LabelID backend.LabelID
}
type Day struct {
	Date  backend.Date
	Tasks []Task
}

type Weekly struct {
	Style   WeeklyStyle
	Labels  backend.Labels
	Tasks   []Day
	Overdue []Task
	TopLine int
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
	var internal strings.Builder
	// Extra blank line at top for spacing from border
	for i, f := range l.Tasks {
		if i > 0 {
			internal.WriteString("\n")
		}
		internal.WriteString(l.GenDay(f))
	}
	if len(l.Overdue) > 0 {
		internal.WriteString("\n\n Overdue \n")
		for _, t := range l.Overdue {
			internal.WriteString(l.GetOverdueLine(t))
			internal.WriteString("\n")
		}
	}

	lines := strings.Split(strings.TrimSuffix(internal.String(), "\n"), "\n")
	start := l.TopLine
	maxLines := 16
	end := min(start+maxLines, len(lines))

	visibleLines := lines[start:end]

	// If there's more content below, replace last visible line with centered 3 dots
	hasMore := end < len(lines)
	if hasMore && len(visibleLines) > 0 {
		indicator := lipgloss.NewStyle().
			Width(48).
			Align(lipgloss.Center).
			Render("● ● ●")
		visibleLines[len(visibleLines)-1] = indicator
	}

	// STRICT: only maxLines lines
	if len(visibleLines) > maxLines {
		visibleLines = visibleLines[:maxLines]
	}

	visible := strings.Join(visibleLines, "\n")

	top := backend.GenBorder(
		" Next 7 Days ",
		l.Style.ColorBorder.ToStyle(),
		backend.Top,
		48,
	)
	box := lipgloss.NewStyle().
		Width(50).
		Height(17).
		Border(
			lipgloss.RoundedBorder(),
			false, true, true, true,
		)

	content := box.Render(visible)
	combined := lipgloss.JoinVertical(lipgloss.Left, top, content)

	return tea.NewView(combined)
}

func (l Weekly) Update(msg tea.Msg) (Weekly, tea.Cmd) {
	var cmd tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "j", "down":
			l.TopLine = min(l.TopLine+1, max(0, len(l.Overdue)+len(l.Tasks)-13))
		case "k", "up":
			l.TopLine = max(l.TopLine-1, 0)
		}
	}
	return l, cmd
}
