package backend

type Calendar struct {
	Months map[Date]map[LabelID]Entry
}

func (c Calendar) getEntry(d Date, l LabelID) (val Entry, ok bool) {
	val, ok = c.Months[d][l]
	return
}
