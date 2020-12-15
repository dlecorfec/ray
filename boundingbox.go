package ray

import (
	"math"
)

// BoundingBox ...
type BoundingBox struct {
	Transform
	childs []Object
	name   string
	min    Point3
	max    Point3
}

// NewBoundingBox ...
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		Transform: IDTransform,
	}
}

// SetName ...
func (bb *BoundingBox) SetName(name string) {
	bb.name = "bb:" + name
}

// Name returns the object's name
func (bb *BoundingBox) Name() string {
	return bb.name
}

// AddObjects ...
func (bb *BoundingBox) AddObjects(obj ...Object) {
	for _, o := range obj {
		bb.mergeBB(o)
		//log.Printf("min=%+v max=%+v", bb.min, bb.max)
		bb.childs = append(bb.childs, o)
	}
}

func (bb *BoundingBox) mergeBB(o Object) {
	lmin, lmax := o.MinMax()
	var x [8]Point3
	// take all the local obj_bb points
	x[0] = Point3{lmin[X], lmin[Y], lmin[Z]}
	x[1] = Point3{lmax[X], lmin[Y], lmin[Z]}
	x[2] = Point3{lmax[X], lmax[Y], lmin[Z]}
	x[3] = Point3{lmin[X], lmax[Y], lmin[Z]}
	x[4] = Point3{lmin[X], lmin[Y], lmax[Z]}
	x[5] = Point3{lmax[X], lmin[Y], lmax[Z]}
	x[6] = Point3{lmax[X], lmax[Y], lmax[Z]}
	x[7] = Point3{lmin[X], lmax[Y], lmax[Z]}
	// transform them to global and adjust bb minmax
	for i := 0; i < 8; i++ {
		x[i] = o.PointToGlobal(x[i])
		if x[i][X] < bb.min[X] {
			bb.min[X] = x[i][X]
		}
		if x[i][Y] < bb.min[Y] {
			bb.min[Y] = x[i][Y]
		}
		if x[i][Z] < bb.min[Z] {
			bb.min[Z] = x[i][Z]
		}
		if x[i][X] > bb.max[X] {
			bb.max[X] = x[i][X]
		}
		if x[i][Y] > bb.max[Y] {
			bb.max[Y] = x[i][Y]
		}
		if x[i][Z] > bb.max[Z] {
			bb.max[Z] = x[i][Z]
		}

	}
}

// Intersect ...
func (bb *BoundingBox) Intersect(r Ray) *Hit {
	locRay := bb.RayToLocal(r)
	if !bb.intersectBB(locRay) {
		return nil
	}
	var closest *Hit
	minD := math.MaxFloat64
	for _, o := range bb.childs {
		h := o.Intersect(locRay)
		if h == nil {
			continue
		}
		d := locRay.pt.Dist(h.globNorm.pt)
		if d < minD {
			closest = h
			minD = d
		}
	}
	if closest == nil {
		return nil
	}
	closest.globRay = bb.RayToGlobal(closest.globRay)
	closest.globNorm = bb.RayToGlobal(closest.globNorm)
	return closest
}

// MinMax ...
func (bb *BoundingBox) MinMax() (Point3, Point3) {
	return bb.min, bb.max
}

// intersectsBB ...
func (bb *BoundingBox) intersectBB(locRay Ray) bool {
	var x, y, z, t float64
	localDir := locRay.dir
	localPoint := locRay.pt

	minT := math.MaxFloat64

	// inters y = -1 et y = 1
	if !isNul(localDir[Y]) {
		// y = -1
		t = (bb.min[Y] - localPoint[Y]) / localDir[Y]
		if t > Epsilon {
			x = localPoint[X] + t*localDir[X]
			if (x > bb.min[X]-Epsilon) && (x < bb.max[X]+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > bb.min[Z]-Epsilon) && (z < bb.max[Z]+Epsilon) {
					return true
				}
			}
		}
		// y =1
		t = (bb.max[Y] - localPoint[Y]) / localDir[Y]
		if (t > Epsilon) && (t < minT) {
			x = localPoint[X] + t*localDir[X]
			if (x > bb.min[X]-Epsilon) && (x < bb.max[X]+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > bb.min[Y]-Epsilon) && (z < bb.max[Y]+Epsilon) {
					return true
				}
			}
		}
	}
	// inters x = -1 and x = 1
	if !isNul(localDir[X]) {
		// x min
		t = (bb.min[X] - localPoint[X]) / localDir[X]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > bb.min[Y]-Epsilon) && (y < bb.max[Y]+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > bb.min[Z]-Epsilon) && (z < bb.max[Z]+Epsilon) {
					return true
				}
			}
		}
		// x max
		t = (bb.max[X] - localPoint[X]) / localDir[X]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > bb.min[Y]-Epsilon) && (y < bb.max[Y]+Epsilon) {
				z = localPoint[Z] + t*localDir[Z]
				if (z > bb.min[Z]-Epsilon) && (z < bb.max[Z]+Epsilon) {
					return true
				}
			}
		}
	}

	// inters z = -1 et z = 1
	if !isNul(localDir[Z]) {
		// z = -1
		t = (bb.min[Z] - localPoint[Z]) / localDir[Z]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > bb.min[Y]-Epsilon) && (y < bb.max[Y]+Epsilon) {
				x = localPoint[X] + t*localDir[X]
				if (x > bb.min[X]-Epsilon) && (x < bb.max[X]+Epsilon) {
					return true
				}
			}
		}
		// z = 1
		t = (bb.max[Z] - localPoint[Z]) / localDir[Z]
		if (t > Epsilon) && (t < minT) {
			y = localPoint[Y] + t*localDir[Y]
			if (y > bb.min[Y]-Epsilon) && (y < bb.max[Y]+Epsilon) {
				x = localPoint[X] + t*localDir[X]
				if (x > bb.min[X]-Epsilon) && (x < bb.max[X]+Epsilon) {
					return true
				}
			}
		}
	}

	return false
}
