package structure

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/point"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
)

var _ feature.CountMerger[Info] = Backend{} // Verify Merger is implemented.

func (b Backend) TryMerge(x, y Info) bool {
	conds := append(x.conds, y.conds...)
	count := b.CountMust(conds.Query(b.Source()))
	if num.Diff(count, x.percent.Count()) > b.maxDiff {
		return false
	}
	return true
}

type trieV struct {
	ctl trie.Control
	Info
}

func (b Backend) CountMerge(infos []Info) (res [][]Info) {
	tr := trie.New[*trieV](0)
	for _, x := range infos {
		tr.Insert(x.JsonPoint(), &trieV{Info: x})
	}
	disableChildren(&tr, b.maxDiff)
	return groupCounts(&tr, b.maxDiff)
}

func disableChildren(tr *trie.Trie[*trieV], maxDiff float64) {
	tr.Walk(func(path string, curV *trieV, depth int) (curCtl trie.Control) {
		curCtl = curV.ctl
		if curCtl != trie.None {
			return
		}
		cur, ok := tr.Find(path)
		if !ok {
			return trie.Next
		}
		cur.Walk(func(childPath string, v *trieV, _ int) (ctl trie.Control) {
			if childPath == "" {
				return trie.Next
			}
			ctl = v.ctl
			if ctl != trie.None {
				return
			}
			if num.Diff(curV.percent.Count(), v.percent.Count()) > maxDiff {
				return trie.Continue
			} else {
				v.ctl = trie.Next
				return trie.Next
			}
		})
		return trie.Next
	})
}

func groupCounts(tr *trie.Trie[*trieV], maxDiff float64) (res [][]Info) {
	tr.Walk(func(curPath string, curV *trieV, _ int) (curCtl trie.Control) {
		curCtl = curV.ctl
		if curCtl != trie.None {
			return
		}
		infos := []Info{curV.Info}
		tr.Walk(func(path string, v *trieV, _ int) (ctl trie.Control) {
			ctl = v.ctl
			if ctl != trie.None {
				return
			}
			rel := point.Relation(curPath, path)
			if rel == point.Child || rel == point.Parent || rel == point.Equal {
				return trie.Continue
			}
			if num.Diff(curV.percent.Count(), v.percent.Count()) <= maxDiff {
				infos = append(infos, v.Info)
				v.ctl = trie.Next
				return trie.Next
			}
			if curV.percent.Count() < v.percent.Count() {
				return trie.Next
			} else {
				return trie.Continue
			}
		})
		res = append(res, infos)
		curV.ctl = trie.Next
		return trie.Next
	})
	return
}
