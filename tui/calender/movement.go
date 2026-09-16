package calender

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func move(c CalenderDate, dir Direction) CalenderDate {
	switch dir {
	case Up:
		return c.addDate(0, 0, -7)
	case Left:
		return c.addDate(0, 0, -1)
	case Right:
		return c.addDate(0, 0, 1)
	case Down:
		return c.addDate(0, 0, 7)
	default:
		return c
	}
}
