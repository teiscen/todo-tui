package backend

type Calendar struct {
	Months map[Date]map[LabelID]Entry
	Labels Labels
}

func (c Calendar) getEntry(d Date, l LabelID) (val Entry, ok bool) {
	val, ok = c.Months[d][l]
	return
}

func (c Calendar) GetCurrentLabel() LabelInfo {
	return c.Labels.Info[c.Labels.Selected]
}
