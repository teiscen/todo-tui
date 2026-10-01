package backend

import "time"

type State struct {
	Calendar      Calendar
	SelectedDate  Date
	SelectedLabel LabelID
	Labels        Labels
	Colors        Colors
	NoteFocused   bool
}

func (s State) getLabel() LabelInfo {
	return s.Labels.Info[s.SelectedLabel]
}

func (s State) PrevLabel() LabelID {
	for i, id := range s.Labels.Order {
		if id == s.SelectedLabel {
			return s.Labels.Order[(i-1+len(s.Labels.Order))%len(s.Labels.Order)]
		}
	}

	return s.SelectedLabel
}

func (s State) NextLabel() LabelID {
	for i, id := range s.Labels.Order {
		if id == s.SelectedLabel {
			return s.Labels.Order[(i+1)%len(s.Labels.Order)]
		}
	}

	return s.SelectedLabel
}

func (s State) getSelectedDateColor() HexCode {
	_, ok := s.Calendar.getEntry(s.SelectedDate, s.SelectedLabel)
	color := s.Colors["valid"]
	if ok {
		color = s.Colors[s.getLabel().ColorID]
	}
	return color
}

func (s State) getColor(cID ColorID) (HexCode, bool) {
	c, ok := s.Colors[cID]
	return c, ok
}

func MockState() State {
	labels, colors := MockConfig()

	selectedLabel := labels.Order[0]

	return State{
		Calendar:      MockCalendar(),
		SelectedDate:  TimeToDate(time.Now()),
		SelectedLabel: selectedLabel,
		Labels:        labels,
		Colors:        colors,
		NoteFocused:   false,
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

	firstDateIdx := int(s.SelectedDate.FirstDay().Weekday())
	lastDateIdx := firstDateIdx + s.SelectedDate.NumDays()

	isValid := func(row, col int) bool {
		idx := row*7 + col
		return firstDateIdx <= idx && idx < lastDateIdx
	}

	date := s.SelectedDate.FirstDay().AddDate(0, 0, -firstDateIdx)

	for row := range grid {
		for col := range grid[row] {
			valid := isValid(row, col)

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

			hovered := date == s.SelectedDate

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

type NoteRenderInfo struct {
	IsFocused bool
	Msg       string
}

func (s State) GetNoteRenderInfo() NoteRenderInfo {
	e, ok := s.Calendar.getEntry(s.SelectedDate, s.SelectedLabel)
	var str string
	if !ok {
		str = "No Entry Found"
	} else {
		str = e.Msg
	}
	return NoteRenderInfo{
		s.NoteFocused,
		str,
	}
}

type PlannerRenderInfo struct {
	Header      ColoredString
	Footer      ColoredString
	BorderColor HexCode
}

func (s State) GetPlannerRenderInfo() PlannerRenderInfo {
	return PlannerRenderInfo{
		ColoredString{
			s.SelectedDate.Format(),
			s.getSelectedDateColor(),
		},
		ColoredString{
			s.getLabel().Name,
			s.Colors[s.getLabel().ColorID],
		},
		s.Colors["valid"],
	}
}

func WriteState(filePath string) error {
	return nil
}

func ReadState(filePath string) error {
	return nil
}
