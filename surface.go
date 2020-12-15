package ray

type Surface struct {
	Ka     float64
	Kd     float64
	Ks     float64
	Color  FloatColor
	Nphong float64
	// textures ...
}

var DefaultSurface = Surface{
	Ka:     0.5,
	Kd:     0.5,
	Ks:     0.5,
	Color:  FloatColor{R: 0.5, G: 0.5, B: 0.5},
	Nphong: 30,
}

var Diffuse = Surface{
	Ka:     0.5,
	Kd:     0.8,
	Ks:     0.2,
	Color:  FloatColor{R: 0.6, G: 0.4, B: 0.3},
	Nphong: 30,
}

func (s *Surface) ColorAt(h *Hit) FloatColor {
	//log.Printf("s %v", s)
	return s.Color
}

var Ocher2 = Surface{
	Ka:     0.7,
	Kd:     0.5,
	Ks:     0.4,
	Color:  FloatColor{R: 1, G: 1 / 1.5, B: 1 / 2},
	Nphong: 50,
}

var Mirror = Surface{
	Ka:     0,
	Kd:     0.1,
	Ks:     0.9,
	Color:  FloatColor{R: 0, G: 0, B: 0},
	Nphong: 1000,
}

var White1 = Surface{
	Ka:     0.5,
	Kd:     0.9,
	Ks:     0.5,
	Color:  FloatColor{R: 1, G: 1, B: 1},
	Nphong: 100,
}

var Building = Surface{
	Ka:     0.5,
	Kd:     0.9,
	Ks:     0.2,
	Color:  FloatColor{R: .5, G: .5, B: .5},
	Nphong: 10,
}
