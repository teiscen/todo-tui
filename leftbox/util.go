package leftbox

import (
	"fmt"
	"sort"
	"strings"

	"todo-tui/backend"
)

// Pure helpers: no lipgloss / bubbletea in here, so they are easy to test.

type Tab int

const (
	TabSmall Tab = iota
	TabLong
	TabRecurring
	tabCount
)

var tabNames = [tabCount]string{"small", "long-term", "recurring"}

type Group int

const (
	GroupOverdue Group = iota
	GroupToday
	GroupWeek
	GroupLater
	GroupNoDate
	GroupDone
)

// Row is one todo, fully resolved by whoever loads it (colors, group, due text).
// The view only draws it.
type Row struct {
	ID    int
	Tab   Tab
	Group Group

	Title   string
	LabelID backend.LabelID
	Color   backend.HexCode // label color, "" = no label

	Done bool

	StepsDone  int // progress bar shows when StepsTotal > 0
	StepsTotal int
	Schedule   string // recurring rows: "M/W/F", "daily", "every 3d"

	Due     string // "2d late", "today", "in 6d" ("" = nothing shown)
	DueTone backend.Tone
}

// headingColor picks the foreground color for a group heading. Unselected
// headings are muted; a selected heading pops in its own color — overdue in the
// late color, every other group in the normal text color.
func headingColor(g Group, selected bool, theme backend.Theme) backend.HexCode {
	if !selected {
		return theme.Dim
	}
	if g == GroupOverdue {
		return theme.Late
	}
	return theme.Text
}

// GroupFor decides which heading a one-time todo sits under.
// due == nil means the todo has no due date.
func GroupFor(due *backend.Date, done bool, today backend.Date) Group {
	switch {
	case done:
		return GroupDone
	case due == nil:
		return GroupNoDate
	}
	switch days := today.DaysUntil(*due); {
	case days < 0:
		return GroupOverdue
	case days == 0:
		return GroupToday
	case days <= 6:
		return GroupWeek
	default:
		return GroupLater
	}
}

// DueInfo is the text for the right-hand column plus how to color it.
func DueInfo(due, today backend.Date) (string, backend.Tone) {
	days := today.DaysUntil(due)
	switch {
	case days < 0:
		return fmt.Sprintf("%dd late", -days), backend.ToneLate
	case days == 0:
		return "today", backend.ToneSoon
	case days == 1:
		return "tmrw", backend.ToneSoon
	default:
		return fmt.Sprintf("in %dd", days), backend.ToneNormal
	}
}

// ScheduleText turns a recurrence weekday bitmask (bit0 = Sunday) into "M/W/F".
// Interval rules are better written by the caller: fmt.Sprintf("every %dd", n).
func ScheduleText(mask int) string {
	if mask&0x7F == 0x7F {
		return "daily"
	}
	short := [7]string{"Su", "M", "Tu", "W", "Th", "F", "Sa"}
	var out []string
	for i, s := range short {
		if mask&(1<<i) != 0 {
			out = append(out, s)
		}
	}
	return strings.Join(out, "/")
}

// visible keeps the rows for this tab (and label filter), ordered by group.
// The sort is stable, so the caller's order holds inside each group.
func visible(rows []Row, tab Tab, filter backend.LabelID) []Row {
	var out []Row
	for _, r := range rows {
		if r.Tab != tab || (filter != "" && r.LabelID != filter) {
			continue
		}
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Group < out[j].Group })
	return out
}

// groupCounts tallies visible rows per group for the heading "(n)" labels.
// It counts every row in the tab, including those hidden by a collapsed group.
func groupCounts(rows []Row, tab Tab, filter backend.LabelID) map[Group]int {
	counts := map[Group]int{}
	for _, r := range rows {
		if r.Tab != tab || (filter != "" && r.LabelID != filter) {
			continue
		}
		counts[r.Group]++
	}
	return counts
}

// openCounts is the number shown on each tab: rows not done yet.
func openCounts(rows []Row, filter backend.LabelID) [tabCount]int {
	var n [tabCount]int
	for _, r := range rows {
		if r.Done || (filter != "" && r.LabelID != filter) {
			continue
		}
		n[r.Tab]++
	}
	return n
}

// windowStart is the first line to show when body is taller than height,
// keeping the cursor line roughly centered.
func windowStart(total, cursorLine, height int) int {
	if total <= height {
		return 0
	}
	return max(0, min(cursorLine-height/2, total-height))
}
