package freq

// import (
// 	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/find"
// 	"github.com/valyala/fastjson"
// )

// type Json []*fastjson.Value

// func (j Json) TopN(n int) List {
// 	filtered := find.TopN(n, j, func(v *fastjson.Value) uint64 {
// 		return v.GetUint64("count")
// 	})
// 	list := make(List, len(filtered))
// 	for i, x := range filtered {
// 		list[i] = NewPair(string(x.GetStringBytes("group")), x.GetUint64("count"))
// 	}
// 	return list
// }
