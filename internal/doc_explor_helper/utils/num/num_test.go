package num

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	floatEqual := func(x, y float64) bool {
		return fmt.Sprintf("%.3f", x) == fmt.Sprintf("%.3f", y)
	}
	type args struct {
		x uint64
		y uint64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"1", args{x: 5, y: 4}, 0.2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.x, tt.args.y); !floatEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		z uint64
		x int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"1", args{1000, 3}, 333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Div(tt.args.z, tt.args.x); got != tt.want {
				t.Errorf("Quo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name string
		x, y uint64
		want float64
	}{
		{"1", 10, 10, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := percent(tt.x, tt.y); got != tt.want {
				t.Errorf("Quo() = %v, want %v", got, tt.want)
			}
		})
	}
}
