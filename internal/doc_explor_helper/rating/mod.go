package rating

type Mod struct {
	Fun  func(Spline) Spline
	Desc string
}

type Mods []Mod

func reverse(s Spline) Spline {
	length := len(s.Y)
	y := make([]float64, length)
	for i, x := range s.Y {
		y[length-1-i] = x
	}
	s.Y = y
	return s
}

var ReverseMod = Mod{
	Fun:  reverse,
	Desc: "Reversed",
}

func (mods Mods) Apply(s Spline) Spline {
	newS := s
	for _, x := range mods {
		newS = x.Fun(newS)
	}
	newS.Name = mods.Desc() + s.Name
	return newS
}

func (mods Mods) Desc() (res string) {
	for _, m := range mods {
		res += m.Desc + " "
	}
	return
}
