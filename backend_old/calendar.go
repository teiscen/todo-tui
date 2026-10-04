package backend

type Calendar struct {
	Months map[Date]map[LabelID]Entry
	Labels Labels
}

func (c Calendar) GetEntry(d Date) (val Entry, ok bool) {
	val, ok = c.Months[d][c.Labels.Selected]
	return
}

func (c Calendar) GetCurrentLabel() LabelInfo {
	return c.Labels.Info[c.Labels.Selected]
}
