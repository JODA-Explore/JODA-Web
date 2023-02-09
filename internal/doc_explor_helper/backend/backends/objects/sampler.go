package objects

import (
	"fmt"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/feature"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/num"
)

var _ feature.Sampler[Info] = Backend{} // Verify Multier is implemented.

func (b Backend) UseSample(info Info) bool {
	return info.real.Total() > 500
}

func (b Backend) CompleteWithSample(info *Info) {
	sampleCount, sampleTotal, err := b.DistinctObjCounts(b.Sample(), info.JsonPoint())
	if err != nil {
		panic(fmt.Errorf("calculate sample DistinctObjCounts of %v:%w", info.JsonPoint(), err))
	}
	info.sample = num.NewPercent(sampleCount, sampleTotal)
	info.real = info.sample.EstimatePercent(info.real.Total())
	info.updateRating()
}
