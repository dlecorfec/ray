package ray

type Light interface {
	RayToLight(Point3) Ray
	Color(Ray) FloatColor
	Sun() bool
}

type PointLight struct {
	Transform
	c   FloatColor
	sun bool
}

func NewPointLight(c FloatColor) *PointLight {
	return &PointLight{Transform: IDTransform, c: c}
}

func (p *PointLight) RayToLight(pt Point3) Ray {
	// XXX cache PointToGlobal result
	return NewRay(pt, p.PointToGlobal(Origin))
}

func (p *PointLight) Color(r Ray) FloatColor {
	return p.c
}

// Translate applies a translation
func (p *PointLight) Translate(x, y, z float64) *PointLight {
	p.Transform.Translate(x, y, z)
	return p
}

func (p *PointLight) Sun() bool {
	return p.sun
}

func (p *PointLight) SetSun(sun bool) {
	p.sun = sun
}
