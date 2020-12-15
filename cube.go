package ray

import "math"

// Cube is a canonical cube, centered of origin, of side 2 (-1 to +1)
type Cube struct {
	Transform
	Surface
	name string
}

// NewCube instantiate a new cube
func NewCube() *Cube {
	return &Cube{
		Transform: IDTransform,
		Surface:   DefaultSurface,
	}
}

// SetName ...
func (c *Cube) SetName(name string) {
	c.name = "cube:" + name
}

// Name returns the Cube's name
func (c *Cube) Name() string {
	return c.name
}

func (c *Cube) Surf() *Surface {
	return &c.Surface
}

// Translate applies a translation to the cube
func (c *Cube) Translate(x, y, z float64) *Cube {
	c.Transform.Translate(x, y, z)
	return c
}

// RotateX applies a rotation around x-axis to the Cube
func (c *Cube) RotateX(x float64) *Cube {
	c.Transform.RotateX(x)
	return c
}

// RotateY applies a rotation around y-axis to the Cube
func (c *Cube) RotateY(y float64) *Cube {
	c.Transform.RotateY(y)
	return c
}

// RotateZ applies a rotation around z-axis to the Cube
func (c *Cube) RotateZ(z float64) *Cube {
	c.Transform.RotateZ(z)
	return c
}

// Scale applies a scaling transform to the Cube
func (c *Cube) Scale(x, y, z float64) *Cube {
	c.Transform.Scale(x, y, z)
	return c
}

// Intersect ...
func (c *Cube) Intersect(r Ray) *Hit {
	var x, y, z, t float64
	locRay := c.RayToLocal(r)
	localDir := locRay.dir
	localPoint := locRay.pt

	minT := math.MaxFloat64
	var h Hit
	h.globRay = r
	h.locRay = locRay

	// inters y = -1 et y = 1
	if !isNul(localDir[Y]) {
		// y = -1
		t = (-1 - localPoint[Y]) / localDir[Y]
		if t > Epsilon {
			x = localPoint[X] + t*localDir[X]
			if (x > -1-Epsilon) && (x < 1-Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > -1-Epsilon) && (z < 1-Epsilon) {
					minT = t
					h.locNorm.pt = Point3{x, -1, z}
					h.locNorm.dir = Vector3{0, -1, 0}
				}
			}
		}
		// y =1
		t = (1 - localPoint[Y]) / localDir[Y]
		if (t > Epsilon) && (t < minT) {
			x = localPoint[X] + t*localDir[X]
			if (x > -1+Epsilon) && (x < 1+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > -1+Epsilon) && (z < 1+Epsilon) {
					minT = t
					h.locNorm.pt = Point3{x, 1, z}
					h.locNorm.dir = Vector3{0, 1, 0}
				}
			}
		}
	}
	// inters x = -1 et x = 1
	if !isNul(localDir[X]) {
		// x = -1
		t = (-1 - localPoint[X]) / localDir[X]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > -1+Epsilon) && (y < 1+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > -1+Epsilon) && (z < 1+Epsilon) {
					minT = t
					h.locNorm.pt = Point3{-1, y, z}
					h.locNorm.dir = Vector3{-1, 0, 0}
				}
			}
		}
		// x = 1
		t = (1 - localPoint[X]) / localDir[X]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > -1-Epsilon) && (y < 1-Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > -1-Epsilon) && (z < 1-Epsilon) {
					minT = t
					h.locNorm.pt = Point3{1, y, z}
					h.locNorm.dir = Vector3{1, 0, 0}
				}
			}
		}
	}

	// inters z = -1 et z = 1
	if !isNul(localDir[Z]) {
		// z = -1
		t = (-1 - localPoint[Z]) / localDir[Z]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > -1+Epsilon) && (y < 1+Epsilon) {
				x = localPoint[X] + t*localDir[X]
				if (x > -1-Epsilon) && (x < 1-Epsilon) {
					minT = t
					h.locNorm.pt = Point3{x, y, -1}
					h.locNorm.dir = Vector3{0, 0, -1}
				}
			}
		}
		// z = 1
		t = (1 - localPoint[Z]) / localDir[Z]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > -1-Epsilon) && (y < 1-Epsilon) {
				x = localPoint[X] + t*localDir[X]
				if (x > -1+Epsilon) && (x < 1+Epsilon) {
					minT = t
					h.locNorm.pt = Point3{x, y, 1}
					h.locNorm.dir = Vector3{0, 0, 1}
				}
			}
		}
	}

	if minT < math.MaxFloat64-Epsilon {
		h.globNorm = c.RayToGlobal(h.locNorm)
		h.globNorm.Normalize()
		h.Surface = &c.Surface
		return &h
	}
	return nil
}

func (c *Cube) MinMax() (Point3, Point3) {
	return Point3{-1, -1, -1}, Point3{1, 1, 1}
}
