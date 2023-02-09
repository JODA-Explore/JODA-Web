package trie

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		name  string
		ps    []string
		check bool
	}{
		{"flat", []string{"a", "b", "c"}, false},
		{"deep", []string{"/a", "/a/b", "/c", "/c/d"}, false},
		{"order", []string{"/a/b", "/a", "/c", "/c/d"}, false},
		{"dup", []string{"/a/b", "/a", "/c", "/c/d", "/a/b"}, false},
		{"complex1", []string{"/a/b/c", "/abc", "/c", "/d", "/a", "/c", "/c/d", "/a/b", "/d/a"}, false},
		{"complex2", []string{"/a/b/c", "/abc", "/c", "/d", "/a", "/c", "/c/d", "/d/a"}, true},
	}
	for _, tt := range tests {
		if !tt.check {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			root := New[int](0)
			for i, x := range tt.ps {
				root.Insert(x, i)
			}
			t.Errorf("\n%v", root)
		})
	}
}

func TestSegment(t *testing.T) {
	tests := []struct {
		name        string
		p           string
		start       int
		wantSegment string
		wantNext    int
	}{
		{"1", "/a/b/c", 0, "a", 2},
		{"1", "a/b/c", 0, "", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeg, gotNext := segment(tt.p, tt.start)
			if tt.wantSegment != gotSeg || tt.wantNext != gotNext {
				t.Errorf("Got:%v,%v,Want:%v,%v", gotSeg, gotNext, tt.wantSegment, tt.wantNext)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		ps        []string
		key       string
		wantFound bool
	}{
		{"complex", []string{"/a/b/c", "/abc", "/c", "/d", "/a", "/c", "/c/d", "/d/a"}, "/a/b/c", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := New[int](0)
			for i, x := range tt.ps {
				root.Insert(x, i)
			}
			_, found := root.Find(tt.key)
			if found != tt.wantFound {
				t.Errorf("wantFound:%v, got:%v", tt.wantFound, found)
			}
		})
	}
}

func TestFind2(t *testing.T) {
	tests := []struct {
		name      string
		ps        []string
		key1      string
		key2      string
		wantFound bool
	}{
		{"complex", []string{"/a/b/c", "/abc", "/c", "/d", "/a", "/c", "/c/d", "/d/a"}, "/a/b", "/c", true},
		{"complex", []string{"/a/b/c", "/abc", "/c", "/d", "/a", "/c", "/c/d", "/d/a"}, "/a/b", "c", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := New[int](0)
			for i, x := range tt.ps {
				root.Insert(x, i)
			}
			child, ok := root.Find(tt.key1)
			if !ok {
				t.Errorf("can not find key:%v in root", tt.key1)
			}
			_, ok = child.Find(tt.key2)
			if ok != tt.wantFound {
				t.Errorf("wantFound:%v, got:%v", tt.wantFound, ok)
			}
		})
	}
}
