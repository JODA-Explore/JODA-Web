package points

// import (
// 	"strconv"
// 	"strings"
// 	"testing"

// 	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/query"
// 	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
// 	"github.com/JODA-Explore/JODA-Web/internal/joda"

// 	"github.com/valyala/fastjson"
// )

// func TestNew(t *testing.T) {
// 	ji := joda.New("http://localhost:5632")
// 	ins := query.NewIns(ji)
// 	twitter, err := ins.GetAttribute("twitter")
// 	var sb strings.Builder
// 	if err != nil {
// 		panic(err)
// 	}
// 	tests := []struct {
// 		name   string
// 		v      *fastjson.Value
// 		walker Walker
// 		check  bool
// 	}{
// 		{
// 			"level 1", twitter,
// 			func(path string, v *Value, depth int) trie.Control {
// 				if depth > 1 {
// 					return trie.Continue
// 				}
// 				sb.WriteString(path)
// 				sb.WriteString("\n")
// 				return trie.Next
// 			},
// 			false,
// 		},
// 		{
// 			"count total", twitter,
// 			func(path string, v *Value, depth int) trie.Control {
// 				if depth > 1 {
// 					return trie.Continue
// 				}
// 				sb.WriteString(path)
// 				sb.WriteString(":")
// 				sb.WriteString(strconv.FormatUint(v.Info.CountTotal, 10))
// 				sb.WriteString("\n")
// 				return trie.Next
// 			},
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		if !tt.check {
// 			continue
// 		}
// 		t.Run(tt.name, func(t *testing.T) {
// 			sb.Reset()
// 			got := New(tt.v)
// 			got.Walk(tt.walker)
// 			t.Errorf("%v", sb.String())
// 		})
// 	}
// }

// func BenchmarkNew(b *testing.B) {
// 	ji := joda.New("http://localhost:5632")
// 	ins := query.NewIns(ji)
// 	twitter, _ := ins.GetAttribute("twitter")
// 	for i := 0; i < 100; i++ {
// 		New(twitter)
// 	}
// }
