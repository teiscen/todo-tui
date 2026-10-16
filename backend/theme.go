package backend

// Tone classifies how urgent a task's due date is, so the UI can color it.
type Tone int

const (
	ToneNormal Tone = iota
	ToneSoon
	ToneLate
)

// Theme is the single set of colors every pane draws with. Panes hold a Theme
// by value; the root constructs one from config (overlaying DefaultTheme) and
// passes it down so colors are never duplicated across packages.
type Theme struct {
	Text    HexCode
	Border  HexCode
	Weekday HexCode
	Weekend HexCode
	Dim     HexCode
	Accent  HexCode
	Late    HexCode
	Soon    HexCode
}

// DefaultTheme returns the built-in color scheme. Config can override the
// Text/Border/Weekday/Weekend entries; the rest keep these defaults.
func DefaultTheme() Theme {
	return Theme{
		Text:    HexCode("#C6D0F5"),
		Border:  HexCode("#C6D0F5"),
		Weekday: HexCode("#C6D0F5"),
		Weekend: HexCode("#939AB7"),
		Dim:     HexCode("#939AB7"),
		Accent:  HexCode("#CA9EE6"),
		Late:    HexCode("#E78284"),
		Soon:    HexCode("#E5C890"),
	}
}

// ToneColor maps a due-date urgency to its display color.
func (t Theme) ToneColor(tone Tone) HexCode {
	switch tone {
	case ToneLate:
		return t.Late
	case ToneSoon:
		return t.Soon
	}
	return t.Dim
}
