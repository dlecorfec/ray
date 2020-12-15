package ray

import (
	"math"
)

// Sphere is a canonical sphere, centered of origin, of radius 1
type Sphere struct {
	Transform
	Surface
	name string
}

// NewSphere instantiate a new sphere
func NewSphere() *Sphere {
	return &Sphere{
		Transform: IDTransform,
		Surface:   DefaultSurface,
	}
}

// SetName ...
func (s *Sphere) SetName(name string) {
	s.name = "sphere:" + name
}

// Name returns the sphere's name
func (s *Sphere) Name() string {
	return s.name
}

// Surf ...
func (s *Sphere) Surf() *Surface {
	return &s.Surface
}

// Translate applies a translation to the sphere
func (s *Sphere) Translate(x, y, z float64) *Sphere {
	s.Transform.Translate(x, y, z)
	return s
}

// RotateX applies a rotation around x-axis to the sphere
func (s *Sphere) RotateX(x float64) *Sphere {
	s.Transform.RotateX(x)
	return s
}

// RotateY applies a rotation around y-axis to the sphere
func (s *Sphere) RotateY(y float64) *Sphere {
	s.Transform.RotateY(y)
	return s
}

// RotateZ applies a rotation around z-axis to the sphere
func (s *Sphere) RotateZ(z float64) *Sphere {
	s.Transform.RotateZ(z)
	return s
}

// Scale applies a scaling transform to the sphere
func (s *Sphere) Scale(x, y, z float64) *Sphere {
	s.Transform.Scale(x, y, z)
	return s
}

// Intersect finds the intersection point (if any) between a global ray
// and this sphere, and the normal at intersection point
func (s *Sphere) Intersect(r Ray) *Hit {
	locRay := s.RayToLocal(r)
	dir := locRay.dir
	p := locRay.pt
	a := square(dir[X]) + square(dir[Y]) + square(dir[Z])
	b := 2 * (dir[X]*p[X] + dir[Y]*p[Y] + dir[Z]*p[Z])
	c := square(p[X]) + square(p[Y]) + square(p[Z]) - 1
	delta := square(b) - 4*a*c
	if delta < BigEpsilon {
		return nil
	}
	sqd := math.Sqrt(delta)
	aa := 2 * a
	t := (sqd - b) / aa
	t2 := (-b - sqd) / aa
	if t > Epsilon && t2 > Epsilon {
		t = math.Min(t, t2)
	} else if t2 > Epsilon {
		t = t2
	} else if t < Epsilon {
		return nil
	}

	var h Hit
	//log.Printf("gr=%v lr=%v", r, locRay)
	h.globRay = r
	h.locRay = locRay

	lp := Point3{p[X] + t*dir[X], p[Y] + t*dir[Y], p[Z] + t*dir[Z]}
	h.locNorm.pt = lp
	h.locNorm.dir = Vector3(lp)
	gp := s.PointToGlobal(lp)
	h.globNorm.pt = gp
	d := -(square(lp[X]) + square(lp[Y]) + square(lp[Z]))
	var A, B Point3
	//log.Printf("INTERSECT lp=%v gp=%v d=%v t=%v ray=%v locray=%v delta=%f", lp, gp, d, t, r, locRay, delta)

	absx := math.Abs(lp[X])
	absy := math.Abs(lp[Y])
	absz := math.Abs(lp[Z])

	switch {
	case absx >= absy && absx >= absz:
		A = Point3{-(lp[Z] + d) / lp[X], 0, 1}
		B = Point3{-(lp[Y] + d) / lp[X], 1, 0}
	case absy >= absx && absy >= absz:
		A = Point3{0, -(lp[Z] + d) / lp[Y], 1}
		B = Point3{1, -(lp[X] + d) / lp[Y], 0}
	default:
		A = Point3{0, 1, -(lp[Y] + d) / lp[Z]}
		B = Point3{1, 0, -(lp[X] + d) / lp[Z]}
	}
	//log.Printf("local A=%v B=%v", A, B)
	A = s.PointToGlobal(A)
	B = s.PointToGlobal(B)
	//log.Printf("gA=%v gB=%v", A, B)
	u := NewVec(gp, A)
	v := NewVec(gp, B)

	h.globNorm.dir = u.Cross(v)
	//log.Printf("globNorm: %v", h.globNorm.dir)
	h.globNorm.Normalize()
	o := s.PointToGlobal(Origin)
	w := Vector3{gp[X] - o[X], gp[Y] - o[Y], gp[Z] - o[Z]}
	//log.Printf("SPHERE u=%v v=%v gn=%v w=%v dot(gn,gp)=%f", u, v, h.globNorm.dir, w, h.globNorm.dir.Dot(w))
	if h.globNorm.dir.Dot(w) < 0 {
		h.globNorm.dir.Reverse()
	}
	h.Surface = &s.Surface
	return &h
}

/*
func debug(r Ray) bool {
	if r.x == 300 && r.y == 220 {
		return true
	}
	return false
}
*/

// MinMax ...
func (s *Sphere) MinMax() (Point3, Point3) {
	return Point3{-1, -1, -1}, Point3{1, 1, 1}
}
