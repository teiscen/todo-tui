package calendar

import (
	"todo-tui/backend"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// Might need to track window size to adjust render
type Planner struct {
	cal  Calendar
	note Note
}

func InitialPlanner(
	cRI backend.CalendarRenderInfo,
	nRI backend.NoteRenderInfo,
) Planner {
	c := Calendar{
		cRI,
		doubleStyle,
	}

	n := getNote()
	return Planner{c, n}
}

func (p Planner) Init() tea.Cmd {
	// p.cal = c.cal.Init()
	// p.note = p.note.Init()
	return tea.Batch(textarea.Blink, tea.RequestBackgroundColor)
}

func (p Planner) View() tea.View {
	calRender := p.cal.Render()
	noteRender := p.note.Render()
	planRender := p.Render(calRender, noteRender)

	return tea.NewView(planRender)
}

func (p Planner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// m.textarea.SetStyles(textarea.DefaultStyles(msg.IsDark()))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "shift-c":
			return p, tea.Quit

		case "tab":
			// m.focus = m.focus.toggleFocus()
			// if m.focus == FocusCal {
			// 	m.notes.Blur()
			// } else {
			// 	cmd = m.notes.Focus()
			// 	cmds = append(cmds, cmd)
			// }
		}
	}

	// if m.focus == FocusNote {
	// 	m.notes, cmd = m.notes.Update(msg)
	// 	cmds = append(cmds, cmd)
	// }

	// m = m.reconcile()

	return p, tea.Batch(cmds...)
}
