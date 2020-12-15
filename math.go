package ray

import (
	"math"
)

//var Epsilon = math.Nextafter(1.0, 2.0) - 1.0
var Epsilon = 1e-8
var BigEpsilon = 1e-3

type Matrix4 [4][4]float64
type Vector3 [3]float64
type Point3 [3]float64

var Origin = Point3{0, 0, 0}

const (
	X int = 0
	Y     = 1
	Z     = 2
)

// isNul returns true if the given number is close enough
// (depending on Epsilon) to 0.
func isNul(x float64) bool {
	return x > -Epsilon && x < Epsilon
}

// square returns the square of a number.
func square(x float64) float64 {
	return x * x
}

// Dist returns the distance between 2 points a and b.
func (a Point3) Dist(b Point3) float64 {
	return math.Sqrt(square(b[0]-a[0]) + square(b[1]-a[1]) + square(b[2]-a[2]))
}

// SquareDist returns the squared distance between points a and b.
func (a Point3) SquareDist(b Point3) float64 {
	return square(b[0]-a[0]) + square(b[1]-a[1]) + square(b[2]-a[2])
}

// NewVec creates a new vector from point a to point b.
func NewVec(a, b Point3) Vector3 {
	return Vector3{b[X] - a[X], b[Y] - a[Y], b[Z] - a[Z]}
}

// Norm returns the norm of the vector.
func (v Vector3) Norm() float64 {
	return math.Sqrt(square(v[X]) + square(v[Y]) + square(v[Z]))
}

// Normalize the vector, which is modified (has a norm of 1).
func (v *Vector3) Normalize() {
	n := v.Norm()
	*v = v.Mult(1 / n)
}

// Mult returns a new vector where all components have been
// multiplied by the given factor.
func (v Vector3) Mult(f float64) Vector3 {
	return Vector3{v[X] * f, v[Y] * f, v[Z] * f}
}

// Sub substracts the second vector from the first and returns
// a new vector.
func (v Vector3) Sub(b Vector3) Vector3 {
	return Vector3{v[X] - b[X], v[Y] - b[Y], v[Z] - b[Z]}
}

// Add adds 2 vectors and returns a new vector.
func (v Vector3) Add(b Vector3) Vector3 {
	return Vector3{v[X] + b[X], v[Y] + b[Y], v[Z] + b[Z]}
}

// Reverse modifies the vector to reverse its direction.
func (v *Vector3) Reverse() {
	v[X] = -v[X]
	v[Y] = -v[Y]
	v[Z] = -v[Z]
}

// Dot calculates the dot product of 2 vectors (equals to |u|*|v|*cos(uv))
func (a Vector3) Dot(b Vector3) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

// Cross returns a new vector which is the cross product of the
// vector parameters.
// The result is a vector orthogonal to the plane formed by
// the 2 input vectors, and such that uvw is a direct ...
func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{
		a[Y]*b[Z] - a[Z]*b[Y],
		a[Z]*b[X] - a[X]*b[Z],
		a[X]*b[Y] - a[Y]*b[X],
	}
}

// Ray contains a point (the origin) and a vector (the direction).
type Ray struct {
	pt  Point3
	dir Vector3
	x   int
	y   int
}

// NewRay creates a Ray going from the starting point to the
// destination point.
func NewRay(src, dst Point3) Ray {
	return Ray{
		pt:  src,
		dir: Vector3{dst[X] - src[X], dst[Y] - src[Y], dst[Z] - src[Z]},
	}
}

// Normalize modifies the ray by normalizing its direction.
func (r *Ray) Normalize() {
	norm := math.Sqrt(square(r.dir[X]) + square(r.dir[Y]) + square(r.dir[Z]))
	if norm <= Epsilon {
		return
	}
	r.dir[X] /= norm
	r.dir[Y] /= norm
	r.dir[Z] /= norm
}

// ID returns a homogeneous, 4x4, row-major, identity matrix.
func ID() Matrix4 {
	return Matrix4{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}}
}

