package rating

import (
	"bytes"
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
	"github.com/cnkei/gospline"
)

func applyScale(x, scale float64) float64 {
	return x * (1 + 0.01*scale)
}

type Spline struct {
	X, Y               []float64
	Name, DescX, DescY string
}

func (s Spline) at(x float64) float64 {
	return gospline.NewCubicSpline(s.X, s.Y).At(x)
}

func (s Spline) ApplyScale(scale float64) Spline {
	y := make([]float64, len(s.Y))
	for i, x := range s.Y {
		y[i] = applyScale(x, scale)
	}
	s.Y = y
	return s
}

func (s Spline) Cal(x, scale float64) float64 {
	if x > s.X[len(s.X)-1] {
		return applyScale(s.Y[len(s.Y)-1], scale)
	}
	return s.ApplyScale(scale).at(x)
}

var splineTmpl = template.Must(template.New("spline").Parse(`<table style="width:60%">
    <caption>{{.name}} Spline</caption>
    <tbody>
        <tr>
            {{range .x}}
            <td>{{.}}</td>
            {{end}}
        </tr>
        <tr>
            {{range .y}}
            <td>{{.}}</td>
            {{end}}
        </tr>
    </tbody>
</table>`))

func (s Spline) Desc() template.HTML {
	toString := func(l []float64, s string) []string {
		res := make([]string, len(l)+1)
		res[0] = s
		for i, x := range l {
			res[i+1] = string(desc.Float(x, 4))
		}
		return res
	}
	x, y := toString(s.X, s.DescX), toString(s.Y, s.DescY)
	data := make(map[string]interface{})
	data["x"] = x
	data["y"] = y
	data["name"] = s.Name
	var tmpRes bytes.Buffer
	err := splineTmpl.Execute(&tmpRes, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(tmpRes.String())
}

func (sa Spline) isEqual(sb Spline) bool {
	return slices.Equal(sa.X, sb.X) && slices.Equal(sa.Y, sb.Y)
}

func (sa Spline) equal(sb Spline) bool {
	return slices.Equal(sa.X, sb.X) && slices.Equal(sa.Y, sb.Y)
}
