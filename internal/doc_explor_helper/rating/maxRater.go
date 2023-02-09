package rating

import (
	"html/template"
	"strconv"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

type MaxRater struct {
	score  float64
	desc   template.HTML
	raters []Rater
	Factors
}

func NewMaxRater(raters []Rater, factors ...Factor) MaxRater {
	return MaxRater{raters: raters, Factors: factors}
}

func (mr *MaxRater) Push(r Rater) {
	mr.raters = append(mr.raters, r)
	mr.Update()
}

func (mr *MaxRater) PushFront(r Rater) {
	mr.raters = slices.PushFront(mr.raters, r)
	mr.Update()
}

func (mr *MaxRater) Score() float64 {
	if mr.score == 0 {
		mr.updateScore()
	}
	return mr.score
}

func (mr *MaxRater) maxScore() float64 {
	_, score := slices.Max(mr.raters, func(r Rater) float64 { return r.Score() })
	return score
}

func (mr *MaxRater) updateScore() {
	if len(mr.raters) == 0 {
		mr.score = 0
	} else {
		score := mr.maxScore()
		score = applyScale(score, mr.sumScale())
		mr.score = normalize(score)
	}
}

func (mr *MaxRater) UpdateDesc() {
	if len(mr.raters) == 1 {
		mr.desc = mr.raters[0].Desc()
		return
	}
	var block template.HTML
	for i, r := range mr.raters {
		block += desc.Collapse(
			desc.InlineRating("Score"+strconv.Itoa(i), r.Score()),
			"",
			r.Desc(),
			"",
		)
	}
	mr.desc = mr.Factors.Desc() +
		desc.Collapse(
			desc.InlineRating("Max Score", mr.maxScore()),
			"<br>It was reduced by the following results, the max score will be picked.",
			block, "")
}

func (mr *MaxRater) Desc() template.HTML {
	if mr.desc == "" {
		mr.Update() // score also need to be updated
	}
	return mr.desc
}

func (mr *MaxRater) Update() {
	mr.updateScore()
	mr.UpdateDesc()
}
