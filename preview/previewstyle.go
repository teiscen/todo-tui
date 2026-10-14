package preview

import "todo-tui/backend"

type PreviewStyle struct {
	BorderColor backend.HexCode
}

func GetDefaultStyle() PreviewStyle {
	return PreviewStyle{
		BorderColor: backend.HexCode("#C6D0F5"),
	}
}
