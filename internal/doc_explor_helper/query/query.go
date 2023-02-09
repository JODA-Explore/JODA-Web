package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/freq"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/heap"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"

	"github.com/JODA-Explore/JODA-Web/internal/joda"

	"github.com/JODA-Explore/JODA-Web/internal/source"
	"github.com/valyala/fastjson"
)

type Ins struct {
	*joda.Joda
}

func NewIns(j joda.Joda) Ins {
	return Ins{&j}
}

// var JodaTime float64

func (ins Ins) Execute(query string) error {
	result, _, err := ins.ExecuteQuery(query)
	if err == nil {
		err = ins.RemoveSource(source.Source{ID: result})
	}
	// if b != nil {
	// 	JodaTime += b.Runtime.Query
	// }
	return err
}

func (ins Ins) getResult(query string) (v *fastjson.Value, err error) {
	result, _, err := ins.ExecuteQuery(query)
	if err != nil {
		return nil, fmt.Errorf("Execute query:%w", err)
	}
	// if b != nil {
	// 	JodaTime += b.Runtime.Query
	// }
	defer ins.RemoveSource(source.Source{ID: result})
	v, err = ins.GetResultDocumentsFastJson(result, 0, 1)
	if err != nil {
		return v, fmt.Errorf("Get result of query:%v\n%w", query, err)
	}
	array, err := v.Array()
	if err != nil {
		return v, fmt.Errorf("Convert result of query %v to array:\n%w", query, err)
	}
	if len(array) != 1 {
		return nil, fmt.Errorf("Expected one result document for query:%v", query)
	}
	return array[0], nil
}

func (ins Ins) AllDatasets() ([]string, error) {
	sources, err := ins.GetSources()
	if err != nil {
		return nil, fmt.Errorf("GetSources:%w", err)
	}
	sourceNames := make([]string, len(sources))
	for i, x := range sources {
		sourceNames[i] = x.Name
	}
	return sourceNames, nil
}

func (ins Ins) FilteredDatasets() ([]string, error) {
	sources, err := ins.GetSources()
	if err != nil {
		return nil, fmt.Errorf("GetSources:%w", err)
	}
	sourceNames := make([]string, 0, len(sources))
	for _, x := range sources {
		if !strings.HasPrefix(x.Name, "__AUTO_GENERATED__") {
			sourceNames = append(sourceNames, x.Name)
		}
	}
	return sourceNames, nil
}

func (ins Ins) GetAttribute(sourceName string) (v *fastjson.Value, err error) {
	q := cmd.Load(sourceName) + cmd.Agg("", cmd.Fun("ATTSTAT", `''`))
	return ins.getResult(q)
}

func (ins Ins) TopDistinctValues(dataset, path string, topN int) (members []string, counts []uint64, err error) {
	q := cmd.Load(dataset) + cmd.Agg("", cmd.Group(`COUNT('')`, "count", cmd.Quote(path)))
	v, err := ins.getResult(q)
	if err != nil {
		return
	}
	va, err := v.Array()
	if err != nil {
		err = fmt.Errorf("Convert to array: %w", err)
		return
	}
	var list []heap.Elem[string, uint64]
	if topN > 0 {
		f := freq.New(topN)
		for _, x := range va {
			f.Insert(string(x.GetStringBytes("group")), x.GetUint64("count"))
		}
		f.SortReverse()
		list = f.List()
	} else {
		for _, x := range va {
			list = append(list, heap.NewElem(string(x.GetStringBytes("group")), x.GetUint64("count")))
		}
		sort.Slice(list, func(i, j int) bool { return list[i].N > list[j].N })
	}
	members = make([]string, len(list))
	counts = make([]uint64, len(list))
	for i, x := range list {
		members[i] = x.Val
		counts[i] = x.N
	}
	return
}

func (ins Ins) AllDistinctValues(dataset string, path string) (members []string, err error) {
	members, _, err = ins.TopDistinctValues(dataset, path, 0)
	return
}

func (ins Ins) Sample(dataset string) (name string, err error) {
	sv, err := ins.newStoreVar(
		"__sample__",
		dataset,
		cmd.Choose(cmd.Fun("MOD", cmd.Fun("ID"), "10")+"== 0"),
	)
	if err != nil {
		err = fmt.Errorf("try to create sample for %v:%w", dataset, err)
		return
	}
	return sv.name, nil
}

func (ins Ins) AllDistinctMembers(dataset string, path string) (members []string, err error) {
	q := cmd.Load(dataset) +
		cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(path))) +
		cmd.Agg("", cmd.Fun("DISTINCT", cmd.Quote("")))
	v, err := ins.getResult(q)
	if err != nil {
		err = fmt.Errorf("AllDistinctMembers:%w", err)
		return
	}
	va, err := v.Array()
	if err != nil {
		err = fmt.Errorf("Convert to array: %w", err)
		return
	}
	members = slices.Map(va, func(x *fastjson.Value) string {
		buf := x.GetStringBytes()
		if buf == nil {
			buf = x.MarshalTo(nil) // make sure that the got type is string
		}
		return string(buf)
	})
	return
}
