package ray

import (
	"image"
	"math"
)

// Camera is the canonical camera object, tracing rays and writing to an image
type Camera struct {
	Transform
	f      float64
	dx     float64
	dy     float64
	Width  int
	Height int
	Image  *image.RGBA
}

// NewCamera creates a camera with a focal length, camera width, camera height and image width.
// Located at origin and points towards negative z in an orthonormal basis
func NewCamera(f, dx, dy float64, w int) *Camera {
	h := int(math.Round(float64(w) * dy / dx))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	return &Camera{
		Transform: IDTransform, f: f, dx: dx, dy: dy,
		Width: w, Height: h, Image: img}
}

// BuildRay creates a Ray in global space given an image pixel position
func (c *Camera) BuildRay(x, y int) Ray {
	X := (float64(x)*c.dx)/float64(c.Width) - c.dx/2
	Y := c.dy/2 - float64(y)*c.dy/float64(c.Height)
	return c.RayToGlobal(Ray{pt: Origin, dir: Vector3{X, Y, -c.f}})
}

// Project ...
func (c *Camera) Project(p Point3) (int, int) {
	lp := c.PointToLocal(p)
	var x, y float64
	if lp[Z] < 0 {
		x = c.f * lp[X] / -lp[Z]
		y = c.f * lp[Y] / -lp[Z]
	} else {
		x = math.Copysign(c.dx/2, lp[X])
		y = math.Copysign(c.dy/2, lp[Y])
	}
	px := (x + c.dx/2) * float64(c.Width) / c.dx
	py := float64(c.Height) - (y+c.dy/2)*float64(c.Height)/c.dy
	return int(math.Round(px)), int(math.Round(py))
}

// Translate applies a translation
func (c *Camera) Translate(x, y, z float64) *Camera {
	c.Transform.Translate(x, y, z)
	return c
}

// RotateX applies a rotation around x-axis
func (c *Camera) RotateX(x float64) *Camera {
	c.Transform.RotateX(x)
	return c
}

// RotateY applies a rotation around y-axis
func (c *Camera) RotateY(y float64) *Camera {
	c.Transform.RotateY(y)
	return c
}

// RotateZ applies a rotation around z-axis
func (c *Camera) RotateZ(z float64) *Camera {
	c.Transform.RotateZ(z)
	return c
}

// Scale applies a scaling transform
func (c *Camera) Scale(x, y, z float64) *Camera {
	c.Transform.Scale(x, y, z)
	return c
}
