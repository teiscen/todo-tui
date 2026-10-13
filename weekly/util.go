package weekly

import (
	"fmt"
	"time"

	"todo-tui/backend"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func FormatDate(d backend.Date) string {
	return d.GetTime().Format("Mon 02")
}

func (w Weekly) GenOverdue(task task) string {
	//  Width = 50 - 2
	//  2d late [8 - Width - 5 ] uni
	// 99d late
	today := backend.NewDate(time.Now())
	daysDiff := task.date.DaysUntil(today)

	var dateStr string
	dateStr = fmt.Sprintf("%dd late", daysDiff)

	text := w.Style.ApplyColorText(task.text)
	lab := w.Style.ApplyColorLabel(task.labelID, &w)

	dateS := lipgloss.NewStyle().Width(8).Align(lipgloss.Right)
	textS := lipgloss.NewStyle().Width(34).Align(lipgloss.Left)
	labS := lipgloss.NewStyle().Width(5).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dateS.Render(dateStr),
		textS.Render(text),
		labS.Render(lab),
	)
}

func (w Weekly) GenDay(t Foo) string {
	var date string
	if t.date.IsWeekend() {
		date = w.Style.ApplyColorWeekend(FormatDate(t.date))
	} else {
		date = w.Style.ApplyColorWeekday(FormatDate(t.date))
	}

	if len(t.tasks) == 0 {
		sep := w.Style.ApplyColorText(w.Style.EmptySep)
		return GetLine(date, sep, "", "")
	}

	var out string
	sep := w.Style.ApplyColorText(w.Style.DefaultSep)
	for i, task := range t.tasks {
		text := w.Style.ApplyColorText(task.text)
		lab := w.Style.ApplyColorLabel(task.labelID, &w)

		if i == 0 {
			out += GetLine(date, sep, text, lab)
		} else {
			out += GetLine("", sep, text, lab)
		}
		out += "\n"
	}

	return out
}

// Need to get the text in preStyled
func GetLine(date, sep, text, lab string) string {
	// Date = 6
	// Seperator = 3
	// LabelID = 3
	// Total width = 50-2(padding)
	// text Name = Total Width - 6 - 3 - 5 = 34
	textLen := 34
	truncatedText := ansi.Truncate(text, textLen, "...")
	truncatedLab := ansi.Truncate(lab, 3, "")
	dateS := lipgloss.NewStyle().Width(4).Align(lipgloss.Left).Padding(0, 1)
	sepS := lipgloss.NewStyle().Width(3).Align(lipgloss.Center).Padding(0, 1)
	textS := lipgloss.NewStyle().Width(textLen).Align(lipgloss.Left).Padding(0, 1)
	labS := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Padding(0, 1)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dateS.Render(date),
		sepS.Render(sep),
		textS.Render(truncatedText),
		labS.Render(truncatedLab),
	)
}
