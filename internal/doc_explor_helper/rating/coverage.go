package rating

import (
	"fmt"
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

var Coverage = Spline{
	X:     []float64{0.0, 0.10, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9, 1},
	Y:     []float64{0.6, 0.8, 0.85, 0.8, 0.7, 0.65, 0.6, 0.45, 0.3, 0},
	DescX: "Percent",
	DescY: "Score",
	Name:  "Coverage",
}

func NewCoverage(x float64, fs ...Factor) Rating {
	return Rating{
		X:       x,
		spline:  Coverage,
		Factors: fs,
	}
}

var reversedCoverage = Spline{
	X:     []float64{0.0, 0.10, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9, 1},
	Y:     slices.NewReversed(Coverage.Y),
	DescX: "Percent",
	DescY: "Score",
	Name:  "Reversed Coverage",
}

func NewReversedCoverage(x float64, fs ...Factor) Rating {
	return Rating{
		X:       x,
		spline:  reversedCoverage,
		Factors: fs,
	}
}

var coverageFactorSpline = Spline{
	X:     []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	Y:     []float64{-50, -30, -5, 0, 5, 8},
	DescX: "Percent",
	DescY: "Score",
	Name:  "Coverage",
}

type CoverageFactor struct {
	Val float64
}

func (c CoverageFactor) ToScale() float64 {
	return coverageFactorSpline.Cal(c.Val, 0)
}

func (c CoverageFactor) Desc() template.HTML {
	return desc.Collapse(
		desc.Scale(c.Name(), c.ToScale()),
		template.HTML(fmt.Sprintf("The coverage of all objects: %.2f%%", c.Val*100)),
		coverageFactorSpline.Desc(),
		"",
	)
}

func (c CoverageFactor) Name() string {
	return "Coverage Factor"
}
