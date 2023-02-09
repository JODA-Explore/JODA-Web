package array

import (
	"fmt"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
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
	return feature.NewBackend[Info](Backend{dataset, 3, feature.NewLimit(3, 0)})
}

func (b Backend) ID() backends.ID {
	return backends.Array
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
	arrayType := jsontype.Invalid
	var total uint64
	cur.Walk(func(path string, v *points.Value, depth int) trie.Control {
		getType := func() jsontype.Type {
			switch {
			case v.Info.CountString == v.Info.CountTotal:
				return jsontype.String
			case v.Info.CountInt == v.Info.CountTotal:
				return jsontype.Int
			case v.Info.CountFloat == v.Info.CountTotal:
				return jsontype.Float
			case v.Info.CountBoolean == v.Info.CountTotal:
				return jsontype.Bool
			default:
				return jsontype.Invalid
			}
		}
		// skip root, i.e. cur self
		if path == "" {
			return trie.Next
		}
		total += v.Info.CountTotal
		newType := getType()
		if arrayType == jsontype.Invalid {
			// init arrayType
			arrayType = newType
			return trie.Continue
		} else {
			if arrayType != newType {
				// if there are more than one type in the array, then the type is invalid
				arrayType = jsontype.Invalid
				return trie.Break
			} else {
				return trie.Continue
			}
		}
	})
	if arrayType == jsontype.Invalid {
		return
	}
	infos = []Info{{
		BasicInfo:  result.NewBasicInfo(path, b.ID()),
		arrayType:  arrayType,
		topPercent: num.NewPercent(0, total),
		coverage:   num.NewPercent(pi.CountTotal-pi.CountNull, b.TotalNumber()),
	}}
	return
}

func (b Backend) complete(info *Info, isSample bool) {
	var source string
	if isSample {
		source = b.Sample()
	} else {
		source = b.Source()
	}
	topK, err := b.AtomicMemberFreq(source, info.JsonPoint(), b.topNumber)
	if err != nil {
		panic(fmt.Errorf("Can not get StringArrayMemberFreq: %w", err))
	}
	var topSum uint64
	if isSample {
		info.sampleTop = topK
		for _, top := range topK.List() {
			if info.arrayType == jsontype.String {
				top.Val = strings.ReplaceAll(top.Val, `"`, `\"`)
				top.Val = cmd.DoubleQuote(top.Val)
			}
			count, err := b.MemberCount(b.Source(), info.JsonPoint(), top.Val)
			if err != nil {
				panic(fmt.Errorf("Can not get member count of %v: %w", top.Val, err))
			}
			topSum += count
		}
	} else {
		for _, top := range topK.List() {
			topSum += top.N
		}
	}
	info.topPercent.SetCount(topSum)
	info.updateRating()
}
