package attstats

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/point"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
	"github.com/valyala/fastjson"
)

type Indices []int

type Attstats struct {
	list     []AttStat
	hashmap  map[string]int
	All      Indices
	NonArray Indices
}

func (res *Attstats) parse(v *fastjson.Value, name string) {
	newName := name
	if n := string(v.GetStringBytes("Key")); n != "" {
		newName += "/" + n
	}
	as := AttStat{
		JsonPoint:    newName,
		CountTotal:   v.GetUint64("Count_Total"),
		CountObject:  v.GetUint64("Count_Object"),
		CountArray:   v.GetUint64("Count_Array"),
		CountNull:    v.GetUint64("Count_Null"),
		CountNumber:  v.GetUint64("Count_Number"),
		CountInt:     v.GetUint64("Count_Int"),
		CountFloat:   v.GetUint64("Count_Float"),
		CountBoolean: v.GetUint64("Count_Boolean"),
		CountTrue:    v.GetUint64("Count_True"),
		CountFalse:   v.GetUint64("Count_False"),
		CountString:  v.GetUint64("Count_String"),
		MaxMember:    v.GetUint64("Max_Member"),
		MaxStrSize:   v.GetUint64("Max_StrSize"),
		MaxSize:      v.GetUint64("Max_Size"),
		MaxInt:       v.GetUint64("Max_Int"),
		MinMember:    v.GetUint64("Min_Member"),
		MinStrSize:   v.GetUint64("Min_StrSize"),
		MinSize:      v.GetUint64("Min_Size"),
		MinInt:       v.GetUint64("Min_Int"),
	}
	res.list = append(res.list, as)
	res.hashmap[newName] = len(res.list) - 1
	if !point.ContainsArray(newName) {
		res.NonArray = append(res.NonArray, len(res.list)-1)
	}
	for _, x := range v.GetArray("Array_Items") {
		res.parse(x, newName)
	}
	for _, x := range v.GetArray("Children") {
		res.parse(x, newName)
	}
}

func New(v *fastjson.Value) (res Attstats) {
	res.hashmap = make(map[string]int)
	res.parse(v, "")
	all := make(Indices, len(res.list))
	for i := range all {
		all[i] = i
	}
	res.All = all
	return
}

func Extract[T any](as *Attstats, indices Indices, f func(AttStat) T) []T {
	l := as.list
	res := make([]T, len(indices))
	for i, x := range indices {
		as := l[x]
		res[i] = f(as)
	}
	return res
}

func FilterByIndices(as *Attstats, indices Indices) []AttStat {
	res := make([]AttStat, len(indices))
	for i, x := range indices {
		res[i] = as.list[x]
	}
	return res
}

func Filter(as *Attstats, indices Indices, keep func(AttStat) bool) (res []AttStat) {
	// slices.Filter is in-place, so the indices should first be copied.
	new := make([]int, len(indices))
	copy(new, indices)
	slices.Filter(&new, func(i int) bool {
		return keep(as.list[i])
	})
	return FilterByIndices(as, new)
}

func Iter(as *Attstats, indices Indices, f func(AttStat)) {
	for _, x := range indices {
		f(as.list[x])
	}
}

func (as *Attstats) Keys(indices Indices) []string {
	return Extract(as, indices, func(as AttStat) string {
		return as.JsonPoint
	})
}

func (as *Attstats) Find(point string) AttStat {
	idx := as.hashmap[point]
	return as.list[idx]
}
