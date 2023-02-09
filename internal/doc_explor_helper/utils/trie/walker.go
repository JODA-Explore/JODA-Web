package trie

import "strings"

type Control int

const (
	None Control = iota
	Next
	Continue
	Break
)

func (tn *Trie[T]) PathWalk(p string, f func(path string, v T, depth int) Control) (curPath string, curT *Trie[T], done bool) {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	curPath = ""
	curT = tn
	if tn.end {
		ctl := f("", curT.val, 0)
		if ctl > Next {
			return
		}
	}
	seg, i := "", 0
	for {
		seg, i = segment(p, i)
		if seg == "" {
			break
		}
		curPath += "/" + seg
		if child, ok := curT.children[seg]; ok {
			curT = child
			if tn.end {
				ctl := f("", curT.val, 0)
				if ctl > Next {
					return
				}
			}
		} else {
			return
		}
	}
	done = true
	return
}

func (tn *Trie[T]) depthWalkHelper(path string, depth int, f func(path string, v T, depth int) Control) (ctl Control) {
	if tn.end {
		ctl = f(path, tn.val, depth)
		if ctl > Next {
			return
		}
	}
	for seg, child := range tn.children {
		ctl = child.depthWalkHelper(path+"/"+seg, depth+1, f)
		if ctl == Break {
			return
		}
	}
	return Next
}

func (tn *Trie[T]) Walk(f func(path string, v T, depth int) Control) {
	tn.depthWalkHelper("", 0, f)
}

