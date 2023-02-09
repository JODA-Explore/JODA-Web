package cond

import "testing"

func TestConditions_Query(t *testing.T) {
	type args struct {
		dataSet string
	}
	tests := []struct {
		name string
		cs   Conditions
		args args
		want string
	}{
		{"1", Conditions{New("/user/profile_background_tile", false, Equal, "false")}, args{"twitter"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.Query(tt.args.dataSet); got != tt.want {
				t.Errorf("Conditions.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
