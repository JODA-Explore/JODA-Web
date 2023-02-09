package explore

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends/array"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends/distinct"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends/objects"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends/structure"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends/types"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

// SearchOpt contails all the options that user can change to get result
// so SearchOpt is also the key of the map that stored QueriesInfo
type SearchOpt struct {
	Dataset         string
	Filter          filter.Filter
	EnabledBackends [backends.Num]bool
}

// resultsInfo is the search result
type resultsInfo struct {
	Results  result.Results
	Backends [backends.Num]backend.Backend // backends that generated the queries
}

// DatasetInfo contains all informations that are relative to dataset
type DatasetInfo struct {
	query.Dataset
	Opts  []*SearchOpt // all search options with the dataset
	Query string       // the query that generated the dataset
	Tree  *DatasetTree
	Pts   *points.Trie
}

func createBackend(dataset *query.Dataset, pts *points.Trie, backendId backends.ID, f filter.Filter) backend.Backend {
	switch backendId {
	case backends.ValueType:
		return types.New(dataset, f.MaxDiff)
	case backends.StructureDiff:
		return structure.New(dataset, f.MaxDiff)
	case backends.Distinct:
		return distinct.New(dataset)
	case backends.Array:
		return array.New(dataset)
	case backends.Objects:
		return objects.New(dataset)
	default:
		panic(fmt.Errorf("the backend %v can not be created", backendId))
	}
}

// newInfo create new DatasetInfo with empty opts,query and tree.
func newInfo(ins query.Ins, dataset string) (res DatasetInfo, err error) {
	attribute, err := ins.GetAttribute(dataset)
	if err != nil {
		return
	}
	res.Pts = points.New(attribute)
	res.Dataset = ins.NewDataset(dataset, res.Pts)
	res.Tree = &DatasetTree{Dataset: dataset}
	return
}

// newDataset generate new dataset and load it in Joda if the query is not empty.
func newDataset(ins query.Ins, opt SearchOpt, query string) error {
	if query != "" {
		err := ins.Execute(query)
		if err != nil {
			return err
		}
	}
	datasetInfo, err := newInfo(ins, opt.Dataset)
	if err != nil {
		return err
	}
	datasetInfo.Query = query
	datasetInfo.AddOpt(&opt)
	datasets[opt.Dataset] = &datasetInfo
	return nil
}

// AddOpt add search option to the DatasetInfo if its not yet exists.
func (di *DatasetInfo) AddOpt(opt *SearchOpt) {
	if _, exists := allResults[*opt]; !exists {
		di.Opts = append(di.Opts, opt)
	}
}

// EnabledBackendsName returns all enabled backends name
func (opt SearchOpt) EnabledBackendsName() (res []string) {
	for i, x := range opt.EnabledBackends {
		if x {
			res = append(res, backends.ID(i).Name())
		}
	}
	return
}

// GetResults calculate the queries if the result not exists and return it.
func Results(ins query.Ins, searchOpt SearchOpt) (info *resultsInfo, err error) {
	info, ok := allResults[searchOpt]
	// ok = false // for debug
	if !ok {
		datasetInfo := GetDatasetInfo(searchOpt.Dataset)
		resultsInfo := resultsInfo{}
		var results result.Results
		// fmt.Printf("joda time:%v\n", query.JodaTime)
		for i, enabled := range searchOpt.EnabledBackends {
			if enabled {
				// startTime := time.Now()
				// query.JodaTime = 0
				// fmt.Printf("%v start at %v\n", backends.ID(i).Name(), startTime.Format("15:04:05 02.01.2006"))
				be := createBackend(&datasetInfo.Dataset, datasetInfo.Pts, backends.ID(i), searchOpt.Filter)
				resultsInfo.Backends[i] = be
				new, err := be.Results(datasetInfo.Pts, searchOpt.Filter)
				if err != nil {
					return nil, err
				}
				// fmt.Printf("%v took time: %v\n", backends.ID(i).Name(), time.Since(startTime).Round(time.Millisecond))
				// fmt.Printf("joda time:%v\n", query.JodaTime)
				results.MergeWithLimit(searchOpt.Filter.MaxNumber, new)
			}
		}
		for i, r := range results {
			be := resultsInfo.Backends[r.Backend()]
			if be == nil {
				err = fmt.Errorf("can not find the backend id: %s(%v)", r.Backend().Name(), r.Backend())
				return
			}
			be.Desc(r, i)
		}
		resultsInfo.Results = results
		info = &resultsInfo
		allResults[searchOpt] = &resultsInfo
	}
	return info, nil
}
