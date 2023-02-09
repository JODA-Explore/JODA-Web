package feature

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

type Sampler[I result.Info] interface {
	UseSample(I) bool
	CompleteWithSample(info *I)
	SampleDesc(I, int) (template.HTML, error)
}
