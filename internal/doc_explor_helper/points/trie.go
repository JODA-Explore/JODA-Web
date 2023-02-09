package points

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
	"github.com/valyala/fastjson"
)

type (
	Trie   = trie.Trie[*Value]
	Walker func(path string, v *Value, depth int) trie.Control
)

func empty() Trie {
	return trie.Empty[*Value]()
}

func parse(t *Trie, v *fastjson.Value) {
	parseChild := func(x *fastjson.Value) {
		name := string(x.GetStringBytes("Key"))
		child := trie.Empty[*Value]()
		child.SetEnd()
		parse(&child, x)
		t.AddChild(name, &child)
	}
	info := Info{
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
	t.SetVal(&Value{Info: info})
	array := v.GetArray("Array_Items")
	children := v.GetArray("Children")
	t.InitLength(len(array) + len(children))
	for _, x := range array {
		parseChild(x)
	}
	for _, x := range children {
		parseChild(x)
	}
}

func New(v *fastjson.Value) *Trie {
	res := trie.Empty[*Value]()
	parse(&res, v)
	return &res
}
