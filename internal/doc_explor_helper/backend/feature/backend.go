package feature

import (
	"html/template"

	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/backend"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/filter"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/points"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/result"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/desc"
	"github.com/JODA-Explore/JODA-Web/internal/doc_explor_helper/utils/trie"
)

// Backend is a wraper of the Main interface.
type Backend[I result.Info] struct {
	Main[I]
}

func NewBackend[I result.Info](m Main[I]) backend.Backend {
	return Backend[I]{m}
}

func (b Backend[I]) Results(pts *points.Trie, f filter.Filter) (rs result.Results, err error) {
	var topK topK[I]
	m, multiOK := b.Main.(Multier[I])
	e, estOK := b.Main.(Estimater[I])
	s, smpOK := b.Main.(Sampler[I])
	if multiOK {
		topK = newTopK[I](m.MultiNum(f))
	}
	complete := func(info *I) {
		if smpOK && s.UseSample(*info) {
			s.CompleteWithSample(info)
		} else {
			e.Complete(info)
		}
	}
	filter := b.Filter()
	pts.Walk(func(path string, v *points.Value, depth int) (ctl trie.Control) {
		if v.Info.CountNull == v.Info.CountTotal {
			return trie.Continue
		}
		ctl = filter(path, v, depth)
		if ctl != trie.None {
			return
		}
		infos, ctl := b.Infos(path, v.Info)
		if smpOK && !estOK {
			for i := range infos {
				info := infos[i]
				if s.UseSample(info) {
					s.CompleteWithSample(&info)
					infos[i] = info
				}
			}
		}
		switch {
		case multiOK && estOK:
			estimated := estimaterTopK(e, infos, f)
			for _, x := range estimated.List() {
				info := x.Val
				complete(&info)
				topK.insert(info)
			}
		case multiOK:
			for _, info := range infos {
				topK.insert(info)
			}
		case estOK:
			estimated := estimaterTopK(e, infos, f)
			for _, x := range estimated.List() {
				info := x.Val
				complete(&info)
				rs.TryInsert(info, info.Rating(), f)
			}
		default:
			for _, info := range infos {
				rs.TryInsert(info, info.Rating(), f)
			}
		}
		return ctl
	})
	if multiOK {
		rs = multierRs(m, topK, f)
	}
	return
}

func (b Backend[I]) QueryMaker(r *result.Result, idx int) (template.HTML, error) {
	if r.QryMaker == "" {
		var err error
		if mixed, ok := r.Info.(*result.MixedInfo[I]); ok {
			multier := b.Main.(Multier[I])
			r.QryMaker, err = multier.MultiQryMaker(mixed.Infos, mixed.Extra(), idx)
			if err != nil {
				return "", err
			}
		} else {
			r.QryMaker, err = b.QryMaker(r.Info.(I), idx)
			if err != nil {
				return "", err
			}
		}
	}
	return r.QryMaker, nil
}

func (b Backend[I]) Desc(r *result.Result, idx int) (err error) {
	if mixed, ok := r.Info.(*result.MixedInfo[I]); ok {
		multier := b.Main.(Multier[I])
		r.Desc, err = multier.MultiDesc(mixed.Infos, mixed.Extra(), idx)
		if err != nil {
			return
		}
	} else {
		info := r.Info.(I)
		r.Desc, err = b.Main.Desc(info, idx)
		if err != nil {
			return
		}
		if s, ok := b.Main.(Sampler[I]); ok {
			if s.UseSample(info) {
				sampleDesc, err := s.SampleDesc(info, idx)
				if err != nil {
					return err
				}
				if sampleDesc != "" {
					r.Desc += sampleDesc
				}
			}
		}
	}
	if _, ok := b.Main.(Estimater[I]); ok {
		r.Desc += desc.LazyCollapse(
			"<b>Guess Rating</b>",
			"guess_rating",
			"track_guess_rating",
			"",
			idx,
		)
	}
	return
}
