package attstats

import (
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/joda"

	"github.com/valyala/fastjson"
)

func TestNew(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	ins := query.NewIns(ji)
	twitter, err := ins.GetAttribute("twitter")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name  string
		v     *fastjson.Value
		fn    func(Attstats) interface{}
		check bool
	}{
		{
			"CountTotal", twitter,
			func(as Attstats) interface{} {
				return Extract(&as, as.NonArray, func(as AttStat) uint64 {
					return as.CountTotal
				})
			},
			false,
		},
		{
			"user/id", twitter,
			func(as Attstats) interface{} {
				return as.Find("/user/id").MaxInt
			},
			false,
		},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.v)
			t.Errorf("%v", tt.fn(got))
		})
	}
}

func BenchmarkNew(b *testing.B) {
	ji := joda.New("http://localhost:5632")
	ins := query.NewIns(ji)
	twitter, _ := ins.GetAttribute("twitter")
	for i := 0; i < 100; i++ {
		New(twitter)
	}
}
