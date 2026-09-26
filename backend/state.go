package backend

import "time"

type State struct {
	Calendar      Calendar
	SelectedDate  Date
	SelectedLabel LabelID
	Labels        Labels
	Colors        Colors
}

func (s State) getLabel() LabelInfo {
	return s.Labels[s.SelectedLabel]
}

func (s State) getColor(cID ColorID) (HexCode, bool) {
	c, ok := s.Colors[cID]
	return c, ok
}

func MockState() State {
	labels, colors := MockConfig()
	// get an arbitrary LabelID to start
	var selectedLabel LabelID
	for k := range labels {
		selectedLabel = k
		break
	}

	return State{
		Calendar:      MockCalendar(),
		SelectedDate:  TimeToDate(time.Now()),
		SelectedLabel: selectedLabel,
		Labels:        labels,
		Colors:        colors,
	}
}

// const (
// 	NoEntry Statys = iota
// 	Partial
// 	Full
// )

type GridCell struct {
	Color    HexCode
	Status   Status
	Selected bool
	IsValid  bool
}
type Grid [6][7]GridCell

func (s State) getGrid() (grid Grid) {
	cValid, _ := s.getColor("valid")
	cInvalid, _ := s.getColor("invalid")
	cLabel, _ := s.getColor(s.getLabel().ColorID)

	isValid := func(row, col int) bool {
		idx := row*7 + col
		firstDateIdx := int(s.SelectedDate.FirstDay().Weekday())
		lastDateIdx := firstDateIdx + s.SelectedDate.NumDays()

		if firstDateIdx < idx && idx < lastDateIdx {
			return true
		} else {
			return false
		}
	}

	date := s.SelectedDate.FirstDay()
	for row := range grid {
		for col := range grid[row] {
			valid := isValid(row, col)

			// Determine Color and Status
			var color HexCode
			var status Status
			val, ok := s.Calendar.getEntry(date, s.SelectedLabel)
			if ok {
				color = cLabel
				status = val.Status
			} else {
				if valid {
					color = cValid
				} else {
					color = cInvalid
				}
				status = NoEntry
			}

			hovered := false
			if date == s.SelectedDate {
				hovered = true
			}

			grid[row][col] = GridCell{
				color,
				status,
				hovered,
				valid,
			}

			date = date.AddDate(0, 0, 1)
		}
	}

	return
}

type CalendarRenderInfo struct {
	Grid Grid
}

// Need to handle the Color not existing
func (s State) GetCalendarRenderInfo() CalendarRenderInfo {
	g := s.getGrid()

	return CalendarRenderInfo{
		// labelStr,
		// dateString,
		g,
	}
}

type NoteRenderInfo struct{}

func (s State) GetNoteRenderInfo() NoteRenderInfo {
	return NoteRenderInfo{}
}

func WriteState(filePath string) error {
	return nil
}

func ReadState(filePath string) error {
	return nil
}
