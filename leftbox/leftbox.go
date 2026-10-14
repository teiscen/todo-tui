package leftbox

import (
	"fmt"
	"strings"

	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Leftbox struct {
	Style  LeftboxStyle
	Labels backend.Labels

	Rows   []Row
	Tab    Tab
	Filter backend.LabelID // "" = all labels
	Cursor int             // index into the visible rows
}

func NewLeftbox() Leftbox {
	return Leftbox{Style: GetDefaultStyle()}
}

func (l Leftbox) Init() tea.Cmd {
	return nil
}

// ReloadState swaps in freshly loaded rows and keeps the cursor on the screen.
func (l *Leftbox) ReloadState(rows []Row) {
	l.Rows = rows
	l.clampCursor()
}

func (l Leftbox) Selected() (Row, bool) {
	vis := visible(l.Rows, l.Tab, l.Filter)
	if l.Cursor < 0 || l.Cursor >= len(vis) {
		return Row{}, false
	}
	return vis[l.Cursor], true
}

func (l Leftbox) View() tea.View {
	s := l.Style
	contentW := s.Width - 4   // border + 1 column of padding on each side
	bodyH := s.Height - 2 - 2 // top line, bottom border, then filter bar + blank line

	vis := visible(l.Rows, l.Tab, l.Filter)

	var body []string
	cursorLine := 0
	prev := Group(-1)
	for i, r := range vis {
		if r.Group != prev {
			if prev != -1 {
				body = append(body, "")
			}
			body = append(body, l.groupHeading(r.Group))
			prev = r.Group
		}
		if i == l.Cursor {
			cursorLine = len(body)
		}
		body = append(body, l.renderRow(r, i == l.Cursor, contentW))
	}
	if len(body) == 0 {
		body = []string{s.ColorDim.ToStyle().Render("nothing here")}
	}
	start := windowStart(len(body), cursorLine, bodyH)
	body = body[start:min(len(body), start+bodyH)]

	lines := append([]string{l.filterBar(contentW), ""}, body...)

	top := backend.GenBorder(
		l.tabBar(),
		s.ColorBorder.ToStyle(),
		backend.Top,
		s.Width-2,
	)
	box := lipgloss.NewStyle().
		Width(s.Width).
		Height(s.Height-1).
		Padding(0, 1).
		BorderForeground(lipgloss.Color(string(s.ColorBorder))).
		Border(lipgloss.RoundedBorder(), false, true, true, true)

	content := box.Render(strings.Join(lines, "\n"))
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, top, content))
}

func (l Leftbox) Update(msg tea.Msg) (Leftbox, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return l, nil
	}
	switch key.String() {
	case "j", "down":
		l.Cursor++
		l.clampCursor()
		return l, l.selectCmd()
	case "k", "up":
		l.Cursor--
		l.clampCursor()
		return l, l.selectCmd()
	case "]":
		l.Tab = (l.Tab + 1) % tabCount
		l.Cursor = 0
		return l, l.selectCmd()
	case "[":
		l.Tab = (l.Tab + tabCount - 1) % tabCount
		l.Cursor = 0
		return l, l.selectCmd()
	case "f":
		l.stepFilter()
		l.Cursor = 0
		return l, l.selectCmd()
	case " ", "space":
		return l, l.toggleCmd()
	}
	return l, nil
}

// ── drawing ─────────────────────────────────────────────────────────────

// tabBar is the title that sits in the top border. Every tab is the same width
// whether active ([x]) or not ( x ), so the bar doesn't shift when you switch.
func (l Leftbox) tabBar() string {
	s := l.Style
	counts := openCounts(l.Rows, l.Filter)
	parts := make([]string, 0, tabCount)
	for i := range tabCount {
		label := fmt.Sprintf("%s %d", tabNames[i], counts[i])
		if Tab(i) == l.Tab {
			parts = append(parts, s.ColorAccent.ToStyle().Bold(true).Render("["+label+"]"))
		} else {
			parts = append(parts, s.ColorDim.ToStyle().Render(" "+label+" "))
		}
	}
	return strings.Join(parts, " ")
}

