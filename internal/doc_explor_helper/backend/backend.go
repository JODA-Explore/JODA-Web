package backend

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

type Backend interface {
	Results(pts *points.Trie, f filter.Filter) (result.Results, error)
	Desc(r *result.Result, idx int) error
	QueryMaker(r *result.Result, idx int) (template.HTML, error)
}
