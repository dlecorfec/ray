package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 800)
	cam.Translate(0, 0, 70)
	//cam.Rotate(-math.Pi/12, 0, math.Pi/4)
	//cam.Rotate(-math.Pi/12, 0, 0)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(0, 0, 25)

	s1 := ray.NewPlane()
	s1.RotateX(math.Pi / 2)
	s1.Scale(10, 10, 10)
	s1.Surface.Color.R = 1

	s := ray.NewScene(cam)
	//s.Ambiant = ray.FloatColor{R: .0, G: .0, B: .0}
	s.Ambiant = ray.FloatColor{R: 0.5, G: 0.5, B: .5}
	_ = l1
	s.AddLights(l1)
	s.AddObjects(s1)
	log.Printf("%v", s1.Surface.Color)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
