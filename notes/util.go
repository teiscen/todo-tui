package notes

import backend "todo-tui/backend"

func (n *Notes) ChangeValue(str string) {
	n.TextArea.SetValue(str)
}

// SetAccent changes the accent color and regenerates the textarea styles.
func (n *Notes) SetAccent(c backend.HexCode) {
	n.Style.ColorAccent = c
	n.Style.Styles = n.Style.GenStyle()
}
