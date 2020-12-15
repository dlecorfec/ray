package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 500)
	cam.RotateX(-math.Pi / 4)
	cam.Translate(-20, 20, 12)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(-5, 20, 0)

	s1 := ray.NewSphere()
	//s1.Rotate(0, 0, math.Pi/4)
	//s1.Scale(5, 5, 5)
	s1.Translate(-10, 10, 0)
	s1.Surface.Color.R = 1

	p1 := ray.NewPlane()
	p1.Scale(100, 100, 100)

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(s1, p1)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
