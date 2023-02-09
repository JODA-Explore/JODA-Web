package status

type Status []bool

func New(n int) Status {
	return make(Status, n)
}

func (s Status) Disable(idx int) {
	s[idx] = true
}

func (s Status) IsDisabled(idx int) bool {
	return s[idx]
}
