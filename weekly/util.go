package weekly

import (
	"fmt"
	"strings"
	"time"

	"todo-tui/backend"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func FormatDate(d backend.Date) string {
	return d.GetTime().Format("Mon 02")
}

func (w Weekly) GetOverdueLine(task Task) string {
	// Box inner width = 48 (50 - 2 borders)
	// Layout: [date:8 right][1 space][text:32 left][label:5 right] = 46, padded to 48
	today := backend.NewDate(time.Now())
	daysDiff := task.Date.DaysUntil(today)

	dateStr := fmt.Sprintf("%dd late", daysDiff)

	text := w.Style.ApplyColorText(task.Text)
	textLen := 35
	truncatedText := ansi.Truncate(text, textLen-3, "...")
	lab := w.Style.ApplyColorLabel(task.LabelID, &w)
	truncatedLab := ansi.Truncate(lab, 3, "")

	dateS := lipgloss.NewStyle().Width(8).MaxWidth(8).Align(lipgloss.Right)
	textS := lipgloss.NewStyle().Width(textLen).MaxWidth(textLen).Align(lipgloss.Left).Padding(0, 1)
	labS := lipgloss.NewStyle().Width(3).MaxWidth(5).Align(lipgloss.Right)

	// No left padding - flush against border
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dateS.Render(dateStr),
		lipgloss.NewStyle().Width(1).Render(" "),
		textS.Render(truncatedText),
		labS.Render(truncatedLab),
	)
}

func (w Weekly) GenDay(t Day) string {
	var date string
	if t.Date.IsWeekend() {
		date = w.Style.ApplyColorWeekend(FormatDate(t.Date))
	} else {
		date = w.Style.ApplyColorWeekday(FormatDate(t.Date))
	}

	if len(t.Tasks) == 0 {
		sep := w.Style.ApplyColorText(w.Style.EmptySep)
		return GetDayLine(date, sep, "", "")
	}

	var lines []string
	sep := w.Style.ApplyColorText(w.Style.DefaultSep)
	for i, task := range t.Tasks {
		text := w.Style.ApplyColorText(task.Text)
		lab := w.Style.ApplyColorLabel(task.LabelID, &w)

		if i == 0 {
			lines = append(lines, GetDayLine(date, sep, text, lab))
		} else {
			lines = append(lines, GetDayLine("", sep, text, lab))
		}
	}

	return strings.Join(lines, "\n")
}

// Need to get the text in preStyled
func GetDayLine(date, sep, text, lab string) string {
	// Layout: [date:8 right][1 space][sep:1][1 space][text:31][label:5] = 48
	textLen := 30
	truncatedText := ansi.Truncate(text, textLen-3, "...")
	truncatedLab := ansi.Truncate(lab, 3, "")

	dateS := lipgloss.NewStyle().Width(8).MaxWidth(8).Align(lipgloss.Right).Padding(0, 1)
	sepS := lipgloss.NewStyle().Width(4).Align(lipgloss.Center).Padding(0, 1)
	textS := lipgloss.NewStyle().Width(textLen).MaxWidth(textLen).Align(lipgloss.Left)
	labS := lipgloss.NewStyle().Width(5).MaxWidth(5).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dateS.Render(date),
		sepS.Render(sep),
		textS.Render(truncatedText),
		labS.Render(truncatedLab),
	)
}
