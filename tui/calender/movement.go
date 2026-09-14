package calender

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func move(currentIndex int, dir Direction, max int) int {
	switch dir {
	case Up:
		if currentIndex > 7 {
			return currentIndex - 7
		}
	case Left:
		if currentIndex%7 != 0 {
			return currentIndex - 1
		}
	case Right:
		if currentIndex%7 != 6 {
			return currentIndex + 1
		}
	case Down:
		if currentIndex < (max - 6) {
			return currentIndex + 7
		}
	}
	return currentIndex
}
