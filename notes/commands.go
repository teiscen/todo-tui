package notes

func (n *Notes) ToggleFocus() {
	if (*n).TextArea.Focused() {
		(*n).TextArea.Blur()
	} else {
		(*n).TextArea.Focus()
	}
}
