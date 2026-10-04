package backend

import "slices"

type (
	LabelID   string
	LabelInfo struct {
		Name  string
		Color HexCode
	}
	Labels struct {
		Selected LabelID
		Info     map[LabelID]LabelInfo
		Order    []LabelID
	}
)

func (l *Labels) Step(forward bool) {
	n := len(l.Order)
	i := slices.Index(l.Order, l.Selected)
	if forward {
		i++
	} else {
		i--
	}
	l.Selected = l.Order[(i+n)%n]
}

// func (m *Model) addLabel(label Label) {
// 	m.labels = append(m.labels, label)
// }
//
// func (m *Model) removeLabel(i int) {
// 	m.labels = append(m.labels[:i], m.labels[i+1:]...)
//
// 	if len(m.labels) == 0 {
// 		m.selLabel = 0
// 	} else if m.selLabel >= len(m.labels) {
// 		m.selLabel = len(m.labels) - 1
// 	}
// }
