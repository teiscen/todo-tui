package preview

import (
	"fmt"
	"strings"

	"todo-tui/backend"
	"todo-tui/taskform"
	"todo-tui/ui"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Detail is the fully-resolved task data the preview pane draws. It holds only
// what the preview needs and has no dependency on leftbox; the root translates
// the selected row into one.
type Detail struct {
	ID         int
	Title      string
	Done       bool
	LabelID    backend.LabelID
	Color      backend.HexCode
	Schedule   string
	Due        string
	DueTone    backend.Tone
	StepsDone  int
	StepsTotal int
}

// Preview displays details for the task currently selected in the left panel.
type Preview struct {
	Style    PreviewStyle
	Selected Detail
	HasTask  bool
	// Form is non-nil while the preview panel is displaying the task editor.
	Form *taskform.Form
}

// New creates the details panel with its configured colors and size.
func New(theme backend.Theme, width, height int) Preview {
	return Preview{Style: NewStyle(theme, width, height)}
}

// SetForm switches the preview panel into task-editing mode.
func (p *Preview) SetForm(form *taskform.Form) {
	p.Form = form
}

// SetTask switches back to read-only details; ok=false shows the empty state.
func (p *Preview) SetTask(detail Detail, ok bool) {
	p.Form = nil
	p.Selected = detail
	p.HasTask = ok
}

func (p Preview) Init() tea.Cmd { return nil }

// View renders task fields, then fits the result into the panel's fixed dimensions.
func (p Preview) View() string {
	if p.Form != nil {
		return p.Form.View()
	}
	width := p.Style.Width
	height := p.Style.Height
	innerWidth := width - 4
	border := ui.Fg(p.Style.Theme.Border)

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
		if r.Title == "" {
			lines = append(lines, "  (untitled)")
		} else if wrapped := ui.WrapText(r.Title, innerWidth, 2); len(wrapped) == 0 {
			lines = append(lines, "  ")
		} else {
			lines = append(lines, strings.Join(wrapped, "\n"))
		}
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
			label = ui.Fg(r.Color).Bold(true).Render("  " + label)
		} else {
			label = "  " + label
		}
		lines = append(lines, label, "")

		lines = append(lines, border.Render("SCHEDULE / DUE DATE"))
		if r.Schedule != "" {
			lines = append(lines, "  "+r.Schedule)
		}
		if r.Due != "" {
			dueStyle := ui.Fg(p.Style.Theme.Dim)
			if r.DueTone == backend.ToneLate {
				dueStyle = ui.Fg(p.Style.Theme.Late)
			} else if r.DueTone == backend.ToneSoon {
				dueStyle = ui.Fg(p.Style.Theme.Soon)
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
			filled := ui.BarFill(r.StepsDone, r.StepsTotal, 20)
			lines = append(lines, "  "+ui.Fg(p.Style.Theme.Accent).Render(strings.Repeat("━", filled))+ui.Fg(p.Style.Theme.Dim).Render(strings.Repeat("─", 20-filled)))
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

	top := ui.GenBorder(" Task details ", border, ui.Top, width-2)
	box := lipgloss.NewStyle().
		Width(width).
		Height(height-1).
		Padding(0, 1).
		BorderForeground(lipgloss.Color(string(p.Style.Theme.Border))).
		Border(lipgloss.RoundedBorder(), false, true, true, true)
	content := box.Render(strings.Join(lines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, top, content)
}

// Update is intentionally a no-op because the root model owns task selection.
func (p Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	return p, nil
}
