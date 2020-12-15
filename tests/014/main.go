package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(12, 5, 4, 800)
	cam.RotateX(-math.Pi / 12).RotateY(math.Pi / 5).RotateZ(0)
	cam.Translate(60, 17, 70)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1.0, G: 0.5, B: 0}).Translate(20, 10, 70)

	groundsize := 100.0
	sol := ray.NewPlane().Scale(groundsize, groundsize, groundsize)
	sol.Surface = ray.Building

	globe := ray.NewSphere().Scale(200, 200, 200).Translate(0, 300, 0)
	globe.Surface = ray.White1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: 0.5}
	s.AddLights(l1)
	s.AddObjects(sol)

	nx, nz := 10, 10
	fnx, fnz := float64(nx), float64(nz)

	for i := 0; i < nx; i++ {
		for j := 0; j < nz; j++ {
			fi, fj := float64(i), float64(j)
			y := 1 + rand.Float64()
			obj := ray.NewCube().Scale(1, y, 1)
			x, z := 7*fi-3*fnx, 7*fj-3*fnz
			obj.Translate(x+rand.Float64(), y, z+rand.Float64())
			fmt.Printf("[%.1f;%.1f] ", x, z)
			obj.Surface = ray.Building
			s.AddObjects(obj)
		}
	}
	fmt.Printf("\n")
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
