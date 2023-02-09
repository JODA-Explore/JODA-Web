package objects

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
)

type Backend struct {
	*query.Dataset
}

func New(dataset *query.Dataset) backend.Backend {
	return feature.NewBackend[Info](Backend{dataset})
}

func (b Backend) ID() backends.ID {
	return backends.Objects
}

func (b Backend) Filter() points.Walker {
	return func(_ string, v *points.Value, _ int) trie.Control {
		if v.Info.CountArray == 0 {
			return trie.Next
		}
		return trie.None
	}
}

func (b Backend) Infos(path string, pi points.Info) (infos []Info, ctl trie.Control) {
	ctl = trie.Continue
	cur, _ := b.Tree().Find(path)
	arrayType := jsontype.Object
	var total uint64
	cur.Walk(func(path string, v *points.Value, depth int) trie.Control {
		// root, i.e. cur self
		if path == "" {
			return trie.Next
		}
		if v.Info.CountObject != v.Info.CountTotal {
			arrayType = jsontype.Invalid
			return trie.Break
		}
		total += v.Info.CountTotal
		return trie.Continue
	})
	if arrayType == jsontype.Invalid || total == 0 {
		return
	}
	info := Info{
		BasicInfo: result.NewBasicInfo(path, 4),
		real:      num.NewPercent(0, total),
		coverage:  num.NewPercent(total, b.TotalNumber()),
	}
	if b.UseSample(info) {
		return []Info{info}, ctl
	}
	realCount, realTotal, err := b.DistinctObjCounts(b.Source(), path)
	if err != nil {
		panic(fmt.Errorf("calculate real DistinctObjCounts of %v:%w", path, err))
	}
	if realTotal != total {
		fmt.Printf("realTotal:%v,total:%v\n", realTotal, total)
	}
	info.real.SetCount(realCount)
	info.updateRating()
	infos = []Info{info}
	return
}
