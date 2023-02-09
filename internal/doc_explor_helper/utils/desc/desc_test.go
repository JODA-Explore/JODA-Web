package desc

import (
	"html/template"
	"testing"
)

func TestButtonBarMaker_Make(t *testing.T) {
	bm1 := ButtonBarMaker{
		1,
		"id",
		"twitter",
		[]struct {
			ButtonName, DocId string
			DivContent        template.HTML
		}{
			{"buttonA", "idA", "<b>testA<b>"},
			{"buttonB", "idB", "<b>testB<b>"},
		},
	}
	tests := []struct {
		name  string
		bm    ButtonBarMaker
		check bool
	}{
		{"1", bm1, false},
	}
	for _, tt := range tests {
		if tt.check {
			t.Run(tt.name, func(t *testing.T) {
				t.Errorf("%v", tt.bm.Make())
			})
		}
	}
}
