package leftbox

import (
	"testing"
	"time"

	"todo-tui/backend"
)

func d(y int, m time.Month, day int) backend.Date {
	return backend.NewDate(time.Date(y, m, day, 0, 0, 0, 0, time.UTC))
}

func TestGroupFor(t *testing.T) {
	today := d(2024, time.January, 15)
	cases := []struct {
		name string
		due  *backend.Date
		done bool
		want Group
	}{
		{"done", nil, true, GroupDone},
		{"no date", nil, false, GroupNoDate},
		{"overdue", datePtr(d(2024, time.January, 14)), false, GroupOverdue},
		{"today", datePtr(today), false, GroupToday},
		{"one day out", datePtr(d(2024, time.January, 16)), false, GroupWeek},
		{"six days out", datePtr(d(2024, time.January, 21)), false, GroupWeek},
		{"seven days out", datePtr(d(2024, time.January, 22)), false, GroupLater},
	}
	for _, c := range cases {
		if got := GroupFor(c.due, c.done, today); got != c.want {
			t.Errorf("%s: GroupFor() = %d, want %d", c.name, got, c.want)
		}
	}
}

func datePtr(d backend.Date) *backend.Date { return &d }

func TestDueInfo(t *testing.T) {
	today := d(2024, time.January, 15)
	cases := []struct {
		due      backend.Date
		wantText string
		wantTone backend.Tone
	}{
		{d(2024, time.January, 14), "1d late", backend.ToneLate},
		{d(2024, time.January, 10), "5d late", backend.ToneLate},
		{today, "today", backend.ToneSoon},
		{d(2024, time.January, 16), "tmrw", backend.ToneSoon},
		{d(2024, time.January, 20), "in 5d", backend.ToneNormal},
	}
	for _, c := range cases {
		text, tone := DueInfo(c.due, today)
		if text != c.wantText || tone != c.wantTone {
			t.Errorf("DueInfo(%s) = (%q, %d), want (%q, %d)",
				c.due.FormatSQL(), text, tone, c.wantText, c.wantTone)
		}
	}
}

func TestScheduleText(t *testing.T) {
	cases := []struct {
		mask int
		want string
	}{
		{0x7F, "daily"}, // all seven days
		{0, ""},
		{42, "M/W/F"},       // bits 1,3,5
		{62, "M/Tu/W/Th/F"}, // bits 1..5
		{1, "Su"},           // bit 0
		{127, "daily"},      // all seven via decimal
	}
	for _, c := range cases {
		if got := ScheduleText(c.mask); got != c.want {
			t.Errorf("ScheduleText(%d) = %q, want %q", c.mask, got, c.want)
		}
	}
}

func TestVisible(t *testing.T) {
	rows := []Row{
		{ID: 1, Tab: TabSmall, Group: GroupLater, LabelID: "uni"},
		{ID: 2, Tab: TabSmall, Group: GroupToday, LabelID: "work"},
		{ID: 3, Tab: TabLong, Group: GroupToday, LabelID: "uni"},
		{ID: 4, Tab: TabSmall, Group: GroupOverdue, LabelID: "gym"},
	}

	// Tab filter excludes ID 3; results are stable-sorted by group ascending.
	got := visible(rows, TabSmall, "")
	want := []int{4, 2, 1}
	if len(got) != len(want) {
		t.Fatalf("visible() returned %d rows, want %d", len(got), len(want))
	}
	for i, id := range want {
		if got[i].ID != id {
			t.Errorf("visible()[%d].ID = %d, want %d", i, got[i].ID, id)
		}
	}

	// Label filter keeps only the uni row.
	got = visible(rows, TabSmall, "uni")
	if len(got) != 1 || got[0].ID != 1 {
		t.Errorf("visible(TabSmall, uni) = %v, want [1]", rowIDs(got))
	}
}

func rowIDs(rows []Row) []int {
	ids := make([]int, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	return ids
}

func TestOpenCounts(t *testing.T) {
	rows := []Row{
		{ID: 1, Tab: TabSmall, LabelID: "uni", Done: false},
		{ID: 2, Tab: TabSmall, LabelID: "work", Done: true},
		{ID: 3, Tab: TabLong, LabelID: "uni", Done: false},
		{ID: 4, Tab: TabRecurring, LabelID: "work", Done: false},
	}

	got := openCounts(rows, "")
	want := [tabCount]int{1, 1, 1}
	if got != want {
		t.Errorf("openCounts(no filter) = %v, want %v", got, want)
	}

	got = openCounts(rows, "uni")
	want = [tabCount]int{1, 1, 0}
	if got != want {
		t.Errorf("openCounts(uni) = %v, want %v", got, want)
	}
}

func TestWindowStart(t *testing.T) {
	cases := []struct {
		total, cursor, height int
		want                  int
	}{
		{5, 0, 10, 0},   // fits, no scroll
		{20, 0, 5, 0},   // at top
		{20, 10, 5, 8},  // centered-ish
		{20, 19, 5, 15}, // clamped to end
		{20, 3, 5, 1},
	}
	for _, c := range cases {
		if got := windowStart(c.total, c.cursor, c.height); got != c.want {
			t.Errorf("windowStart(%d,%d,%d) = %d, want %d", c.total, c.cursor, c.height, got, c.want)
		}
	}
}