// Translation returns a matrix describing a translation along all 3 axis in space.
func Translation(x, y, z float64) Matrix4 {
	return Matrix4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

// Scaling returns a matrix describing a scaling along all 3 axis in space.
func Scaling(x, y, z float64) Matrix4 {
	return Matrix4{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

// RotationX returns a matrix describing a rotation in radians around X-axis
// in right-hand coordinates.
func RotationX(a float64) Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, math.Cos(a), -math.Sin(a), 0},
		{0, math.Sin(a), math.Cos(a), 0},
		{0, 0, 0, 1},
	}
}

// RotationY returns a matrix describing a rotation in radians around Y-axis
// in right-hand coordinates.
func RotationY(a float64) Matrix4 {
	return Matrix4{
		{math.Cos(a), 0, math.Sin(a), 0},
		{0, 1, 0, 0},
		{-math.Sin(a), 0, math.Cos(a), 0},
		{0, 0, 0, 1},
	}
}

// RotationZ returns a matrix describing a rotation in radians around Z-axis
// in right-hand coordinates.
func RotationZ(a float64) Matrix4 {
	return Matrix4{
		{math.Cos(a), -math.Sin(a), 0, 0},
		{math.Sin(a), math.Cos(a), 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// MulM multiply a with b, return resulting matrix
func (a Matrix4) MulM(b Matrix4) Matrix4 {
	return Matrix4{
		{
			a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0] + a[0][3]*b[3][0],
			a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1] + a[0][3]*b[3][1],
			a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2] + a[0][3]*b[3][2],
			a[0][0]*b[0][3] + a[0][1]*b[1][3] + a[0][2]*b[2][3] + a[0][3]*b[3][3],
		},
		{
			a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0] + a[1][3]*b[3][0],
			a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1] + a[1][3]*b[3][1],
			a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2] + a[1][3]*b[3][2],
			a[1][0]*b[0][3] + a[1][1]*b[1][3] + a[1][2]*b[2][3] + a[1][3]*b[3][3],
		},
		{
			a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0] + a[2][3]*b[3][0],
			a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1] + a[2][3]*b[3][1],
			a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2] + a[2][3]*b[3][2],
			a[2][0]*b[0][3] + a[2][1]*b[1][3] + a[2][2]*b[2][3] + a[2][3]*b[3][3],
		},
		{
			a[3][0]*b[0][0] + a[3][1]*b[1][0] + a[3][2]*b[2][0] + a[3][3]*b[3][0],
			a[3][0]*b[0][1] + a[3][1]*b[1][1] + a[3][2]*b[2][1] + a[3][3]*b[3][1],
			a[3][0]*b[0][2] + a[3][1]*b[1][2] + a[3][2]*b[2][2] + a[3][3]*b[3][2],
			a[3][0]*b[0][3] + a[3][1]*b[1][3] + a[3][2]*b[2][3] + a[3][3]*b[3][3],
		},
	}
}

// MulP returns a point resulting from the application of the
// transformation matrix to the given point.
func (a Matrix4) MulP(b Point3) Point3 {
	return Point3{
		a[0][0]*b[0] + a[0][1]*b[1] + a[0][2]*b[2] + a[0][3],
		a[1][0]*b[0] + a[1][1]*b[1] + a[1][2]*b[2] + a[1][3],
		a[2][0]*b[0] + a[2][1]*b[1] + a[2][2]*b[2] + a[2][3],
	}
}

// MulV returns a vector resulting from the application of the
// transformation matrix to the given vector.
func (a Matrix4) MulV(b Vector3) Vector3 {
	return Vector3{
		a[0][0]*b[0] + a[0][1]*b[1] + a[0][2]*b[2],
		a[1][0]*b[0] + a[1][1]*b[1] + a[1][2]*b[2],
		a[2][0]*b[0] + a[2][1]*b[1] + a[2][2]*b[2],
	}
}

// Transpose returns a new matrix which is the transposition of the given matrix.
func Transpose(m Matrix4) Matrix4 {
	return Matrix4{
		{m[0][0], m[1][0], m[2][0], m[3][0]},
		{m[0][1], m[1][1], m[2][1], m[3][1]},
		{m[0][2], m[1][2], m[2][2], m[3][2]},
		{m[0][3], m[1][3], m[2][3], m[3][3]},
	}
}
