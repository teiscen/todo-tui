package leftbox

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Leftbox struct{}

func NewLeftbox() Leftbox {
	return Leftbox{}
}

func (l Leftbox) Init() tea.Cmd {
	return nil
}

func (l Leftbox) View() tea.View {
	box := lipgloss.NewStyle().
		Width(50).
		Height(43).
		Border(lipgloss.RoundedBorder())

	content := box.Render("HEHE")

	return tea.NewView(content)
}

func (l Leftbox) Update(msg tea.Msg) (Leftbox, tea.Cmd) {
	var cmd tea.Cmd
	return l, cmd
}
