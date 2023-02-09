package types

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
	return feature.NewBackend[Info](Backend{dataset, maxDiff, feature.NewLimit(4.5, 100)})
}

func (b Backend) ID() backends.ID {
	return backends.ValueType
}

func (b Backend) Filter() points.Walker {
	return func(path string, v *points.Value, depth int) (ctl trie.Control) {
		if v.Info.CountArray > 0 {
			return trie.Continue
		}
		return trie.None
	}
}

func newConditions(p string, neg bool, condType cond.Type) (res cond.Conditions) {
	if condType == cond.Equal {
		var right string
		if neg {
			right = "false"
		} else {
			right = "true"
		}
		res = append(res, cond.New(p, false, condType, right))
		return
	}
	res = append(res, cond.New(p, neg, condType, nil))
	if neg {
		// exists && is not null
		res = append(res, cond.New(p, false, cond.Exists, nil))
	}
	return
}

func (b Backend) Infos(path string, pi points.Info) (infos []Info, ctl trie.Control) {
	null := pi.CountNull
	total := pi.CountTotal
	nonNull := pi.CountTotal - pi.CountNull
	boolean := pi.CountBoolean
	args := []struct {
		nom, denom uint64
		condType   cond.Type
		neg        bool
	}{
		{null, total, cond.IsNull, false},
		{nonNull, total, cond.IsNull, true},
		{pi.CountObject, nonNull, cond.IsObject, false},
		{pi.CountArray, nonNull, cond.IsArray, false},
		{pi.CountNumber, nonNull, cond.IsNumber, false},
		{boolean, nonNull, cond.IsBool, false},
		{pi.CountTrue, boolean, cond.Equal, false},
		{pi.CountFalse, boolean, cond.Equal, true},
		{pi.CountString, nonNull, cond.IsString, false},
	}
	for _, arg := range args {
		if arg.nom == 0 || arg.nom == arg.denom {
			continue
		}
		info := Info{
			BasicInfo: result.NewBasicInfo(path, b.ID()),
			percent:   num.NewPercent(arg.nom, arg.denom),
			conds:     newConditions(path, arg.neg, arg.condType),
		}
		info.updateRating()
		infos = append(infos, info)
	}
	return
}
