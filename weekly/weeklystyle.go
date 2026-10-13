package weekly

import "todo-tui/backend"

type WeeklyStyle struct {
	DefaultSep string
	EmptySep   string

	ColorText backend.HexCode

	ColorBorder  backend.HexCode
	ColorWeekday backend.HexCode
	ColorWeekend backend.HexCode
}

func GetDefaultStyle() WeeklyStyle {
	return WeeklyStyle{
		DefaultSep: "○",
		EmptySep:   "—",

		ColorBorder:  backend.HexCode("#c6d0h5"),
		ColorWeekday: backend.HexCode("#C6D0F5"),
		ColorWeekend: backend.HexCode("#939AB7"),
	}
}

func (w WeeklyStyle) ApplyColorText(text string) string {
	return w.ColorText.ToStyle().Render(text)
}

func (w WeeklyStyle) ApplyColorBorder(text string) string {
	return w.ColorBorder.ToStyle().Render(text)
}

func (w WeeklyStyle) ApplyColorWeekday(text string) string {
	return w.ColorBorder.ToStyle().Render(text)
}

func (w WeeklyStyle) ApplyColorWeekend(text string) string {
	return w.ColorBorder.ToStyle().Render(text)
}

// TODO: NOT IMPLEMENTED
func (w WeeklyStyle) ApplyColorLabel(lID backend.LabelID, we *Weekly) string {
	if lID == "" {
		return ""
	}
	we.Labels.Selected = lID
	return string(we.Labels.GetColor().ToStyle().Render(string(lID)))
}
