package preview

import "todo-tui/backend"

// PreviewStyle contains the preview panel's size and visual styling.
type PreviewStyle struct {
	Width int
	Height int
	BorderColor backend.HexCode
}

func GetDefaultStyle() PreviewStyle {
	return PreviewStyle{
		Width: 50,
		Height: 43,
		BorderColor: backend.HexCode("#C6D0F5"),
	}
}
