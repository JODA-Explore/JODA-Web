package distinct

import (
	"fmt"
	"strings"

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
	topNumber int
	feature.Limit
}

func New(dataset *query.Dataset) backend.Backend {
	return feature.NewBackend[Info](Backend{dataset, 3, feature.NewLimit(6, 0)})
}

func (b Backend) ID() backends.ID {
	return backends.Distinct
}

func (b Backend) Filter() points.Walker {
	return func(path string, v *points.Value, depth int) trie.Control {
		if v.Info.CountArray > 0 {
			return trie.Continue
		}
		if v.Info.CountString == 0 && v.Info.CountNumber == 0 {
			return trie.Next
		}
		return trie.None
	}
}

func (b Backend) Infos(path string, pi points.Info) ([]Info, trie.Control) {
	info := Info{BasicInfo: result.NewBasicInfo(path, b.ID())}
	total := pi.CountTotal - pi.CountNull
	switch {
	case pi.CountString == total:
		info.valueType = jsontype.String
		info.min = pi.MinStrSize
		info.max = pi.MaxStrSize
	case pi.CountInt == total:
		info.valueType = jsontype.Int
		info.min = pi.MinInt
		info.max = pi.MaxInt
	case pi.CountFloat == total:
		info.valueType = jsontype.Float
	default:
		return nil, trie.Next
	}
	info.topPercent.SetTotal(total)
	info.coverage = num.NewPercent(total, b.TotalNumber())
	return []Info{info}, trie.None
}

func (b Backend) complete(info *Info, isSample bool) {
	var source string
	if isSample {
		source = b.Sample()
	} else {
		source = b.Source()
	}
	topk, err := b.DistinctValuesFreq(source, info.JsonPoint(), b.topNumber)
	if err != nil {
		panic(fmt.Errorf("Can not get DistinctValuesFreq: %w", err))
	}
	var topSum uint64
	if isSample {
		info.sampleTop = topk
		for _, top := range topk.List() {
			top.Val = strings.ReplaceAll(top.Val, `"`, `\"`)
			count, err := b.Ins.ValueCount(b.Source(), info.JsonPoint(), top.Val, info.valueType)
			if err != nil {
				// something wrong happened, skip this info
				info.topPercent = num.NewPercent(0, 0)
				info.updateRating()
				return
			}
			topSum += count
		}
	} else {
		for _, x := range topk.List() {
			topSum += x.N
		}
	}
	info.topPercent.SetCount(topSum)
	info.updateRating()
}
