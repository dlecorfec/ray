package main

import (
	"log"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(16, 16, 9, 800).Translate(0, 0, 30)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1}).Translate(30, 0, 30)

	//bottom := ray.NewPlane().RotateX(math.Pi).Scale(100, 1, 100).Translate(0, -100, 0)
	//back := ray.NewPlane().RotateX(math.Pi/2).Scale(100, 1, 100).Translate(0, 0, -100)
	//top := ray.NewPlane().Scale(100, 1, 100).Translate(0, 100, 0)

	s1 := ray.NewSphere().Scale(6, 6, 6)
	s1.Surface.Color.R = 1

	s := ray.NewScene(cam)
	//s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(s1)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
