package preview

import "todo-tui/backend"

// PreviewStyle contains the preview panel's size and visual styling.
type PreviewStyle struct {
	Width  int
	Height int
	Theme  backend.Theme
}

// NewStyle builds a preview style from the shared theme and dimensions.
func NewStyle(theme backend.Theme, width, height int) PreviewStyle {
	return PreviewStyle{
		Width:  width,
		Height: height,
		Theme:  theme,
	}
}
