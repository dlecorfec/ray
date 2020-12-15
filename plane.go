package ray

import "log"

// Object is the interface for all raytraced objects
type Object interface {
	Transformer
	Intersect(Ray) *Hit
	Name() string
	MinMax() (Point3, Point3) // in local coords
}

// Plane is a plane of all points with y = 0
type Plane struct {
	Transform
	Surface
	name string
}

// NewPlane creates a xz plane of 2 units sides centered on origin
func NewPlane() *Plane {
	return &Plane{
		Transform: IDTransform,
		Surface:   DefaultSurface,
	}
}

// SetName ...
func (p *Plane) SetName(name string) {
	p.name = "plane:" + name
}

// Name returns the object name
func (p *Plane) Name() string {
	return p.name
}

// Surf ...
func (p *Plane) Surf() *Surface {
	return &p.Surface
}

// Translate applies a translation
func (p *Plane) Translate(x, y, z float64) *Plane {
	p.Transform.Translate(x, y, z)
	return p
}

// RotateX applies a rotation around x-axis
func (p *Plane) RotateX(x float64) *Plane {
	p.Transform.RotateX(x)
	return p
}

// RotateY applies a rotation around y-axis
func (p *Plane) RotateY(y float64) *Plane {
	p.Transform.RotateY(y)
	return p
}

// RotateZ applies a rotation around z-axis
func (p *Plane) RotateZ(z float64) *Plane {
	p.Transform.RotateZ(z)
	return p
}

// Scale applies a scaling transform
func (p *Plane) Scale(x, y, z float64) *Plane {
	p.Transform.Scale(x, y, z)
	return p
}

// Intersect returns intersected object and normal ray
func (p *Plane) Intersect(r Ray) *Hit {
	locRay := p.RayToLocal(r)
	if isNul(locRay.dir[Y]) {
		return nil
	}
	t := -locRay.pt[Y] / locRay.dir[Y]
	if t < Epsilon {
		return nil
	}
	x := t*locRay.dir[X] + locRay.pt[X]
	if x < -1-Epsilon || x > 1+Epsilon {
		return nil
	}
	z := t*locRay.dir[Z] + locRay.pt[Z]
	if z < -1-Epsilon || z > 1+Epsilon {
		return nil
	}
	if p.debug(r) {
		log.Printf("%s inter: locpt=%v t=%f", p.Name(), locRay.pt, t)
	}
	var h Hit
	h.locRay = locRay
	h.locNorm.pt[X] = x
	h.locNorm.pt[Z] = z
	h.locNorm.dir[Y] = 1
	if locRay.pt[Y] < 0 {
		h.locNorm.dir[Y] = -1
	}
	h.globNorm = p.RayToGlobal(h.locNorm)
	h.globNorm.Normalize()
	h.Surface = &p.Surface
	h.globRay = r
	if p.debug(r) {
		log.Printf("%s inter norm=%v x,z=%.3f,%.3f ", p.Name(), h.globNorm.dir, x, z)
	}
	//log.Printf("%#v %#v\n", p, &sd)
	return &h
}

// MinMax returns the min max points of the bounding box in local coords
func (*Plane) MinMax() (Point3, Point3) {
	return Point3{-1, 0, -1}, Point3{1, 0, 1}
}

func (*Plane) debug(r Ray) bool {
	if r.x == 360 && (r.y == 200 || r.y == 250) {
		return true
	}
	return false
}
