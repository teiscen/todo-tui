package preview

import (
	"fmt"
	"strings"

	"todo-tui/backend"
	"todo-tui/leftbox"
	"todo-tui/taskform"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Preview displays details for the task currently selected in the left panel.
type Preview struct {
	Style    PreviewStyle
	Selected leftbox.Row
	HasTask  bool
	// Form is non-nil while the preview panel is displaying the task editor.
	Form *taskform.Form
}

// NewPreview creates the details panel with its default colors.
func NewPreview() Preview {
	return Preview{Style: GetDefaultStyle()}
}

// SetForm switches the preview panel into task-editing mode.
func (p *Preview) SetForm(form *taskform.Form) {
	p.Form = form
}

// SetTask switches back to read-only details; ok=false shows the empty state.
func (p *Preview) SetTask(row leftbox.Row, ok bool) {
	p.Form = nil
	p.Selected = row
	p.HasTask = ok
}

func (p Preview) Init() tea.Cmd { return nil }

// View renders task fields, then fits the result into the panel's fixed dimensions.
func (p Preview) View() tea.View {
	if p.Form != nil {
		return p.Form.View()
	}
	width := p.Style.Width
	height := p.Style.Height
	innerWidth := width - 4
	border := p.Style.BorderColor.ToStyle()

	var lines []string
	if !p.HasTask {
		lines = []string{
			border.Bold(true).Render("TASK DETAILS"),
			"",
			"  Select a task on the left",
			"  to inspect its details.",
		}
	} else {
		r := p.Selected
		lines = append(lines, border.Bold(true).Render("TASK DETAILS"), "")
		lines = append(lines, border.Render("TITLE"))
		lines = append(lines, wrapText(r.Title, innerWidth))
		lines = append(lines, "")

		lines = append(lines, border.Render("STATUS"))
		status := "Open"
		if r.Done {
			status = "Completed"
		}
		lines = append(lines, "  "+status, "")

		lines = append(lines, border.Render("LABEL"))
		label := "Unlabelled"
		if r.LabelID != "" {
			label = string(r.LabelID)
		}
		if r.Color != "" {
			label = r.Color.ToStyle().Bold(true).Render("  " + label)
		} else {
			label = "  " + label
		}
		lines = append(lines, label, "")

		lines = append(lines, border.Render("SCHEDULE / DUE DATE"))
		if r.Schedule != "" {
			lines = append(lines, "  "+r.Schedule)
		}
		if r.Due != "" {
			dueStyle := backend.HexCode("#939AB7").ToStyle()
			if r.DueTone == leftbox.ToneLate {
				dueStyle = backend.HexCode("#E78284").ToStyle()
			} else if r.DueTone == leftbox.ToneSoon {
				dueStyle = backend.HexCode("#E5C890").ToStyle()
			}
			lines = append(lines, dueStyle.Render("  "+r.Due))
		}
		if r.Schedule == "" && r.Due == "" {
			lines = append(lines, "  No schedule or due date")
		}
		lines = append(lines, "")

		if r.StepsTotal > 0 {
			lines = append(lines, border.Render("PROGRESS"))
			lines = append(lines, fmt.Sprintf("  %d of %d steps completed", r.StepsDone, r.StepsTotal))
			filled := barFill(r.StepsDone, r.StepsTotal, 20)
			lines = append(lines, "  "+backend.HexCode("#CA9EE6").ToStyle().Render(strings.Repeat("━", filled))+backend.HexCode("#939AB7").ToStyle().Render(strings.Repeat("─", 20-filled)))
			lines = append(lines, "")
		}

		lines = append(lines, border.Render("TASK ID"))
		lines = append(lines, fmt.Sprintf("  #%d", r.ID))
	}

	// Keep the panel height stable so the dashboard columns stay aligned.
	for i := range lines {
		lines[i] = lipgloss.NewStyle().Width(innerWidth).Render(lines[i])
	}
	for len(lines) > height-4 {
		lines = lines[:len(lines)-1]
	}
	for len(lines) < height-4 {
		lines = append(lines, "")
	}

	top := backend.GenBorder(" Task details ", border, backend.Top, width-2)
	box := lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Padding(0, 1).
		BorderForeground(lipgloss.Color(string(p.Style.BorderColor))).
		Border(lipgloss.RoundedBorder(), false, true, true, true)
	content := box.Render(strings.Join(lines, "\n"))
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, top, content))
}

// wrapText wraps a title at word boundaries while accounting for terminal cell width.
func wrapText(text string, width int) string {
	if text == "" {
		return "  (untitled)"
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return "  "
	}
	var lines []string
	line := "  "
	for _, word := range words {
		if lipgloss.Width(line)+1+lipgloss.Width(word) > width && line != "  " {
			lines = append(lines, line)
			line = "  " + word
		} else if line == "  " {
			line += word
		} else {
			line += " " + word
		}
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}

// barFill converts step completion into a bounded number of filled segments.
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

// Update is intentionally a no-op because the root model owns task selection.
func (p Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	return p, nil
}
