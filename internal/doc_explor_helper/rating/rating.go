package rating

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type Rating struct {
	X      float64
	spline Spline
	Factors
	score float64
	desc  template.HTML
}

func (r *Rating) Score() float64 {
	if r.score == 0 {
		r.UpdateScore()
	}
	return r.score
}

func (r *Rating) Update() {
	r.UpdateScore()
	r.UpdateDesc()
}

func (r *Rating) UpdateScore() float64 {
	s := r.spline.Cal(r.X, r.Factors.sumScale())
	r.score = normalize(s)
	return r.score
}

var ratingTmpl = template.Must(template.New("rating").Parse(`
<div>
    {{.Factors}}
    <hr>
    <div>
        Original:<br> {{ .splineOri}}
    </div>
    <p> Scaled: </p>
    {{.splineScaled}}
    x: <b>{{.x}}</b> <span style="font-size:140%;">&nbsp; <b>&#8658;</b>&nbsp;  </span> Result: <b>{{.score}}</b>
</div>
`))

func (r *Rating) Desc() template.HTML {
	if r.desc == "" {
		r.UpdateDesc()
	}
	return r.desc
}

func (r *Rating) UpdateDesc() {
	splineOri := r.spline
	scaleSum := r.sumScale()
	splineScaled := splineOri.ApplyScale(scaleSum)
	data := map[string]interface{}{
		"x":            desc.Float(r.X, 6),
		"splineOri":    splineOri.Desc(),
		"Factors":      r.Factors.Desc(),
		"scale":        desc.Scale("Scale", scaleSum),
		"splineScaled": splineScaled.Desc(),
		"score":        desc.Float(r.score, 4),
	}
	r.desc = desc.ExecuteTemplateWeb(*ratingTmpl, data)
}
