package ray

import (
	"image/color"
	"math"
)

// FloatColor implements the color.Color interface by having a RGBA() method
type FloatColor struct {
	R float64
	G float64
	B float64
}

// Red ...
var Red = FloatColor{R: 1}

// Green ...
var Green = FloatColor{G: 1}

// Blue ...
var Blue = FloatColor{B: 1}

// White ...
var White = FloatColor{R: 1, G: 1, B: 1}

// RGBA return the color as RGBA
func (fc FloatColor) RGBA() (uint32, uint32, uint32, uint32) {
	r := math.Max(fc.R, 0)
	g := math.Max(fc.G, 0)
	b := math.Max(fc.B, 0)
	r = math.Min(r, 1)
	g = math.Min(g, 1)
	b = math.Min(b, 1)
	return uint32(math.Round(float64(0xffff) * r)),
		uint32(math.Round(float64(0xffff) * g)),
		uint32(math.Round(float64(0xffff) * b)),
		0xffff
}

// Color ...
func (fc FloatColor) Color() color.RGBA {
	r, g, b, _ := fc.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
}

// MulC multiplies each RGB component of a color by another color component
func (fc FloatColor) MulC(b FloatColor) FloatColor {
	fc.R *= b.R
	fc.G *= b.G
	fc.B *= b.B
	return fc
}

// Add adds each RGB component of a color with another color component
func (fc *FloatColor) Add(b FloatColor) FloatColor {
	fc.R += b.R
	fc.G += b.G
	fc.B += b.B
	return *fc
}

// MulF multiplies each RGB component by a scalar
func (fc FloatColor) MulF(b float64) FloatColor {
	fc.R *= b
	fc.G *= b
	fc.B *= b
	return fc
}
