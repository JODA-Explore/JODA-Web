package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

type CountMerger[I result.Info] interface {
	Multier[I]

	// CountMerge group all infos by its count,
	// so that the elements in each sub-slice have same count.
	CountMerge([]I) [][]I

	// The given two infos x and y must already have same count,
	// which means they are in the same sub-slice of the result of CountMerge.
	// TryMerge check whether the merged info(x && y) still has same count.
	TryMerge(x, y I) bool
}

func countRs[I result.Info](cm CountMerger[I], topK topK[I], f filter.Filter) (rs result.Results) {
	all := cm.CountMerge(topK.Vals())
	// allNum, groupedNum, avg, max := 0, 0, 0, 0
	// defer func() {
	// 	fmt.Println("-------------------------------------------------")
	// 	fmt.Printf("allNum: %v, groupedNum: %v, avg: %v, max:%v\n", allNum, groupedNum, float64(avg)/float64(groupedNum), max)
	// }()
	for _, counts := range all {
		// for _, x := range grouped {
		// 	allNum += 1
		// 	if len(x) > 1 {
		// 		groupedNum += 1
		// 		avg += len(x)
		// 	}
		// 	if len(x) > max {
		// 		max = len(x)
		// 	}
		// }
		grouped := GroupInfo(counts, cm.TryMerge)
		for _, g := range grouped {
			if len(g) == 1 {
				rs.TryInsert(g[0], g[0].Rating(), f)
			} else {
				mixed := mixedResult[I](cm, g)
				if mixed.Rating().Score() >= f.MinRating {
					rs.InsertWithLimit(f.MaxNumber, mixed)
				}
			}
		}
	}
	return
}

func groupSingleInfo[I any](s []I, group func(x, y I) bool) (grouped, single, rest []I) {
	if len(s) <= 1 {
		return nil, s, nil
	}
	j := 0
	for i, x := range s {
		j = i + 1
		for j < len(s) {
			y := s[j]
			if group(x, y) {
				grouped = []I{x, y}
				break
			}
			j++
		}
		if len(grouped) != 0 {
			single = s[:i]
			rest = s[i+1 : j]
			break
		}
	}
	if len(grouped) == 0 {
		return nil, nil, s
	}
	for i := j + 1; i < len(s); i++ {
		new := s[i]
		if group(grouped[0], new) {
			grouped = append(grouped, new)
		} else {
			rest = append(rest, new)
		}
	}
	return
}

func GroupInfo[I any](s []I, group func(i, j I) bool) (res [][]I) {
	for {
		grouped, single, rest := groupSingleInfo(s, group)
		for _, x := range single {
			res = append(res, []I{x})
		}
		if len(grouped) != 0 {
			res = append(res, grouped)
		} else {
			for _, x := range rest {
				res = append(res, []I{x})
			}
			return
		}
		s = rest
	}
}

// brute-force version
// func GroupInfo[I any](ori []I, group func(i, j I) bool) (res [][]I) {
// 	for i, x := range ori {
// 		g := []I{x}
// 		for j := i + 1; j < len(ori); j++ {
// 			y := ori[j]
// 			if group(x, y) {
// 				g = append(g, y)
// 			}
// 		}
// 		res = append(res, g)
// 	}
// 	return
// }

type Counts[I result.Info] struct {
	Infos   [][]I
	counts  []uint64
	count   func(I) uint64
	maxDiff float64
}

func NewCounts[I result.Info](count func(I) uint64, maxDiff float64) Counts[I] {
	return Counts[I]{count: count, maxDiff: maxDiff}
}

func (cs *Counts[I]) Insert(info I) {
	ct := cs.count(info)
	for i, c := range cs.counts {
		if num.Diff(ct, c) <= cs.maxDiff {
			cs.Infos[i] = append(cs.Infos[i], info)
			return
		}
	}
	cs.counts = append(cs.counts, ct)
	cs.Infos = append(cs.Infos, []I{info})
}
