package main

import (
	"log"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(5, 16, 9, 800)
	cam.Translate(0, 0, 20)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(0, 50, 0)

	s1 := ray.NewSphere()
	s1.Scale(5, 5, 5)
	s1.Surface.Color.R = 1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(s1)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
