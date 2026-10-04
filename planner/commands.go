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
	p.NotesModel.ToggleFocus()
}

func (p *Planner) UpdateMonth(next bool) {
	curr := p.CalendarModel.SelectedDate
	if next {
		p.CalendarModel.SelectedDate = curr.AddDate(0, 1, 0)
	} else {
		p.CalendarModel.SelectedDate = curr.AddDate(0, -1, 0)
	}
	p.CalendarModel.UpdateGrid(p.Calendar)
}

func (p *Planner) UpdateLabel(next bool) {
	if next {
		p.Calendar.Labels.Step(true)
	} else {
		p.Calendar.Labels.Step(false)
	}
	p.CalendarModel.Style.ColorAccent = p.Calendar.GetCurrentLabel().Color
	p.CalendarModel.UpdateGrid(p.Calendar)
}
