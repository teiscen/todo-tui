package calendar

type (
	Label     string
	LabelInfo struct {
		Name      string
		ColorName ColorName
	}
)

func (m *Model) addLabel(label Label) {
	m.labels = append(m.labels, label)
}

func (m *Model) removeLabel(i int) {
	m.labels = append(m.labels[:i], m.labels[i+1:]...)

	if len(m.labels) == 0 {
		m.selLabel = 0
	} else if m.selLabel >= len(m.labels) {
		m.selLabel = len(m.labels) - 1
	}
}
