package slices

import "testing"

func TestPushFront(t *testing.T) {
	tests := []struct {
		name string
		l    []int
		new  int
		want []int
	}{
		{"1", []int{2, 3}, 1, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PushFront(tt.l, tt.new); !Equal(got, tt.want) {
				t.Errorf("got:%v, want:%v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		l    []int
		keep func(int) bool
		want []int
	}{
		{"1", []int{1, 2, 3, 0, 4}, func(i int) bool { return i > 1 }, []int{2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Filter(&tt.l, tt.keep)
			if !Equal(tt.l, tt.want) {
				t.Errorf("got:%v, want:%v", tt.l, tt.want)
			}
		})
	}
}
