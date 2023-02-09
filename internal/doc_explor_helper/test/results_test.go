package test

import (
	"strings"
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/explore"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

var ins query.Ins

func TestValueType(t *testing.T) {
	ins = query.NewIns(joda.New("http://localhost:5632"))
	searchOpt := explore.SearchOpt{Dataset: "twitter", Filter: filter.Filter{MaxNumber: 20, MaxDiff: 0.1}}
	searchOpt.EnabledBackends[backends.ValueType] = true
	tests := []struct {
		name      string
		opt       explore.SearchOpt
		resultIdx int
		loopTime  int
		check     bool
	}{
		{"ori", searchOpt, 0, 0, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		for i := 0; i < tt.loopTime; i++ {
			explore.HandleNewDataset(ins, "exists", "", "", tt.opt)
			resultsInfo, err := explore.Results(ins, tt.opt)
			if err != nil {
				panic(err)
			}
			if !strings.Contains(string(resultsInfo.Results[0].Desc), "/geo") {
				t.Errorf("%v", resultsInfo.Results[tt.resultIdx])
			}
		}
	}
}

func TestStructure(t *testing.T) {
	ins = query.NewIns(joda.New("http://localhost:5632"))
	searchOpt := explore.SearchOpt{Dataset: "twitter", Filter: filter.Filter{MaxNumber: 20, MaxDiff: 0.1}}
	searchOpt.EnabledBackends[backends.StructureDiff] = true
	tests := []struct {
		name      string
		opt       explore.SearchOpt
		resultIdx int
		check     bool
	}{
		{"ori", searchOpt, 0, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			explore.HandleNewDataset(ins, "exists", "", "", tt.opt)
			resultsInfo, err := explore.Results(ins, tt.opt)
			if err != nil {
				panic(err)
			}
			t.Errorf("%v", resultsInfo.Results[tt.resultIdx])
		})
	}
}

func TestDistinct(t *testing.T) {
	ins = query.NewIns(joda.New("http://localhost:5632"))
	searchOpt := explore.SearchOpt{Dataset: "twitter", Filter: filter.Filter{MaxNumber: 20, MaxDiff: 0.1}}
	searchOpt.EnabledBackends[backends.Distinct] = true
	tests := []struct {
		name      string
		opt       explore.SearchOpt
		resultIdx int
		loopTime  int
		check     bool
	}{
		{"ori", searchOpt, 0, 1, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		for i := 0; i < tt.loopTime; i++ {
			t.Run(tt.name, func(t *testing.T) {
				explore.HandleNewDataset(ins, "exists", "", "", tt.opt)
				resultsInfo, err := explore.Results(ins, tt.opt)
				if err != nil {
					panic(err)
				}
				t.Errorf("%v", resultsInfo.Results[tt.resultIdx])
			})
		}
	}
}

func TestArray(t *testing.T) {
	ins = query.NewIns(joda.New("http://localhost:5632"))
	searchOpt := explore.SearchOpt{Dataset: "movies", Filter: filter.Filter{MaxNumber: 20, MaxDiff: 0.1}}
	searchOpt.EnabledBackends[backends.Array] = true
	tests := []struct {
		name      string
		opt       explore.SearchOpt
		resultIdx int
		loopTime  int
		check     bool
	}{
		{"ori", searchOpt, 0, 1, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		for i := 0; i < tt.loopTime; i++ {
			t.Run(tt.name, func(t *testing.T) {
				explore.HandleNewDataset(ins, "exists", "", "", tt.opt)
				resultsInfo, err := explore.Results(ins, tt.opt)
				if err != nil {
					panic(err)
				}
				t.Errorf("%v", resultsInfo.Results[tt.resultIdx])
			})
		}
	}
}

func TestObjects(t *testing.T) {
	ins = query.NewIns(joda.New("http://localhost:5632"))
	searchOpt := explore.SearchOpt{Dataset: "twitter", Filter: filter.Filter{MaxNumber: 20, MaxDiff: 0.1}}
	searchOpt.EnabledBackends[backends.Objects] = true
	tests := []struct {
		name      string
		opt       explore.SearchOpt
		resultIdx int
		loopTime  int
		check     bool
	}{
		{"ori", searchOpt, 0, 1, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		for i := 0; i < tt.loopTime; i++ {
			t.Run(tt.name, func(t *testing.T) {
				explore.HandleNewDataset(ins, "exists", "", "", tt.opt)
				resultsInfo, err := explore.Results(ins, tt.opt)
				if err != nil {
					panic(err)
				}
				t.Errorf("%v", resultsInfo.Results[tt.resultIdx])
			})
		}
	}
}
