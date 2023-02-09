package rating

import (
	"reflect"
	"testing"
)

func Test_applyScale(t *testing.T) {
	type args struct {
		x     float64
		scale float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyScale(tt.args.x, tt.args.scale); got != tt.want {
				t.Errorf("applyScale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpline_applyScale(t *testing.T) {
	tests := []struct {
		name   string
		spline Spline
		scale  float64
		want   Spline
	}{
		{"coverage", Coverage, 45, Coverage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spline.ApplyScale(tt.scale); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Spline.applyScale() = %v, want %v", got, tt.want)
			}
		})
	}
}
