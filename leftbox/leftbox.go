package leftbox

import (
	"fmt"
	"strings"

	"todo-tui/backend"
	"todo-tui/ui"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// Leftbox owns task-list selection, filtering, grouping, and row expansion state.
type Leftbox struct {
	Style  LeftboxStyle
	Labels backend.Labels

	Rows      []Row
	Tab       Tab
	Filter    backend.LabelID // "" = all labels
	Cursor    int             // index into the visible rows
	Collapsed map[Group]bool
	Expanded  map[int]bool
}

// New initializes an empty list and its per-group/per-task UI state.
func New(theme backend.Theme, labels backend.Labels, width, height int) Leftbox {
	style := NewStyle(theme)
	style.Width = width
	style.Height = height
	return Leftbox{Style: style, Labels: labels, Collapsed: map[Group]bool{}, Expanded: map[int]bool{}}
}

func (l Leftbox) Init() tea.Cmd {
	return nil
}

// ReloadState swaps in freshly loaded rows and keeps the same task selected
// (by ID) when it is still visible; otherwise the cursor is clamped into range.
func (l *Leftbox) ReloadState(rows []Row) {
	var keepID int
	keep := false
	if cur, ok := l.Selected(); ok {
		keepID, keep = cur.ID, true
	}
	l.Rows = rows
	if keep {
		for i, r := range l.selectableRows() {
			if r.ID == keepID {
				l.Cursor = i
				return
			}
		}
	}
	l.clampCursor()
}

// Selected returns the row at Cursor after applying the active tab and label filter.
func (l Leftbox) Selected() (Row, bool) {
	rows := l.selectableRows()
	if l.Cursor < 0 || l.Cursor >= len(rows) {
		return Row{}, false
	}
	return rows[l.Cursor], true
}

// selectableRows excludes tasks hidden inside collapsed groups.
func (l Leftbox) selectableRows() []Row {
	rows := visible(l.Rows, l.Tab, l.Filter)
	out := make([]Row, 0, len(rows))
	for _, row := range rows {
		if !l.Collapsed[row.Group] {
			out = append(out, row)
		}
	}
	return out
}

// ── layout ─────────────────────────────────────────────────────────────

// lineKind labels each body line so drawing and hover hit-testing consume the
// same structure and can never disagree.
type lineKind uint8

const (
	blankLine lineKind = iota
	headingLine
	itemLine
	detailLine
)

// bodyLine is one rendered line of the task body. Item lines carry the
// selectable index; the other kinds are not interactive.
type bodyLine struct {
	kind  lineKind
	text  string // rendered content, already styled
	index int    // selectable index for itemLine; -1 otherwise
}

// header returns the tab bar, the filter chips, and the trailing blank line
// that separates them from the task body.
func (l Leftbox) header(contentW int) []string {
	header := append([]string{l.tabBar()}, l.filterLines(contentW)...)
	return append(header, "")
}

// layout renders the grouped task body exactly once and returns it alongside
// the body-line index of the selected task. View and HoverAt both consume this
// single layout, so the keyboard cursor, drawing, and mouse hover always agree.
func (l Leftbox) layout(contentW int) ([]bodyLine, int) {
	vis := visible(l.Rows, l.Tab, l.Filter)
	body := make([]bodyLine, 0, len(vis))
	cursorLine := 0
	selectableIndex := 0
	prev := Group(-1)

	for _, r := range vis {
		if r.Group != prev {
			if prev != -1 {
				body = append(body, bodyLine{kind: blankLine})
			}
			body = append(body, bodyLine{kind: headingLine, text: l.groupHeading(r.Group, vis)})
			prev = r.Group
		}
		if l.Collapsed[r.Group] {
			continue
		}
		if selectableIndex == l.Cursor {
			cursorLine = len(body)
		}
		marker := "▸ "
		if l.Expanded[r.ID] {
			marker = "▾ "
		}
		body = append(body, bodyLine{
			kind:  itemLine,
			index: selectableIndex,
			text:  marker + l.renderRow(r, selectableIndex == l.Cursor, contentW-2),
		})
		selectableIndex++
		if l.Expanded[r.ID] {
			body = append(body, bodyLine{kind: detailLine, text: l.detailLine(r, contentW)})
		}
	}

	if len(body) == 0 {
		body = []bodyLine{{kind: blankLine, text: ui.Fg(l.Style.Theme.Dim).Render("nothing here")}}
	}
	return body, cursorLine
}

// detailLine is the expanded row's sub-line (progress or "No subtasks").
func (l Leftbox) detailLine(r Row, contentW int) string {
	detail := "  No subtasks"
	if r.StepsTotal > 0 {
		detail = fmt.Sprintf("  Progress: %d/%d steps completed", r.StepsDone, r.StepsTotal)
	}
	if r.Schedule != "" {
		detail += " · " + r.Schedule
	}
	return ui.Fg(l.Style.Theme.Dim).Render(ansi.Truncate(detail, contentW, "…"))
}

// View renders the tab/filter header, grouped task rows, and keyboard hints.
func (l Leftbox) View() string {
	s := l.Style
	contentW := s.Width - 4 // border + 1 column of padding on each side
	header := l.header(contentW)
	bodyH := s.Height - 2 - len(header)

	body, cursorLine := l.layout(contentW)
	start := windowStart(len(body), cursorLine, bodyH)
	body = body[start:min(len(body), start+bodyH)]

	lines := make([]string, 0, len(header)+len(body))
	lines = append(lines, header...)
	for _, line := range body {
		lines = append(lines, line.text)
	}

	top := ui.GenBorder(
		" Tasks ",
		ui.Fg(s.Theme.Border),
		ui.Top,
		s.Width-2,
	)
	box := lipgloss.NewStyle().
		Width(s.Width).
		Height(s.Height-1).
		Padding(0, 1).
		BorderForeground(lipgloss.Color(string(s.Theme.Border))).
		Border(lipgloss.RoundedBorder(), false, true, true, true)

	content := box.Render(strings.Join(lines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, top, content)
}

// HoverAt maps a terminal Y coordinate to a task row in the left panel, using
// the same layout the View draws.
func (l *Leftbox) HoverAt(y int) bool {
	contentW := l.Style.Width - 4
	header := l.header(contentW)
	bodyH := l.Style.Height - 2 - len(header)

	body, cursorLine := l.layout(contentW)
	start := windowStart(len(body), cursorLine, bodyH)
	body = body[start:min(len(body), start+bodyH)]

	index := y - (1 + len(header))
	if index < 0 || index >= len(body) {
		return false
	}
	line := body[index]
	if line.kind != itemLine || line.index == l.Cursor {
		return false
	}
	l.Cursor = line.index
	return true
}

// Update handles list navigation and emits commands for selection or completion.
func (l Leftbox) Update(msg tea.Msg) (Leftbox, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return l, nil
	}
	switch key.String() {
	case "j", "down":
		l.Cursor++
		l.clampCursor()
		return l, nil
	case "k", "up":
		l.Cursor--
		l.clampCursor()
		return l, nil
	case "]", "right_bracket", "close_bracket":
		l.Tab = (l.Tab + 1) % tabCount
		l.Cursor = 0
		return l, nil
	case "[", "left_bracket", "open_bracket":
		l.Tab = (l.Tab + tabCount - 1) % tabCount
		l.Cursor = 0
		return l, nil
	case "f":
		l.stepFilter()
		l.Cursor = 0
		return l, nil
	case " ", "space":
		return l, l.toggleCmd()
	case "enter":
		if r, ok := l.Selected(); ok {
			if l.Collapsed == nil {
				l.Collapsed = map[Group]bool{}
			}
			l.Collapsed[r.Group] = !l.Collapsed[r.Group]
			l.clampCursor()
		}
	case "e":
		if r, ok := l.Selected(); ok {
			if l.Expanded == nil {
				l.Expanded = map[int]bool{}
			}
			l.Expanded[r.ID] = !l.Expanded[r.ID]
		}
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
			parts = append(parts, ui.Fg(s.Theme.Accent).Bold(true).Render("["+label+"]"))
		} else {
			parts = append(parts, ui.Fg(s.Theme.Dim).Render(" "+label+" "))
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
		chips = append(chips, ui.Fg(color).Bold(active).Render(text))
		chipWidths = append(chipWidths, lipgloss.Width(text))
	}
	add("all", s.Theme.Text, l.Filter == "")
	for _, id := range l.Labels.Order {
		add(string(id), l.Labels.Info[id].Color, l.Filter == id)
	}
	var lines []string
	line := ""
	lineWidth := 0
	for i, chip := range chips {
		w := chipWidths[i]
		gap := 1
		if line == "" {
			gap = 0
		}
		if line != "" && lineWidth+gap+w > width {
			lines = append(lines, line)
			line, lineWidth = chip, w
			continue
		}
		if gap > 0 {
			line += " "
			lineWidth++
		}
		line += chip
		lineWidth += w
	}
	if line != "" {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

// groupHeading styles a group title using its semantic urgency.
func (l Leftbox) groupHeading(g Group, rows []Row) string {
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
	count := 0
	for _, row := range rows {
		if row.Group == g {
			count++
		}
	}
	marker := "▾"
	if l.Collapsed[g] {
		marker = "▸"
	}
	name = fmt.Sprintf("%s %s (%d)", marker, name, count)
	color := s.Theme.Dim
	if g == GroupOverdue {
		color = s.Theme.Late
	}
	return ui.Fg(color).Bold(true).Render(name)
}

// renderRow lays out one row as fixed columns:
//
//	▸ [ ] Title…            tag  ▰▰▱▱ 3/5   2d late
//
// renderRow aligns each task's title, label, progress/schedule, and due-date columns.
func (l Leftbox) renderRow(r Row, selected bool, width int) string {
	s := l.Style
	dim := ui.Fg(s.Theme.Dim)

	marker := " "
	if selected {
		marker = ui.Fg(s.Theme.Accent).Bold(true).Render(s.CursorChar)
	}
	check := s.CheckOpen
	if r.Done {
		check = s.CheckDone
	}

	// marker + space + checkbox + space, then 3 single-space gaps between columns.
	titleW := width - (1 + 1 + lipgloss.Width(check) + 1) - 3 - s.TagWidth - s.ProgWidth - s.DueWidth

	titleStyle := ui.Fg(s.Theme.Text)
	switch {
	case r.Done:
		titleStyle = dim.Strikethrough(true)
	case selected:
		titleStyle = titleStyle.Bold(true)
	}
	title := titleStyle.Render(ansi.Truncate(r.Title, titleW, "…"))

	tagStyle, accent := dim, s.Theme.Accent
	if r.Color != "" {
		tagStyle, accent = ui.Fg(r.Color), r.Color
	}
	tag := tagStyle.Render(ansi.Truncate(string(r.LabelID), s.TagWidth, ""))

	prog := ""
	switch {
	case r.StepsTotal > 0:
		n := ui.BarFill(r.StepsDone, r.StepsTotal, s.BarSegs)
		prog = ui.Fg(accent).Render(strings.Repeat(s.BarFull, n)) +
			dim.Render(strings.Repeat(s.BarEmpty, s.BarSegs-n)+
				fmt.Sprintf(" %d/%d", r.StepsDone, r.StepsTotal))
	case r.Schedule != "":
		prog = dim.Render(ansi.Truncate(r.Schedule, s.ProgWidth, ""))
	}

	due := ui.Fg(s.Theme.ToneColor(r.DueTone)).Render(r.Due)

	return strings.Join([]string{
		marker + " " + check + " " + ui.PadRight(title, titleW),
		ui.PadRight(tag, s.TagWidth),
		ui.PadRight(prog, s.ProgWidth),
		ui.PadLeft(due, s.DueWidth),
	}, " ")
}

// ── small helpers ───────────────────────────────────────────────────────

// clampCursor keeps selection within the currently visible row range.
func (l *Leftbox) clampCursor() {
	n := len(l.selectableRows())
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
