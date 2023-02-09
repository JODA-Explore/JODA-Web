package feature

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

func Test_groupSingleInfo(t *testing.T) {
	genPred := func(invalids []int) func(x, y int) bool {
		return func(x, y int) bool {
			return !slices.Contains(invalids, x) &&
				!slices.Contains(invalids, y)
		}
	}
	tests := []struct {
		name        string
		counts      []int
		wantRest    []int
		wantGrouped []int
	}{
		{"a", []int{1, 2, 3, 4, 5}, []int{1}, []int{2, 3, 4, 5}},
		{"b", []int{1, 2, 3, 4, 5}, []int{1, 3}, []int{2, 4, 5}},
		{"c", []int{1, 2, 3, 4, 5}, []int{4, 3}, []int{1, 2, 5}},
		{"c", []int{1, 2, 3, 4, 5}, []int{2, 3, 4}, []int{1, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, single, rest := groupSingleInfo(tt.counts, genPred(tt.wantRest))
			invalid := append(single, rest...)
			if !slices.OrderlessEqual(valid, tt.wantGrouped) || !slices.OrderlessEqual(invalid, tt.wantRest) {
				t.Errorf("got valid:%v,want valid:%v\ngot invalid:%v,want invalid:%v", valid, tt.wantGrouped, invalid, tt.wantRest)
			}
		})
	}
}

func Test_GroupInfo(t *testing.T) {
	orderless2DSlicesEq := func(a, b [][]int) bool {
		return slices.OrderlessEqualFn(a, b, func(x, y []int) bool {
			return reflect.DeepEqual(x, y)
		})
	}
	genPred := func(groups [][]int) func(x, y int) bool {
		return func(x, y int) bool {
			for _, g := range groups {
				if slices.Contains(g, x) && slices.Contains(g, y) {
					return true
				}
			}
			return false
		}
	}
	tests := []struct {
		name        string
		counts      []int
		wantRest    []int
		wantGrouped [][]int
	}{
		{"a", []int{1, 2, 3, 4, 5}, []int{1}, [][]int{{2, 3, 4, 5}}},
		{"b", []int{1, 2, 3, 4, 5, 6}, []int{1, 3}, [][]int{{2, 4}, {5, 6}}},
		{"c", []int{1, 2, 3, 4, 5}, []int{4}, [][]int{{1, 2}, {3, 5}}},
		{"d", []int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4}, [][]int{{5, 6}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupInfo(tt.counts, genPred(tt.wantGrouped))
			want := tt.wantGrouped
			for _, x := range tt.wantRest {
				want = append(want, []int{x})
			}
			if !orderless2DSlicesEq(got, want) {
				t.Errorf("got:%v,want:%v", got, want)
			}
		})
	}
	for i := 0; i < 50; i++ {
		t.Run("random", func(t *testing.T) {
			l := 10
			ori := make([]int, l)
			for i := range ori {
				ori[i] = i
			}
			var grouped [][]int
			rdm := rand.Perm(l)
			j := 0
			for {
				newJ := j + 2 + rand.Intn(l-j-3)/2
				elem := make([]int, newJ-j)
				copy(elem, rdm[j:newJ])
				sort.Ints(elem)
				grouped = append(grouped, elem)
				j = newJ
				if j >= l-3 {
					break
				}
			}
			sort.Slice(grouped, func(i, j int) bool { return grouped[i][0] < grouped[j][0] })
			// fmt.Printf("grouped:%v\n", grouped)
			rest := rdm[j:]
			got := GroupInfo(ori, genPred(grouped))
			want := grouped
			sort.Ints(rest)
			for _, x := range rest {
				want = append(want, []int{x})
			}
			if !orderless2DSlicesEq(got, want) {
				t.Errorf("got:%v,want:%v", got, want)
			}
		})
	}
}

func BenchmarkMerger(b *testing.B) {
	avg := 0
	genPred := func(groups [][]int) func(x, y int) bool {
		return func(x, y int) bool {
			avg++
			time.Sleep(50 * time.Millisecond)
			for _, g := range groups {
				if slices.Contains(g, x) && slices.Contains(g, y) {
					return true
				}
			}
			return false
		}
	}
	length := 20
	ori := make([]int, length)
	for i := range ori {
		ori[i] = i
	}
	for i := 0; i < b.N; i++ {
		rdm := rand.Perm(length)
		j := 0
		var grouped [][]int
		for {
			newJ := j + 2 + rand.Intn(length-j-3)/2
			elem := make([]int, newJ-j)
			copy(elem, rdm[j:newJ])
			grouped = append(grouped, elem)
			j = newJ
			if j >= length-3 {
				break
			}
		}
		GroupInfo(ori, genPred(grouped))
	}
	// fmt.Printf("avg:%v\n", float64(avg)/float64(b.N))
}

func BenchmarkMergerBF(b *testing.B) {
	avg := 0
	genPred := func(groups [][]int) func(x, y int) bool {
		return func(x, y int) bool {
			time.Sleep(50 * time.Millisecond)
			avg++
			for _, g := range groups {
				if slices.Contains(g, x) && slices.Contains(g, y) {
					return true
				}
			}
			return false
		}
	}
	length := 20
	ori := make([]int, length)
	for i := range ori {
		ori[i] = i
	}
	for i := 0; i < b.N; i++ {
		rdm := rand.Perm(length)
		j := 0
		var grouped [][]int
		for {
			newJ := j + 2 + rand.Intn(length-j-3)/2
			elem := make([]int, newJ-j)
			copy(elem, rdm[j:newJ])
			grouped = append(grouped, elem)
			j = newJ
			if j >= length-3 {
				break
			}
		}
		var res [][]int
		for i, x := range ori {
			g := []int{x}
			for j := i + 1; j < len(ori); j++ {
				y := ori[j]
				if genPred(grouped)(x, y) {
					g = append(g, y)
				}
			}
			res = append(res, g)
		}
	}
	// fmt.Printf("avg:%v\n", float64(avg)/float64(b.N))
}
