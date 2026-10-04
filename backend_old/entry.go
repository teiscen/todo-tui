package backend

type (
	Status int
	Entry  struct {
		Status Status
		Msg    string
	}
)

const (
	NoEntry Status = iota
	Partial
	Full
)
