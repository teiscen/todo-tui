package calendar

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (m *Model) moveCursor(dir Direction) {
	switch dir {
	case Up:
		m.selDate = m.selDate.addDate(0, 0, -7)
	case Down:
		m.selDate = m.selDate.addDate(0, 0, 7)
	case Left:
		m.selDate = m.selDate.addDate(0, 0, -1)
	case Right:
		m.selDate = m.selDate.addDate(0, 0, 1)
	}
}

func (m *Model) movePage(dir Direction) {
	switch dir {
	case Left:
		m.selLabel--
		if m.selLabel < 0 {
			m.selLabel = len(m.labels) - 1
		}

	case Right:
		m.selLabel++
		if m.selLabel >= len(m.labels) {
			m.selLabel = 0
		}
	}
}

// TODO: Shift focus logic to movement
func (m *Model) moveFocus(dir Direction) {
}