// filterBar: [all] uni work gym ..., each label in its own color.
func (l Leftbox) filterBar(width int) string {
	s := l.Style
	item := func(name string, color backend.HexCode, active bool) string {
		if active {
			return color.ToStyle().Bold(true).Render("[" + name + "]")
		}
		return color.ToStyle().Render(" " + name + " ")
	}
	parts := []string{item("all", s.ColorText, l.Filter == "")}
	for _, id := range l.Labels.Order {
		parts = append(parts, item(string(id), l.Labels.Info[id].Color, l.Filter == id))
	}
	return ansi.Truncate(strings.Join(parts, ""), width, "…")
}

func (l Leftbox) groupHeading(g Group) string {
	s := l.Style
	name := map[Group]string{
		GroupOverdue: "OVERDUE",
		GroupToday:   "TODAY",
		GroupWeek:    "THIS WEEK",
		GroupLater:   "LATER",
		GroupNoDate:  "NO DATE",
		GroupDone:    "DONE",
	}[g]
	if l.Tab == TabRecurring && g == GroupLater {
		name = "OTHER DAYS"
	}
	color := s.ColorDim
	if g == GroupOverdue {
		color = s.ColorLate
	}
	return color.ToStyle().Bold(true).Render(name)
}

// renderRow lays out one row as fixed columns:
//
//	▸ [ ] Title…            tag  ▰▰▱▱ 3/5   2d late
func (l Leftbox) renderRow(r Row, selected bool, width int) string {
	s := l.Style
	dim := s.ColorDim.ToStyle()

	marker := " "
	if selected {
		marker = s.ColorAccent.ToStyle().Bold(true).Render(s.CursorChar)
	}
	check := s.CheckOpen
	if r.Done {
		check = s.CheckDone
	}

	// marker + space + checkbox + space, then 3 single-space gaps between columns.
	titleW := width - (1 + 1 + lipgloss.Width(check) + 1) - 3 - s.TagWidth - s.ProgWidth - s.DueWidth

	titleStyle := s.ColorText.ToStyle()
	switch {
	case r.Done:
		titleStyle = dim.Strikethrough(true)
	case selected:
		titleStyle = titleStyle.Bold(true)
	}
	title := titleStyle.Render(ansi.Truncate(r.Title, titleW, "…"))

	tagStyle, accent := dim, s.ColorAccent
	if r.Color != "" {
		tagStyle, accent = r.Color.ToStyle(), r.Color
	}
	tag := tagStyle.Render(ansi.Truncate(string(r.LabelID), s.TagWidth, ""))

	prog := ""
	switch {
	case r.StepsTotal > 0:
		n := barFill(r.StepsDone, r.StepsTotal, s.BarSegs)
		prog = accent.ToStyle().Render(strings.Repeat(s.BarFull, n)) +
			dim.Render(strings.Repeat(s.BarEmpty, s.BarSegs-n)+
				fmt.Sprintf(" %d/%d", r.StepsDone, r.StepsTotal))
	case r.Schedule != "":
		prog = dim.Render(ansi.Truncate(r.Schedule, s.ProgWidth, ""))
	}

	due := s.ToneColor(r.DueTone).ToStyle().Render(r.Due)

	return strings.Join([]string{
		marker + " " + check + " " + padRight(title, titleW),
		padRight(tag, s.TagWidth),
		padRight(prog, s.ProgWidth),
		padLeft(due, s.DueWidth),
	}, " ")
}

// ── small helpers ───────────────────────────────────────────────────────

func padRight(s string, w int) string {
	return s + strings.Repeat(" ", max(0, w-lipgloss.Width(s)))
}

func padLeft(s string, w int) string {
	return strings.Repeat(" ", max(0, w-lipgloss.Width(s))) + s
}

func (l *Leftbox) clampCursor() {
	n := len(visible(l.Rows, l.Tab, l.Filter))
	l.Cursor = max(0, min(l.Cursor, n-1))
}

// stepFilter cycles all -> each label in order -> all.
func (l *Leftbox) stepFilter() {
	order := append([]backend.LabelID{""}, l.Labels.Order...)
	for i, id := range order {
		if id == l.Filter {
			l.Filter = order[(i+1)%len(order)]
			return
		}
	}
	l.Filter = ""
}
