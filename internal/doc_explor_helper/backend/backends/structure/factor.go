package structure

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/cond"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
)

type hybridCondsBonus struct {
	exists, nonExists bool
}

func newHybridCondsBonus(conds cond.Conditions) hybridCondsBonus {
	exists, nonExists := conds.GroupByNeg()
	return hybridCondsBonus{exists: len(exists) > 0, nonExists: len(nonExists) > 0}
}

func (hcb hybridCondsBonus) hasBonus() bool {
	return hcb.exists && hcb.nonExists
}

func (hcb hybridCondsBonus) ToScale() float64 {
	if hcb.hasBonus() {
		return 12
	}
	return 0
}

func (hcb hybridCondsBonus) Desc() template.HTML {
	var block template.HTML
	if hcb.hasBonus() {
		block = "Bonused because this query has both exists and non-exists conditions"
	} else {
		if hcb.exists {
			block = "Zero because this query has only exists conditions"
		} else {
			block = "Zero because this query has only non-exists conditions"
		}
	}
	return desc.Collapse(
		desc.Scale(hcb.Name(), hcb.ToScale()),
		"",
		block,
		"",
	)
}

func (hcb hybridCondsBonus) Name() string {
	return "Hybrid Conditions Bonus"
}
