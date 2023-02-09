package point

import (
	"regexp"
	"strconv"
	"strings"
)

type JsonPoint = string

var arrayReg = regexp.MustCompile(`\/[0-9]+(\/.*)?$`)

func Nth(x JsonPoint, idx int) JsonPoint {
	return x + "/" + strconv.Itoa(idx)
}

func Depth(x JsonPoint) int {
	return strings.Count(x, "/")
}

func Split(x JsonPoint) []string {
	return strings.Split(x, "/")
}

func Root(x JsonPoint) string {
	i := strings.IndexRune(x, '/')
	if i == -1 {
		return ""
	}
	return x[:i]
}

func AvgDepth(jps []JsonPoint) (res int) {
	l := len(jps)
	if l == 0 {
		return 0
	}
	for _, x := range jps {
		res += Depth(x)
	}
	res /= l
	return
}

func GetParent(x JsonPoint) JsonPoint {
	labelIdx := strings.LastIndex(x, "/")
	if labelIdx == -1 {
		return ""
	}
	return x[:labelIdx]
}

func CurrentPoint(x JsonPoint) string {
	labelIdx := strings.LastIndex(x, "/")
	if labelIdx == -1 {
		return x
	}
	return x[labelIdx+1:]
}

func ContainsArray(x JsonPoint) bool {
	return arrayReg.Match([]byte(x))
}

func IsArray(x JsonPoint) bool {
	return ContainsArray("/" + CurrentPoint(x))
}

// IsUnnestedArray checks if x contains and only contains one array
func IsUnnestedArray(x JsonPoint) bool {
	if IsArray(x) {
		pre, _, _ := ParseArray(x)
		return !ContainsArray(pre)
	}
	return false
}

// /entities/hashtags/14/text
// => /entities/hashtags, /text, 14
// ParseArray find the last array index and split it.
// ParseArray assumed that x contains array.
func ParseArray(x JsonPoint) (prefix, suffix JsonPoint, idx int) {
	s := strings.Split(x, "/")
	for i := len(s) - 1; i >= 0; i-- {
		if ContainsArray("/" + s[i]) {
			prefix = strings.Join(s[:i], "/")
			suffix = strings.Join(s[i+1:], "/")
			var err error
			idx, err = strconv.Atoi(s[i])
			if err != nil {
				panic(err)
			}
			return
		}
	}
	return
}

func HasSameParent(x, y JsonPoint) bool {
	return GetParent(x) == GetParent(y)
}

func HasSameParentSlice(s []JsonPoint) bool {
	var fst JsonPoint
	for i, x := range s {
		if i == 0 {
			fst = x
		} else {
			if !HasSameParent(fst, x) {
				return false
			}
		}
	}
	return true
}

// IsSemiEqual check if x and y both contains array, and the only difference between them is the index of this array.
func IsSemiEqual(x, y JsonPoint) bool {
	if ContainsArray(x) && ContainsArray(y) {
		xPre, xSuf, _ := ParseArray(x)
		yPre, ySuf, _ := ParseArray(y)
		return xPre == yPre && xSuf == ySuf
	}
	return false
}

type relation int

const (
	None relation = iota
	Child
	Parent
	SameParent
	SemiEqual
	Equal
)

func Relation(x, y JsonPoint) relation {
	switch {
	case x == y:
		return Equal
	case IsSemiEqual(x, y):
		return SemiEqual
	case strings.HasPrefix(x, y):
		return Child
	case strings.HasPrefix(y, x):
		return Parent
	case HasSameParent(x, y):
		return SameParent
	default:
		return None
	}
}
