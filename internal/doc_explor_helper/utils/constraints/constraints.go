package constraints

type Number interface {
	~int | ~uint64 | ~float64
}
