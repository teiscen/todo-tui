package preview

import (
	"sort"
	"strings"

	"todo-tui/backend"
	"todo-tui/leftbox"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Preview struct {
	Style   PreviewStyle
	Rows    []leftbox.Row
	Entries map[string]string
}

func NewPreview() Preview {
	return Preview{Style: GetDefaultStyle()}
}

func (p *Preview) SetRows(rows []leftbox.Row) {
	p.Rows = rows
}

func (p *Preview) SetEntries(entries map[string]string) {
	p.Entries = entries
}

func (p Preview) Init() tea.Cmd { return nil }

func (p Preview) View() tea.View {
	const width = 50
	const height = 43
	innerWidth := width - 4

	var lines []string
	lines = append(lines, p.Style.BorderColor.ToStyle().Bold(true).Render("UPCOMING"))
	lines = append(lines, "")

	count := 0
	for _, row := range p.Rows {
		if row.Done || (row.Group != leftbox.GroupToday && row.Group != leftbox.GroupWeek) {
			continue
		}
		title := ansi.Truncate(row.Title, 27, "…")
		label := strings.TrimSpace(string(row.LabelID))
		due := row.Due
		line := "  ○ " + title
		if due != "" {
			line = ansi.Truncate(line, 29, "…")
			line += "  " + due
		}
		if label != "" {
			line += "  " + label
		}
		lines = append(lines, ansi.Truncate(line, innerWidth, "…"))
		count++
	}
	if count == 0 {
		lines = append(lines, "  No upcoming tasks")
	}

	lines = append(lines, "", p.Style.BorderColor.ToStyle().Bold(true).Render("CALENDAR ENTRIES"), "")
	keys := make([]string, 0, len(p.Entries))
	for key := range p.Entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	entryCount := 0
	for _, key := range keys {
		msg := strings.TrimSpace(p.Entries[key])
		if msg == "" {
			continue
		}
		lines = append(lines, "  "+ansi.Truncate(key+"  "+strings.Split(msg, "\n")[0], innerWidth-2, "…"))
		entryCount++
		if entryCount == 10 {
			break
		}
	}
	if entryCount == 0 {
		lines = append(lines, "  No saved entries in this month")
	}

	for len(lines) > height-4 {
		lines = lines[:len(lines)-1]
	}
	for len(lines) < height-4 {
		lines = append(lines, "")
	}

	top := backend.GenBorder(" Next 7 days ", p.Style.BorderColor.ToStyle(), backend.Top, width-2)
	box := lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Padding(0, 1).
		BorderForeground(lipgloss.Color(string(p.Style.BorderColor))).
		Border(lipgloss.RoundedBorder(), false, true, true, true)
	content := box.Render(strings.Join(lines, "\n"))
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, top, content))
}

func (p Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	return p, nil
}
