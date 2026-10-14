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
	header := append([]string{l.tabBar()}, l.filterLines(contentW)...)
	header = append(header, "")
	bodyH := s.Height - 2 - len(header)

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

	lines := append(header, body...)

	top := backend.GenBorder(
		" Tasks ",
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

// HoverAt maps a terminal Y coordinate to a task row in the left panel.
func (l *Leftbox) HoverAt(y int) bool {
	contentW := l.Style.Width - 4
	header := append([]string{l.tabBar()}, l.filterLines(contentW)...)
	header = append(header, "")
	bodyH := l.Style.Height - 2 - len(header)
	vis := visible(l.Rows, l.Tab, l.Filter)

	var body []int
	prev := Group(-1)
	cursorLine := 0
	for i, r := range vis {
		if r.Group != prev {
			if prev != -1 {
				body = append(body, -1)
			}
			body = append(body, -1)
			prev = r.Group
		}
		if i == l.Cursor {
			cursorLine = len(body)
		}
		body = append(body, i)
	}
	start := windowStart(len(body), cursorLine, bodyH)
	body = body[start:min(len(body), start+bodyH)]
	index := y - (1 + len(header))
	if index < 0 || index >= len(body) || body[index] < 0 || body[index] == l.Cursor {
		return false
	}
	l.Cursor = body[index]
	return true
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

// filterLines renders colored label chips and wraps them onto as many lines as needed.
func (l Leftbox) filterLines(width int) []string {
	s := l.Style
	chips := []string{}
	chipWidths := []int{}
	add := func(name string, color backend.HexCode, active bool) {
		text := " " + name + " "
		if active {
			text = "[" + name + "]"
		}
		chips = append(chips, color.ToStyle().Bold(active).Render(text))
		chipWidths = append(chipWidths, lipgloss.Width(text))
	}
	add("all", s.ColorText, l.Filter == "")
	for _, id := range l.Labels.Order {
		add(string(id), l.Labels.Info[id].Color, l.Filter == id)
	}
	var lines []string
	line := ""
	lineWidth := 0
	for i, chip := range chips {
		w := chipWidths[i]
		gap := 1
		if line == "" { gap = 0 }
		if line != "" && lineWidth+gap+w > width {
			lines = append(lines, line)
			line, lineWidth = chip, w
			continue
		}
		if gap > 0 { line += " "; lineWidth++ }
		line += chip
		lineWidth += w
	}
	if line != "" { lines = append(lines, line) }
	if len(lines) == 0 { return []string{""} }
	return lines
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
