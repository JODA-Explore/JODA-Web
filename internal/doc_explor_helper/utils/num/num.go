package num

import (
	"math"
	"math/big"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/constraints"
)

func Min[T constraints.Number](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max[T constraints.Number](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}

func Diff(x, y uint64) float64 {
	if x > y {
		x, y = y, x
	}
	res := percent(x, y)
	return 1 - res
}

func PercentBig(x, y big.Int) (res float64) {
	var r big.Rat
	r.SetFrac(&x, &y)
	res, _ = r.Float64()
	return
}

func percent(x, y uint64) (res float64) {
	a, b := NewBigInt(x), NewBigInt(y)
	return PercentBig(a, b)
}

func MaxBigInt(x, y big.Int) big.Int {
	if x.Cmp(&y) < 0 {
		return y
	} else {
		return x
	}
}

func NewBigInt(x uint64) (r big.Int) {
	r.SetUint64(x)
	return
}

func Div(z uint64, x int) uint64 {
	bigZ := NewBigInt(z)
	var res big.Int
	res.Div(&bigZ, big.NewInt(int64(x)))
	return res.Uint64()
}

func Mul(x int, y float64) int {
	return int(math.Round(float64(x) * y))
}
