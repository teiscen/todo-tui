package backend

type (
	LabelID   string
	LabelInfo struct {
		Name    string
		ColorID ColorID
	}
	Labels struct {
		Info  map[LabelID]LabelInfo
		Order []LabelID
	}
)

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
