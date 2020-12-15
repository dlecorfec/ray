package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(16, 16, 9, 800).RotateX(-math.Pi/16).Translate(10, 80, 400)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1}).Translate(0, 0, 400)

	bottom := ray.NewPlane().Scale(100, 1, 100).Translate(0, -100, 0)
	back := ray.NewPlane().Scale(100, 1, 100).RotateX(math.Pi/2).Translate(0, 0, -100)
	top := ray.NewPlane().Scale(100, 1, 100).Translate(0, 100, 0)

	log.Printf("%v", back.Transform.Check())

	s1 := ray.NewSphere().Scale(15, 15, 15)
	s1.Surface.Color.R = 1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(back, top, bottom, s1)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
