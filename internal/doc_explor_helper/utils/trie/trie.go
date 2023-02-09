package trie

import (
	"fmt"
	"strings"
)

type Trie[T any] struct {
	val      T
	children map[string]*Trie[T]
	end      bool
}

func Empty[T any]() (res Trie[T]) {
	return
}

func New[T any](l int) (res Trie[T]) {
	res.children = make(map[string]*Trie[T], l)
	return
}

// based on https://github.com/dghubble/trie/blob/master/common.go
func segment(p string, start int) (segment string, next int) {
	if len(p) == 0 || start < 0 || start > len(p)-1 {
		return "", -1
	}
	end := strings.IndexRune(p[start+1:], '/') // next '/' after 0th rune
	if end == -1 {
		return p[1+start:], -1
	}
	return p[1+start : start+end+1], start + end + 1
}

func (tn *Trie[T]) InitLength(l int) {
	tn.children = make(map[string]*Trie[T], l)
}

func (tn *Trie[T]) Val() T {
	return tn.val
}

func (tn *Trie[T]) SetVal(v T) {
	tn.val = v
}

func (tn *Trie[T]) AddChild(key string, child *Trie[T]) {
	tn.children[key] = child
}

func (tn *Trie[T]) Len() int {
	return len(tn.children)
}

func (tn *Trie[T]) Insert(p string, v T) {
	cur := tn
	seg, i := "", 0
	for {
		seg, i = segment(p, i)
		if seg == "" {
			break
		}
		if child, ok := cur.children[seg]; ok {
			cur = child
		} else {
			new := New[T](0)
			cur.children[seg] = &new
			cur = &new
		}
	}
	cur.end = true
	cur.val = v
}

func (tn *Trie[T]) Find(p string) (target *Trie[T], found bool) {
	walker := func(path string, v T, depth int) Control {
		return None
	}
	_, target, found = tn.PathWalk(p, walker)
	return
}

func (tn *Trie[T]) Delete() {
	tn.end = false
}

func (tn *Trie[T]) SetEnd() {
	tn.end = true
}

func (t Trie[T]) String() string {
	var sb strings.Builder
	t.Walk(func(path string, v T, depth int) Control {
		sb.WriteString(strings.Repeat(" ", 4*depth))
		sb.WriteString(path)
		sb.WriteString("\n")
		return Next
	})
	return sb.String()
}

func (t Trie[T]) StringWithValue(fn func(T) any) string {
	var sb strings.Builder
	t.Walk(func(path string, v T, depth int) Control {
		sb.WriteString(strings.Repeat(" ", 4*depth))
		sb.WriteString(path)
		sb.WriteString(fmt.Sprintf("(%v)", fn(v)))
		sb.WriteString("\n")
		return Next
	})
	return sb.String()
}
