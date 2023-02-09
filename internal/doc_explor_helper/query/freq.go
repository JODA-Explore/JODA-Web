package query

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/freq"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
)

func (ins Ins) memberFreq(query string, topNumber int) (f freq.Freq, err error) {
	v, err := ins.getResult(query)
	if err != nil {
		err = fmt.Errorf("Member Freq:%w", err)
		return
	}
	va, err := v.Array()
	if err != nil {
		err = fmt.Errorf("Convert to array: %w", err)
		return
	}
	f = freq.New(topNumber)
	for _, x := range va {
		name := x.GetStringBytes("group")
		if name == nil {
			name = x.Get("group").MarshalTo(nil) // make sure that the got type is string
		}
		f.Insert(string(name), x.GetUint64("count"))
	}
	return
}

func (ins Ins) DistinctValuesFreq(dataset, jsonPoint string, topNumber int) (f freq.Freq, err error) {
	return ins.memberFreq(cmd.Load(dataset)+cmd.As("", cmd.Quote(jsonPoint))+cmd.Agg("", cmd.Group(`COUNT('')`, "count", cmd.Quote(""))), topNumber)
}

func (ins Ins) AtomicMemberFreq(
	dataset, jsonPoint string,
	topNumber int,
) (f freq.Freq, err error) {
	q := cmd.Load(dataset) +
		cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(jsonPoint))) +
		cmd.Agg("", cmd.Group(`COUNT('')`, "count", cmd.Quote("")))
	return ins.memberFreq(q, topNumber)
}

func (ins Ins) ArrayObjMemberFreq(dataset, jsonPoint string, topNumber int) (f freq.Freq, err error) {
	flattened, err := ins.newStoreVar(
		"__FLATTENED__",
		dataset,
		cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(jsonPoint))),
	)
	if err != nil {
		err = fmt.Errorf("try to flatten %v:%w", dataset, err)
		return
	}
	q := cmd.Load(flattened.name) + cmd.As("", `HASH('')`) + cmd.Agg("", cmd.Group(`COUNT('')`, "count", `''`))
	f, err = ins.memberFreq(q, topNumber)
	flattened.clear()
	if err != nil {
		err = fmt.Errorf("ArrayObjMemberFreq:%w", err)
	}
	return
}
