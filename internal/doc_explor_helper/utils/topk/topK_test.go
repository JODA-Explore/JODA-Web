package topk

import (
	"math/rand"
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

func TestInsert(t *testing.T) {
	l, rdm := make([]int, 26), make([]int, 26)
	for i := 0; i < 26; i++ {
		rdm[i] = i
		l[i] = i
	}
	rand.Shuffle(len(rdm), func(i, j int) {
		rdm[i], rdm[j] = rdm[j], rdm[i]
	})
	tests := []struct {
		name string
		list []int
		cap  int
		want []int
	}{
		{"5 to 5", l[:5], 5, []int{0, 1, 2, 3, 4}},
		{"26 to 5", l, 5, []int{21, 22, 23, 24, 25}},
		{"rdm", rdm, 5, []int{21, 22, 23, 24, 25}},
		{"4 to 5", l[:4], 5, []int{0, 1, 2, 3}},
		{"1 to 5", []int{1}, 5, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := New[int, int](tt.cap)
			for _, x := range tt.list {
				tp.Insert(x, x)
			}
			l := tp.List()
			got := make([]int, len(l))
			for i, x := range l {
				got[i] = x.Val
			}
			if !slices.OrderlessEqual(got, tt.want) {
				t.Errorf("TopN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	l, rdm := make([]int, 26), make([]int, 26)
	for i := 0; i < 26; i++ {
		rdm[i] = i
		l[i] = i
	}
	rand.Shuffle(len(rdm), func(i, j int) {
		rdm[i], rdm[j] = rdm[j], rdm[i]
	})
	tests := []struct {
		name string
		list []int
		cap  int
		want []int
	}{
		{"5 to 5", l[:5], 5, []int{0, 1, 2, 3, 4}},
		{"26 to 5", l, 5, []int{21, 22, 23, 24, 25}},
		{"rdm", rdm, 5, []int{21, 22, 23, 24, 25}},
		{"4 to 5", l[:4], 5, []int{0, 1, 2, 3}},
		{"1 to 5", []int{1}, 5, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := New[int, int](tt.cap)
			for _, x := range tt.list {
				tp.Insert(x, x)
			}
			tp.Sort()
			l := tp.List()
			got := make([]int, len(l))
			for i, x := range l {
				got[i] = x.Val
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("TopN() = %v, want %v", got, tt.want)
			}
		})
	}
}
