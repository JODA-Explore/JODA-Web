package slices

import "github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/constraints"

func Equal[T comparable](a, b []T) bool {
	if a == nil {
		return b == nil
	}
	if len(a) != len(b) {
		return false
	}
	for i, ax := range a {
		if ax != b[i] {
			return false
		}
	}
	return true
}

func Contains[T comparable](l []T, e T) bool {
	for _, x := range l {
		if x == e {
			return true
		}
	}
	return false
}

func ContainsFn[T any](l []T, e T, eq func(x, y T) bool) bool {
	for _, x := range l {
		if eq(x, e) {
			return true
		}
	}
	return false
}

// see https://github.com/golang/go/wiki/SliceTricks#Reversing
func Reverse[T any](a []T) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func NewReversed[T any](a []T) []T {
	res := make([]T, len(a))
	copy(res, a)
	Reverse(res)
	return res
}

func Map[T, R any](l []T, fn func(T) R) []R {
	res := make([]R, len(l))
	for i, x := range l {
		res[i] = fn(x)
	}
	return res
}

func Max[Elem any, N constraints.Number](l []Elem, val func(Elem) N) (e Elem, v N) {
	for i, new := range l {
		newV := val(new)
		if i == 0 || newV > v {
			e = new
			v = newV
		}
	}
	return
}

func PushFront[T any](l []T, new T) []T {
	res := make([]T, len(l)+1)
	copy(res[1:], l)
	res[0] = new
	return res
}

// see https://github.com/golang/go/wiki/SliceTricks#filter-in-place
func Filter[T any](l *[]T, keep func(T) bool) {
	n := 0
	for _, x := range *l {
		if keep(x) {
			(*l)[n] = x
			n++
		}
	}
	*l = (*l)[:n]
}

func OrderlessEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		if !Contains(b, x) {
			return false
		}
	}
	return true
}

func OrderlessEqualFn[T any](a, b []T, eq func(x, y T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		if !ContainsFn(b, x, eq) {
			return false
		}
	}
	return true
}

func OrderlessRemove[T any](s []T, i int) []T {
	l := len(s)
	if i != l-1 {
		s[i] = s[l-1]
	}
	return s[:l-1]
}
