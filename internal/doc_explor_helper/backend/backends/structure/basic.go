package structure

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
)

type Backend struct {
	*query.Dataset
	maxDiff float64
	feature.Limit
}

func New(dataset *query.Dataset, maxDiff float64) backend.Backend {
	return feature.NewBackend[Info](Backend{dataset, maxDiff, feature.NewLimit(5, 400)})
}

func (b Backend) ID() backends.ID {
	return backends.StructureDiff
}

func (b Backend) Filter() points.Walker {
	return func(path string, v *points.Value, depth int) (ctl trie.Control) {
		if v.Info.CountArray > 0 {
			return trie.Continue
		}
		pi := v.Info
		if pi.CountNull*2 > pi.CountTotal {
			return trie.Continue
		}
		if pi.CountTotal == 0 || pi.CountTotal == b.TotalNumber() {
			return trie.Next
		}
		return trie.None
	}
}

func (b Backend) newPercent(count uint64) num.Percent {
	return num.NewPercent(count, b.TotalNumber())
}

func (b Backend) newInfo(path string, pi points.Info, totalNumber uint64) (info Info) {
	basic := result.NewBasicInfo(path, backends.StructureDiff)
	if 2*pi.CountTotal > totalNumber {
		info = Info{
			BasicInfo: basic,
			conds:     cond.Conditions{cond.New(path, true, cond.Exists, nil)},
			percent:   b.newPercent(totalNumber - pi.CountTotal),
		}
	} else {
		info = Info{
			BasicInfo: basic,
			conds:     cond.Conditions{cond.New(path, false, cond.Exists, nil)},
			percent:   b.newPercent(pi.CountTotal),
		}
	}
	info.updateRating()
	return
}

func (b Backend) Infos(path string, pi points.Info) (infos []Info, ctl trie.Control) {
	info := b.newInfo(path, pi, b.TotalNumber())
	return []Info{info}, trie.None
}
