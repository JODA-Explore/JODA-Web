package query

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/cmd"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/jsontype"
	"github.com/JODA-Explore/JODA-Web/internal/joda"
)

func (ins Ins) Count(query string) (uint64, error) {
	query += cmd.Agg("", cmd.Fun("COUNT", `''`))
	v, err := ins.getResult(query)
	if err != nil {
		return 0, fmt.Errorf("get count of query:%v\n%w", query, err)
	}
	return v.Uint64()
}

func (ins Ins) CountMust(query string) uint64 {
	res, err := ins.Count(query)
	if err != nil {
		panic(err)
	}
	return res
}

func (ins Ins) CountOfDataset(dataset string) (uint64, error) {
	return ins.Count(cmd.Load(dataset))
}

func (ins Ins) CountOfDistinctValues(dataset string, candidats []string) (counts []uint64, err error) {
	left1 := func(i int) string {
		return cmd.Quote("/" + strconv.Itoa(i))
	}
	right1 := func(i int) string {
		return cmd.Fun("DISTINCT", cmd.Quote(candidats[i]))
	}
	distinct, err := ins.newStoreVar(
		"__DistinctValues__",
		dataset,
		cmd.AggFunc(len(candidats), left1, right1),
	)
	if err != nil {
		return nil, fmt.Errorf("Store DistinctValues: %w", err)
	}

	// get distinct values count
	left2 := func(i int) string {
		return cmd.Quote("/" + strconv.Itoa(i))
	}
	right2 := func(i int) string {
		return cmd.Fun("SIZE", cmd.Quote("/"+strconv.Itoa(i)))
	}
	as := cmd.AsFunc(len(candidats), left2, right2)
	q2 := cmd.Load(distinct.name) + as
	v, err := ins.getResult(q2)
	if err != nil {
		return nil, fmt.Errorf("Get distinct values: %w", err)
	}
	va, err := v.Array()
	if err != nil {
		err = fmt.Errorf("Convert to array: %w", err)
		return
	}
	counts = make([]uint64, len(va))
	for i, x := range va {
		counts[i], err = x.Uint64()
		if err != nil {
			err = fmt.Errorf("convert %v to uint64:%w", x, err)
			return
		}
	}
	distinct.clear()
	return
}

func (ins Ins) DistinctObjCounts(dataset, jsonPoint string) (count, total uint64, err error) {
	flattened, err := ins.newStoreVar("__flattened__", dataset, cmd.As("", cmd.Fun("FLATTEN", cmd.Quote(jsonPoint))))
	if err != nil {
		if errors.Is(err, joda.EmptyResError) {
			return 0, 0, nil
		}
		err = fmt.Errorf("try to flatten %v:%w", dataset, err)
		return
	}
	total, err = ins.Count(cmd.Load(flattened.name))
	if err != nil {
		err = fmt.Errorf("try to count total members for %v:%w", dataset, err)
		return
	}
	distinct, err := ins.newStoreVar("__DISTINCT__", flattened.name, cmd.As("", `HASH('')`)+cmd.Agg("", `DISTINCT('')`))
	if err != nil {
		if errors.Is(err, joda.EmptyResError) {
			return 0, 0, nil
		}
		err = fmt.Errorf("try to hash and pick distinct values of  %v:%w", dataset, err)
		return
	}
	q := cmd.Load(distinct.name) + cmd.As("", cmd.Fun("SIZE", `''`))
	v, err := ins.getResult(q)
	if err != nil {
		err = fmt.Errorf("try to count distinct objs for %v:%w", dataset, err)
		return
	}
	flattened.clear()
	distinct.clear()
	count = v.GetUint64()
	return
}

func (ins Ins) MemberCount(dataset, jsonPoint, member string) (c uint64, err error) {
	q := cmd.Load(dataset) +
		cmd.Choose(cmd.Fun("IN", member, cmd.Quote(jsonPoint)))
	return ins.Count(q)
}

func (ins Ins) ValueCount(dataset, jsonPoint, value string, jsonType jsontype.Type) (c uint64, err error) {
	switch jsonType {
	case jsontype.Object:
		return 0, fmt.Errorf("can not calculate the count of object: %v", jsonPoint)
	case jsontype.String:
		value = cmd.DoubleQuote(value)
	default:
	}
	conds := cond.Conditions{cond.New(jsonPoint, false, cond.Equal, value)}
	return ins.Count(conds.Query(dataset))
}
