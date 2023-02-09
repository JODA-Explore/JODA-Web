package rating

import (
	"fmt"
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/point"
)

type Path struct {
	jsonPoint  string
	depthScale float64
	reversed   bool
	result     float64
}

func PathFactor(depthScale float64, reversed bool, jsonPoint string) *Path {
	return &Path{jsonPoint: jsonPoint, depthScale: depthScale, reversed: reversed, result: -101}
}

var pathSpline = Spline{
	X:     []float64{1, 2, 3, 5, 8, 10, 20},
	Y:     []float64{2, 1, 0, -1, -3, -8, -10},
	DescX: "Depth",
	DescY: "Scale",
	Name:  "Path Depth",
}

func (p *Path) ToScale() float64 {
	if p.result == -101 {
		depth := point.Depth(p.jsonPoint)
		var spline Spline
		if p.reversed {
			spline = reverse(pathSpline)
		} else {
			spline = pathSpline
		}
		p.result = spline.Cal(float64(depth), p.depthScale)
	}
	return p.result
}

func (p *Path) Name() string {
	return "Path Depth Factor"
}

func (p *Path) Desc() template.HTML {
	return desc.Collapse(
		desc.Scale(p.Name(), p.ToScale()),
		template.HTML(fmt.Sprintf(
			"Json Point: <b>%v</b> <br>Depth: <b>%v</b>, Scale:%.2f",
			p.jsonPoint,
			point.Depth(p.jsonPoint),
			p.depthScale,
		)),
		pathSpline.ApplyScale(p.depthScale).Desc(),
		"",
	)
}
