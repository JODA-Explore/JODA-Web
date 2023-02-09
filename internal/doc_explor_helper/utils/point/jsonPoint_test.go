package point

import (
	"testing"
)

func TestDepth(t *testing.T) {
	tests := []struct {
		name string
		p    JsonPoint
		want int
	}{
		{"1", "/entities/hashtags/14/text", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Depth(tt.p); got != tt.want {
				t.Errorf("ContainsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestRoot(t *testing.T) {
	tests := []struct {
		name string
		p    JsonPoint
		want string
	}{
		{"1", "/entities/hashtags/14/text", ""},
		{"1", "entities/hashtags/14/text", "entities"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(tt.p); got != tt.want {
				t.Errorf("ContainsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsArray(t *testing.T) {
	type args struct {
		x JsonPoint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{x: "/entities/hashtags/14/text"}, true},
		{"2", args{x: "/entities/hashtags/14"}, true},
		{"3", args{x: "/14"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsArray(tt.args.x); got != tt.want {
				t.Errorf("ContainsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnnestedArray(t *testing.T) {
	type args struct {
		x JsonPoint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{x: "/entities/hashtags/14/text"}, false},
		{"2", args{x: "/entities/hashtags/14"}, true},
		{"2", args{x: "/entities/1/hashtags/14"}, false},
		{"3", args{x: "/14/1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnnestedArray(tt.args.x); got != tt.want {
				t.Errorf("IsUnnestedArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
