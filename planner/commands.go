package planner

// Command:
// Switch between the labels (updates calendar)
// Swtich between the months (updates calendar)

func (p *Planner) ToggleFocus() {
	if p.Focus == FocusCalendar {
		p.Focus = FocusNotes
	} else {
		p.Focus = FocusCalendar
	}
}

func (p *Planner) UpdateMonth(next bool) {
	if next {
		p.CalendarModel.SelectedDate.AddDate(0, 1, 0)
	} else {
		p.CalendarModel.SelectedDate.AddDate(0, -1, 0)
	}
	p.CalendarModel.UpdateGrid(p.Calendar)
}

func (p *Planner) UpdateLabel(next bool) {
	if next {
		p.Calendar.Labels.Step(true)
	} else {
		p.Calendar.Labels.Step(false)
	}
	p.CalendarModel.UpdateGrid(p.Calendar)
}
