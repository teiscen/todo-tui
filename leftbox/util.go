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

type Tone int

const (
	ToneNormal Tone = iota
	ToneSoon
	ToneLate
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
	DueTone Tone
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
func DueInfo(due, today backend.Date) (string, Tone) {
	days := today.DaysUntil(due)
	switch {
	case days < 0:
		return fmt.Sprintf("%dd late", -days), ToneLate
	case days == 0:
		return "today", ToneSoon
	case days == 1:
		return "tmrw", ToneSoon
	default:
		return fmt.Sprintf("in %dd", days), ToneNormal
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

// barFill is how many of segs segments are filled. It never shows an empty bar
// once something is done, or a full one before everything is.
func barFill(done, total, segs int) int {
	if total <= 0 {
		return 0
	}
	n := (done*segs + total/2) / total
	if done > 0 && n == 0 {
		n = 1
	}
	if done < total && n == segs {
		n = segs - 1
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
