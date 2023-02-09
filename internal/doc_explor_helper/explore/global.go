package explore

import (
	"errors"
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
)

var (
	allResults = make(map[SearchOpt]*resultsInfo)
	datasets   = make(map[string]*DatasetInfo)
	Current    SearchOpt
	// BackendsName = anlys.AllBackendsName()
)

func init() {
	root := DatasetInfo{}
	root.Tree = &DatasetTree{}
	datasets[""] = &root
}

func GetOpt(dataset string, idx int) *SearchOpt {
	return datasets[dataset].Opts[idx]
}

// addChild find the tree of the parent, and add child to it.
func addChild(parentName, childName string) error {
	root := datasets[""].Tree
	return root.add(parentName, childName)
}

func GetDatasetInfo(dataset string) *DatasetInfo {
	return datasets[dataset]
}

func CurrentResults() *resultsInfo {
	return allResults[Current]
}

// HandleNewDataset receives the user inputs and load the new dataset if needed.
func HandleNewDataset(
	ins query.Ins,
	sourceType, parent, queryContent string,
	searchOpt SearchOpt,
) error {
	var query string
	addNewDataset := func() (err error) {
		err = newDataset(ins, searchOpt, query)
		if err != nil {
			return
		}
		err = addChild(parent, searchOpt.Dataset)
		return
	}
	switch sourceType {
	case "exists":
		if datasetInfo, exists := datasets[searchOpt.Dataset]; exists {
			datasetInfo.AddOpt(&searchOpt)
		} else {
			err := addNewDataset()
			if err != nil {
				return err
			}
		}
	case "new":
		query = cmd.Load(parent) + queryContent + cmd.Store(searchOpt.Dataset)
		err := addNewDataset()
		if err != nil {
			return err
		}
	case "load":
		query = cmd.Load(searchOpt.Dataset) + " FROM " + queryContent
		err := addNewDataset()
		if err != nil {
			return err
		}
	}
	Current = searchOpt
	return nil
}

// NewChildName return the next child name by its parent.
func NewChildName(parentName string) (string, error) {
	root := datasets[""].Tree
	parent, found := root.find(parentName)
	if !found {
		return "", errors.New("can not find dataset:" + parentName)
	}
	return parentName + "_" + strconv.Itoa(len(parent.Children)), nil
}

func CurrentResult(resultIdx int) *result.Result {
	resultsInfo := CurrentResults()
	return resultsInfo.Results[resultIdx]
}

func CurrentBackend(resultIdx int) backend.Backend {
	backendId := CurrentResult(resultIdx).Backend()
	return CurrentResults().Backends[backendId]
}

func CurrentQueryMaker(resultIdx int) (template.HTML, error) {
	be := CurrentBackend(resultIdx)
	return be.QueryMaker(CurrentResult(resultIdx), resultIdx)
}
