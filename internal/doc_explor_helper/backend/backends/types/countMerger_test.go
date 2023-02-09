package types

import (
	"fmt"
	"testing"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

func TestCountMerge(t *testing.T) {
	genInfos := func(counts []uint64) (infos []Info) {
		infos = make([]Info, len(counts))
		for i, x := range counts {
			infos[i] = Info{percent: num.NewPercent(x, 0)}
		}
		return
	}
	tests := []struct {
		name   string
		counts []uint64
		check  bool
	}{
		{"a", []uint64{1, 2, 1, 4, 5}, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		b := Backend{maxDiff: 0.1}
		t.Run(tt.name, func(t *testing.T) {
			got := b.CountMerge(genInfos(tt.counts))
			for i, x := range got {
				fmt.Printf("%v:", i)
				for _, y := range x {
					fmt.Printf("%v ", y.percent.Count())
				}
				fmt.Println("")
			}
			t.Error("check")
		})
	}
}
