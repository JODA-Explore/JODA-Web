package distinct

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

var sizeSpline = rating.Spline{
	X:     []float64{0, 5, 10, 15, 20, 50},
	Y:     []float64{0, 0, -5, -10, -20, -30},
	DescX: "Size",
	DescY: "Scale",
	Name:  "Size",
}

type sizeFactor uint64

func (s sizeFactor) Desc() template.HTML {
	return desc.Collapse(
		desc.Scale(s.Name(), s.ToScale()),
		template.HTML(fmt.Sprintf("The Maximum Size: <b> %v </b>, smaller is better", s)),
		sizeSpline.Desc(),
		"",
	)
}

func (s sizeFactor) Name() string {
	return "Size Factor"
}

func (s sizeFactor) ToScale() float64 {
	return sizeSpline.Cal(float64(s), 0)
}

type minMaxDiffFactor struct {
	min, max uint64
}

var (
	minMaxAbsDiffSpline = rating.Spline{
		X:     []float64{0, 2, 5, 8, 10},
		Y:     []float64{20, 10, 0, -10, -20},
		DescX: "Diff",
		DescY: "Scale",
		Name:  "Abolute Diff",
	}
	minMaxRelativDiffSpline = rating.Spline{
		X:     []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		Y:     []float64{-30, -25, -20, -15, -10, 0, 5, 10, 15, 20, 30},
		DescX: "Diff",
		DescY: "Scale",
		Name:  "Relative Diff",
	}
	minMaxThreshold uint64 = 10
)

func (mmf minMaxDiffFactor) useAbsolute() bool {
	return mmf.max < minMaxThreshold
}

func (mmf minMaxDiffFactor) Desc() template.HTML {
	var block, d template.HTML
	if mmf.useAbsolute() {
		block = minMaxAbsDiffSpline.Desc()
		d = "Absolute Difference will be used since the max size is smaller than " + template.HTML(
			strconv.FormatUint(
				minMaxThreshold,
				10,
			),
		)
	} else {
		block = minMaxRelativDiffSpline.Desc()
		d = "Relative Difference will be used since the max size is greater than " + template.HTML(
			strconv.FormatUint(
				minMaxThreshold,
				10,
			),
		)
	}
	return desc.Collapse(
		desc.Scale(mmf.Name(), mmf.ToScale()),
		template.HTML(fmt.Sprintf("Min:<b>%v</b>,Max:<b>%v</b><br>", mmf.min, mmf.max))+d,
		block,
		"",
	)
}

func (mmf minMaxDiffFactor) Name() string {
	return "Min Max Diff Factor"
}

func (mmf minMaxDiffFactor) ToScale() float64 {
	if mmf.max < minMaxThreshold {
		return minMaxAbsDiffSpline.Cal(float64(mmf.max-mmf.min), 0)
	}
	p := num.NewPercent(mmf.min, mmf.max)
	return minMaxRelativDiffSpline.Cal(p.Ratio(), 0)
}

func newMinMaxFactor(min, max uint64) minMaxDiffFactor {
	return minMaxDiffFactor{min: min, max: max}
}
