package topk

import (
	"sort"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/constraints"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/heap"
)

type TopK[V any, N constraints.Number] struct {
	h        heap.Heap[V, N]
	len, cap int
}

func New[V any, N constraints.Number](cap int) (res TopK[V, N]) {
	res.cap = cap
	res.h = make(heap.Heap[V, N], cap, cap)
	return
}

func (tp *TopK[V, N]) Insert(v V, n N) {
	if tp.len < tp.cap {
		tp.h[tp.len] = heap.NewElem(v, n)
		tp.len++
		return
	}
	if tp.h[0].N < n {
		tp.h[0] = heap.NewElem(v, n)
		tp.h.Fix(0)
	}
}

func (tp *TopK[V, N]) List() []heap.Elem[V, N] {
	if tp.len < tp.cap {
		return tp.h[:tp.len]
	}
	return tp.h
}

func (tp *TopK[V, N]) Vals() []V {
	list := tp.List()
	vals := make([]V, len(list))
	for i, x := range list {
		vals[i] = x.Val
	}
	return vals
}

func (tp *TopK[V, N]) Len() int {
	return tp.len
}

func (tp *TopK[V, N]) clean() {
	if tp.len < tp.cap {
		tp.cap = tp.len
		tp.h = tp.h[:tp.len]
	}
}

func (tp *TopK[V, N]) Sort() {
	tp.clean()
	sort.Slice(tp.h, func(i, j int) bool { return tp.h[i].N < tp.h[j].N })
}

func (tp *TopK[V, N]) SortReverse() {
	tp.clean()
	sort.Slice(tp.h, func(i, j int) bool { return tp.h[i].N > tp.h[j].N })
}
