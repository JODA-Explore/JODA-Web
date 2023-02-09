package result

import (
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend/backends"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/rating"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/slices"
)

type Info interface {
	Rating() rating.Rater
	Backend() backends.ID
}

type Infos []Info

func (is Infos) Results(f filter.Filter, rs *Results) {
	for _, info := range is {
		rs.TryInsert(info, info.Rating(), f)
	}
	return
}

type BasicInfo struct {
	jsonPoint string
	backendID backends.ID
	rating    rating.Rater
}

func NewBasicInfo(jsonPoint string, backendID backends.ID) BasicInfo {
	return BasicInfo{jsonPoint, backendID, nil}
}

func (bi BasicInfo) JsonPoint() string {
	return bi.jsonPoint
}

func (bi BasicInfo) Backend() backends.ID {
	return bi.backendID
}

func (bi BasicInfo) Rating() rating.Rater {
	return bi.rating
}

func (bi *BasicInfo) SetRating(r rating.Rater) {
	bi.rating = r
}

type MixedInfo[I Info] struct {
	rater   rating.MaxRater
	Infos   []I
	backend backends.ID
	extra   ExtraInfo[I]
}

type ExtraInfo[I Info] interface {
	Info() I
	Factors() rating.Factors
}

func NewMixedInfo[I Info](infos []I, extra ExtraInfo[I]) (res MixedInfo[I]) {
	res.extra = extra
	res.Infos = infos
	res.backend = infos[0].Backend()
	var factors rating.Factors
	if extra != nil {
		factors = extra.Factors()
	}
	raters := slices.Map(infos, func(info I) rating.Rater { return info.Rating() })
	res.rater = rating.NewMaxRater(raters, factors...)
	return
}

func (mi *MixedInfo[I]) Extra() ExtraInfo[I] {
	return mi.extra
}

func (mi *MixedInfo[I]) Rating() rating.Rater {
	return &mi.rater
}

func (mi *MixedInfo[I]) Backend() (ids backends.ID) {
	return mi.backend
}
