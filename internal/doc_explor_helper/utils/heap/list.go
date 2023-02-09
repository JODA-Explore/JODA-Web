package heap

// import "github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/constraints"

// type List[T constraints.Number] []T

// func push[T constraints.Number](l *List[T], x T) {
// 	*l = append(*l, x)
// }

// func less[T constraints.Number](l *List[T], i, j int) bool {
// 	return (*l)[i] < (*l)[j]
// }

// func swap[T constraints.Number](l *List[T], i, j int) {
// 	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
// }

// func pop[T constraints.Number](l *List[T]) (v T) {
// 	*l, v = (*l)[:len(*l)-1], (*l)[len(*l)-1]
// 	return
// }

// // Init establishes the heap invariants required by the other routines in this package.
// // Init is idempotent with respect to the heap invariants
// // and may be called whenever the heap invariants may have been invalidated.
// // The complexity is O(n) where n = h.Len().
// func Init[V constraints.Number](h *List[V]) {
// 	// heapify
// 	n := len(*h)
// 	for i := n/2 - 1; i >= 0; i-- {
// 		down(h, i, n)
// 	}
// }

// // Push pushes the element x onto the heap.
// // The complexity is O(log n) where n = h.Len().
// func Push[V constraints.Number](h *List[V], x V) {
// 	push(h, x)
// 	up(h, len(*h)-1)
// }

// // Pop removes and returns the minimum element (according to Less) from the heap.
// // The complexity is O(log n) where n = len(*h).
// // Pop is equivalent to Remove(h, 0).
// func Pop[V constraints.Number](h *List[V]) V {
// 	n := len(*h) - 1
// 	swap(h, 0, n)
// 	down(h, 0, n)
// 	return pop(h)
// }

// // Remove removes and returns the element at index i from the heap.
// // The complexity is O(log n) where n = len(*h).
// func Remove[V constraints.Number](h *List[V], i int) V {
// 	n := len(*h) - 1
// 	if n != i {
// 		swap(h, i, n)
// 		if !down(h, i, n) {
// 			up(h, i)
// 		}
// 	}
// 	return pop(h)
// }

// // Fix re-establishes the heap ordering after the element at index i has changed its value.
// // Changing the value of the element at index i and then calling Fix is equivalent to,
// // but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// // The complexity is O(log n) where n = len(*h).
// func Fix[V constraints.Number](h *List[V], i int) {
// 	if !down(h, i, len(*h)) {
// 		up(h, i)
// 	}
// }

// func up[V constraints.Number](h *List[V], j int) {
// 	for {
// 		i := (j - 1) / 2 // parent
// 		if i == j || !less(h, j, i) {
// 			break
// 		}
// 		swap(h, i, j)
// 		j = i
// 	}
// }

// func down[V constraints.Number](h *List[V], i0, n int) bool {
// 	i := i0
// 	for {
// 		j1 := 2*i + 1
// 		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
// 			break
// 		}
// 		j := j1 // left child
// 		if j2 := j1 + 1; j2 < n && less(h, j2, j1) {
// 			j = j2 // = 2*i + 2  // right child
// 		}
// 		if !less(h, j, i) {
// 			break
// 		}
// 		swap(h, i, j)
// 		i = j
// 	}
// 	return i > i0
// }
