package structure

import (
	"fmt"
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

func Test_disableChildren(t *testing.T) {
	ins := query.NewIns(joda.New("http://localhost:5632"))
	dataset := ins.NewDataset("twitter", nil)
	b := Backend{&dataset, 0.1, feature.NewLimit(5, 400)}
	a, err := ins.GetAttribute("twitter")
	if err != nil {
		panic(err)
	}
	pts := points.New(a)
	tests := []struct {
		name     string
		pt       string
		maxCount int
		check    bool
	}{
		{"/quoted_status", "/quoted_status", 20, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			child, ok := pts.Find(tt.pt)
			if !ok {
				t.Error("can not find " + tt.pt)
			}
			tr := trie.New[*trieV](0)
			child.Walk(func(path string, v *points.Value, depth int) (curCtl trie.Control) {
				tt.maxCount--
				if tt.maxCount < 0 {
					return trie.Break
				}
				if depth > 2 {
					return trie.Continue
				}
				tr.Insert(path, &trieV{Info: b.newInfo(path, v.Info, b.TotalNumber())})
				return trie.Next
			})
			fmt.Println("before:")
			fmt.Printf("%v\n", tr.StringWithValue(func(v *trieV) interface{} { return v.percent.Count() }))
			fmt.Println("-----------------------")
			disableChildren(&tr, 0.1)
			fmt.Println("now:")
			fmt.Printf("%v\n", tr.StringWithValue(func(v *trieV) interface{} { return fmt.Sprintf("%v,%v", v.percent.Count(), v.ctl) }))
			t.Error("check")
		})
	}
}

func Test_groupCounts(t *testing.T) {
	ins := query.NewIns(joda.New("http://localhost:5632"))
	dataset := ins.NewDataset("twitter", nil)
	b := Backend{&dataset, 0.1, feature.NewLimit(5, 400)}
	a, err := ins.GetAttribute("twitter")
	if err != nil {
		panic(err)
	}
	pts := points.New(a)
	tests := []struct {
		name     string
		pt       string
		maxCount int
		check    bool
	}{
		{"/quoted_status", "/quoted_status", 20, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			child, ok := pts.Find(tt.pt)
			if !ok {
				t.Error("can not find " + tt.pt)
			}
			tr := trie.New[*trieV](0)
			child.Walk(func(path string, v *points.Value, depth int) (curCtl trie.Control) {
				if path == "" {
					return trie.Next
				}
				tt.maxCount--
				if tt.maxCount < 0 {
					return trie.Break
				}
				if depth > 2 {
					return trie.Continue
				}
				tr.Insert(path, &trieV{Info: b.newInfo(path, v.Info, b.TotalNumber())})
				return trie.Next
			})
			disableChildren(&tr, 0.1)
			fmt.Printf("%v\n", tr.StringWithValue(func(v *trieV) interface{} { return v.percent.Count() }))
			fmt.Println("-----------------------")
			got := groupCounts(&tr, 0.1)
			for i, x := range got {
				fmt.Printf("%v:\n", i)
				for _, y := range x {
					fmt.Printf("%v(%v)\n", y.JsonPoint(), y.percent.Count())
				}
				fmt.Println("")
			}
			t.Error("check")
		})
	}
}
