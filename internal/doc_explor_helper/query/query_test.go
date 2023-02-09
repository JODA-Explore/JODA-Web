package query

import (
	"fmt"
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

func TestIns_GetAllDatasets(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	type fields struct {
		Joda *joda.Joda
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"1", fields{&ji}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{
				Joda: tt.fields.Joda,
			}
			got, err := ins.AllDatasets()
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got)
		})
	}
}

func TestIns_Count(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name  string
		query string
		check bool
	}{
		{"1", "LOAD twitter", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.Count(tt.query)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got)
		})
	}
}

func TestIns_GetAttribute(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name    string
		dataset string
		check   bool
	}{
		{"1", "twitter", false},
	}
	ins := Ins{Joda: &ji}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := ins.GetAttribute(tt.dataset)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got)
		})
	}
}

func TestIns_DistinctValuesFreq(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name      string
		dataset   string
		jsonPoint string
		topN      int
		check     bool
	}{
		{"2", "twitter", "/lang", 10, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.DistinctValuesFreq(tt.dataset, tt.jsonPoint, tt.topN)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got.List())
		})
	}
}

func TestIns_AtomicMemberFreq(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name      string
		dataset   string
		jsonPoint string
		topN      int
		check     bool
	}{
		{"1", "movies", "/genres", 4, false},
		{"2", "twitter", "/retweeted_status/geo/coordinates", 10, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.AtomicMemberFreq(tt.dataset, tt.jsonPoint, tt.topN)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got.List())
		})
	}
}

func TestIns_DistinctObjCounts(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name      string
		dataset   string
		jsonPoint string
		check     bool
	}{
		{"1", "twitter", "/retweeted_status/entities/hashtags", false},
		{"2", "twitter", "/quoted_status/place/bounding_box/coordinates", false},
		{"3", "__AUTO_GENERATED__twitter__sample__", "/quoted_status/place/bounding_box/coordinates", false},
		{"4", "__AUTO_GENERATED__twitter__sample__", "/retweeted_status/quoted_status/place/bounding_box/coordinates", false},
		{"4", "twitter", "/retweeted_status/extended_tweet/entities/urls", true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			count, total, err := ins.DistinctObjCounts(tt.dataset, tt.jsonPoint)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v,%v", count, total)
		})
	}
}

func TestIns(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	t.Run("test", func(t *testing.T) {
		ins := Ins{Joda: &ji}
		ins.memberFreq(cmd.Load("hashed"), 1)
		// v, err := ins.getResult(cmd.Load("hashed"))
		// if err != nil {
		// 	panic(err)
		// }
		// va, err := v.Array()
		// if err != nil {
		// 	panic(err)
		// }
		// buf := va[0].Get("group").MarshalTo(nil)
		t.Errorf("%v", "string(buf)")
	})
}

func TestIns_CountOfDistinctValues(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name    string
		dataset string
		jps     []string
		check   bool
	}{
		{"1", "twitter", []string{"/retweet_count", "/filter_level"}, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.CountOfDistinctValues(tt.dataset, tt.jps)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", got)
		})
	}
}

func TestIns_TopDistinctValues(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name    string
		dataset string
		jp      string
		n       int
		check   bool
	}{
		{"1", "twitter", "/lang", 10, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			members, counts, err := ins.TopDistinctValues(tt.dataset, tt.jp, tt.n)
			if err != nil {
				panic(err)
			}
			for i, x := range members {
				fmt.Printf("member:%v count:%v\n", x, counts[i])
			}
			t.Error("check")
		})
	}
}

func TestIns_AllDistinctValues(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name    string
		dataset string
		path    string
		check   bool
	}{
		{"1", "twitter", "/lang", true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.AllDistinctValues(tt.dataset, tt.path)
			if err != nil {
				panic(err)
			}
			size := num.Min(10, len(got))
			t.Errorf("%v", got[:size])
		})
	}
}

func TestIns_AllDistinctMembers(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name    string
		dataset string
		path    string
		check   bool
	}{
		{"1", "movies", "/genres", true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			got, err := ins.AllDistinctMembers(tt.dataset, tt.path)
			if err != nil {
				panic(err)
			}
			size := num.Min(10, len(got))
			t.Errorf("%v", got[:size])
		})
	}
}

func TestIns_ValueCount(t *testing.T) {
	ji := joda.New("http://localhost:5632")
	tests := []struct {
		name      string
		dataset   string
		jsonPoint string
		value     string
		jsontype  jsontype.Type
		check     bool
	}{
		{"1", "twitter", "/favorite_count", "0", jsontype.Int, false},
		{"2", "twitter", "/lang", "en", jsontype.String, false},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			ins := Ins{Joda: &ji}
			count, err := ins.ValueCount(tt.dataset, tt.jsonPoint, tt.value, tt.jsontype)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", count)
		})
	}
}
