package query

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
)

type Dataset struct {
	name, sample string
	totalNumber  uint64
	*Ins
	tree *points.Trie
}

func (ins Ins) NewDataset(name string, tree *points.Trie) Dataset {
	return Dataset{name: name, Ins: &ins, tree: tree}
}

func (d Dataset) Source() string {
	return d.name
}

func (d Dataset) Tree() *points.Trie {
	if d.tree == nil {
		panic(fmt.Errorf("the trie of dataset %v is not set", d.Source()))
	}
	return d.tree
}

func (d Dataset) Sample() string {
	if d.sample == "" {
		var err error
		d.sample, err = d.Ins.Sample(d.name)
		if err != nil {
			panic(err)
		}
	}
	return d.sample
}

func (d Dataset) TotalNumber() uint64 {
	if d.totalNumber == 0 {
		var err error
		d.totalNumber, err = d.Ins.CountOfDataset(d.name)
		if err != nil {
			panic(err)
		}
	}
	return d.totalNumber
}
