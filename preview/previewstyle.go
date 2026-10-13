package preview

import "todo-tui/backend"

type PreviewStyle struct {
	BorderColor backend.HexCode
}

func GetDefaultStyle() PreviewStyle {
	return PreviewStyle{
		BorderColor: backend.HexCode("#c6d0h5"),
	}
}
