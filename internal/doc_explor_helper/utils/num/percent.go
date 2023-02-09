package num

type Percent struct {
	count, total uint64
	ratio        float64
}

func (p Percent) Count() uint64 {
	return p.count
}

func (p *Percent) SetCount(count uint64) {
	p.count = count
}

func (p Percent) Total() uint64 {
	return p.total
}

func (p *Percent) SetTotal(total uint64) {
	p.total = total
}

func (p Percent) Empty() bool {
	return p.total == 0 && p.count == 0
}

func (p *Percent) Ratio() float64 {
	if p.ratio == 0 {
		if p.total != 0 {
			p.ratio = percent(p.count, p.total)
		}
	}
	return p.ratio
}

func NewPercent(count, total uint64) Percent {
	return Percent{count, total, 0}
}

func (old *Percent) EstimatePercent(newTotal uint64) Percent {
	return NewPercent(uint64(float64(newTotal)*old.Ratio()), newTotal)
}
