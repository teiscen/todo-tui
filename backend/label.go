package backend

import "slices"

type (
	LabelID   string
	LabelInfo struct {
		Name  string  `yaml:"name"`
		Color HexCode `yaml:"color"`
	}
	Labels struct {
		Selected LabelID               `yaml:"selected"`
		Info     map[LabelID]LabelInfo `yaml:"info"`
		Order    []LabelID             `yaml:"order"`
	}
)

func (l Labels) GetColor() HexCode {
	return l.Info[l.Selected].Color
}

func (l Labels) GetName() string {
	return l.Info[l.Selected].Name
}

func (l *Labels) Step(forward bool) {
	n := len(l.Order)
	if n == 0 {
		return
	}
	i := slices.Index(l.Order, l.Selected)
	if forward {
		i++
	} else {
		i--
	}
	l.Selected = l.Order[(i+n)%n]
}
