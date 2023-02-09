package rating

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type hybridConds struct {
	hasExists, hasNonExists bool
}

func (hc hybridConds) ToScale() float64 {
	if hc.hasExists && hc.hasNonExists {
		return 25
	}
	return 0
}

func (hc hybridConds) Desc() template.HTML {
	var block template.HTML
	if hc.hasExists && hc.hasNonExists {
		block = "Bonused because this query has both exists and non-exists conditions"
	} else {
		if hc.hasExists {
			block = "Zero because this query has only exists conditions"
		} else {
			block = "Zero because this query has only non-exists conditions"
		}
	}
	return desc.Collapse(
		desc.Scale(hc.Name(), hc.ToScale()),
		"",
		block,
		"",
	)
}

func (hc hybridConds) Name() string {
	return "Hybrid Conditions Bonus"
}

func (hc *hybridConds) merge(new hybridConds) {
	hc.hasExists = hc.hasExists || new.hasExists
	hc.hasNonExists = hc.hasNonExists || new.hasNonExists
}
