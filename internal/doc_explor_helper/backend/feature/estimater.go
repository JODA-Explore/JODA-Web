package feature

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
)

// Estimater is a feature to first estimate rating score with less time.
type Estimater[I result.Info] interface {
	Main[I]
	Estimate(*I) rating.Rater
	Complete(info *I)

	// EstNum return the limit number of info-candidates.
	EstNum(filter.Filter) int
}

type EstimaterInfo interface {
	EstimatedRating() rating.Rater
}

func estimaterTopK[I result.Info](e Estimater[I], infos []I, f filter.Filter) topK[I] {
	estimated := newTopK[I](e.EstNum(f))
	for i := range infos {
		info := infos[i]
		estimated.Insert(info, e.Estimate(&info).Score())
	}
	return estimated
}

// func infosToRs[I result.Info](infos []*I, f filter.Filter) (rs result.Results) {
// 	for _, x := range infos {
// 		rs.TryInsert(*x, (*x).Rating(), f)
// 	}
// 	return
// }

// func infosToTopK[I result.Info](infos []*I, n int) (t topK[I]) {
// 	t = newTopK[I](n)
// 	for _, x := range infos {
// 		t.insert(*x)
// 	}
// 	return t
// }
