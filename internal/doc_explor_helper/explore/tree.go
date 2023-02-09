package explore

import (
	"fmt"
)

type DatasetTree struct {
	Dataset  string
	Children []*DatasetTree
}

func (tr *DatasetTree) find(dataset string) (*DatasetTree, bool) {
	if tr.Dataset == dataset {
		return tr, true
	}
	for _, x := range tr.Children {
		if res, found := x.find(dataset); found {
			return res, found
		}
	}
	return nil, false
}

func (root *DatasetTree) add(parentName, childName string) error {
	parent, foundParent := root.find(parentName)
	if !foundParent {
		return fmt.Errorf("can not find parent: %v for the child: %v", parentName, childName)
	}
	if _, foundChild := root.find(childName); foundChild {
		return fmt.Errorf("The dataset: [%v] already exists", childName)
	}
	child := DatasetTree{Dataset: childName}
	parent.Children = append(parent.Children, &child)
	// fmt.Printf("parent:%v\n", parent)
	// fmt.Printf("parent.Children:%v\n", parent.Children)
	return nil
}

