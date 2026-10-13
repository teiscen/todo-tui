package preview

import (
	"todo-tui/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Preview struct {
	Style PreviewStyle
}

func NewPreview() Preview {
	return Preview{
		Style: GetDefaultStyle(),
	}
}

func (p Preview) Init() tea.Cmd {
	return nil
}

func (p Preview) View() tea.View {
	box := lipgloss.NewStyle().
		Width(50).
		Height(24).
		Border(
			lipgloss.RoundedBorder(),
			false, true, true, true,
		)

	top := backend.GenBorder(
		" Preview ",
		p.Style.BorderColor.ToStyle(),
		backend.Top,
		48,
	)

	content := box.Render("Weekly")
	combined := lipgloss.JoinVertical(lipgloss.Left, top, content)

	return tea.NewView(combined)
}

func (l Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	var cmd tea.Cmd
	return l, cmd
}
