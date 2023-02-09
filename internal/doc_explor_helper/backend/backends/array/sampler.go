package array

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
)

var _ feature.Sampler[Info] = Backend{} // Verify Merger is implemented.

func (b Backend) CompleteWithSample(info *Info) {
	b.complete(info, true)
}

func (b Backend) UseSample(info Info) bool {
	return info.topPercent.Total() > 500
}
