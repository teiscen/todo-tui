package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestBarFill(t *testing.T) {
	cases := []struct {
		done, total, segs int
		want              int
	}{
		{0, 0, 4, 0}, // zero total
		{0, 5, 4, 0}, // nothing done
		{1, 5, 4, 1},
		{3, 5, 4, 2},
		{4, 5, 4, 3},
		{5, 5, 4, 4},    // fully done
		{1, 40, 4, 1},   // tiny progress floors at 1
		{2, 6, 4, 1},    // rounds to nearest segment
		{1, 100, 4, 1},  // tiny progress floors at 1
		{99, 100, 4, 3}, // near-complete caps below segs
	}
	for _, c := range cases {
		if got := BarFill(c.done, c.total, c.segs); got != c.want {
			t.Errorf("BarFill(%d,%d,%d) = %d, want %d", c.done, c.total, c.segs, got, c.want)
		}
	}
}

func TestGenBorderWidth(t *testing.T) {
	plain := lipgloss.NewStyle()
	for _, text := range []string{" Tasks ", "[small 4]  long-term 2   recurring 3", "x", ""} {
		for _, posn := range []Posn{Top, Mid, Bot} {
			if got := lipgloss.Width(GenBorder(text, plain, posn, 48)); got != 50 {
				t.Errorf("GenBorder(%q) width = %d, want 50", text, got)
			}
		}
	}
}

func TestPad(t *testing.T) {
	if got := PadRight("ab", 5); got != "ab   " {
		t.Errorf("PadRight(ab, 5) = %q, want %q", got, "ab   ")
	}
	if got := PadLeft("ab", 5); got != "   ab" {
		t.Errorf("PadLeft(ab, 5) = %q, want %q", got, "   ab")
	}
	if got := PadRight("abcd", 3); got != "abcd" {
		t.Errorf("PadRight(abcd, 3) = %q, want %q", got, "abcd")
	}
	if got := PadLeft("abcd", 3); got != "abcd" {
		t.Errorf("PadLeft(abcd, 3) = %q, want %q", got, "abcd")
	}
}

func TestWrap(t *testing.T) {
	got := WrapJoin([]string{"aa", "bb", "cc", "dd"}, " · ", 9)
	if len(got) != 2 || got[0] != "aa · bb" || got[1] != "cc · dd" {
		t.Errorf("WrapJoin = %q", got)
	}
	lines := WrapText("one two three four", 12, 2)
	if strings.Join(lines, "|") != "  one two|  three four" {
		t.Errorf("WrapText = %q", lines)
	}
}
