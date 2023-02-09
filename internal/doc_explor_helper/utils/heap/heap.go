// a generic version based on the go official heap implementation

package heap

import "github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/constraints"

type Elem[V any, N constraints.Number] struct {
	Val V
	N   N
}

type Heap[V any, N constraints.Number] []Elem[V, N]

func push[V any, N constraints.Number](l *Heap[V, N], x Elem[V, N]) {
	*l = append(*l, x)
}

func less[V any, N constraints.Number](l *Heap[V, N], i, j int) bool {
	return (*l)[i].N < (*l)[j].N
}

func swap[V any, N constraints.Number](l *Heap[V, N], i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

func pop[V any, N constraints.Number](l *Heap[V, N]) (v Elem[V, N]) {
	*l, v = (*l)[:len(*l)-1], (*l)[len(*l)-1]
	return
}

func NewElem[V any, N constraints.Number](v V, n N) Elem[V, N] {
	return Elem[V, N]{v, n}
}

func New[V any, N constraints.Number]() *Heap[V, N] {
	return &Heap[V, N]{}
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func (h *Heap[V, N]) Init() {
	// heapify
	n := len(*h)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[V, N]) Push(x Elem[V, N]) {
	push(h, x)
	h.up(len(*h) - 1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = len(*h).
// Pop is equivalent to Remove(h, 0).
func (h *Heap[V, N]) Pop() Elem[V, N] {
	n := len(*h) - 1
	swap(h, 0, n)
	h.down(0, n)
	return pop(h)
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = len(*h).
func (h *Heap[V, N]) Remove(i int) Elem[V, N] {
	n := len(*h) - 1
	if n != i {
		swap(h, i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return pop(h)
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = len(*h).
func (h *Heap[V, N]) Fix(i int) {
	if !h.down(i, len(*h)) {
		h.up(i)
	}
}

func (h *Heap[V, N]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(h, j, i) {
			break
		}
		swap(h, i, j)
		j = i
	}
}

func (h *Heap[V, N]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(h, j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(h, j, i) {
			break
		}
		swap(h, i, j)
		i = j
	}
	return i > i0
}
