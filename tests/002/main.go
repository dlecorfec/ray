package main

import (
	"log"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 800)
	cam.Translate(0, 0, 70)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(0, 100, 0)

	s1 := ray.NewSphere()
	s1.Scale(10, 10, 10)
	s1.Surface.Color.R = 1
	s2 := ray.NewSphere()
	s2.Scale(10, 10, 10)
	s2.Translate(-20, 0, 0)
	s2.Surface.Color.G = 1
	s3 := ray.NewSphere()
	s3.Scale(10, 10, 10)
	s3.Translate(20, 0, 0)
	s3.Surface.Color.B = 1
	s4 := ray.NewSphere()
	s4.Scale(10, 10, 10)
	s4.Translate(0, 20, 0)
	s4.Surface.Color.R = 1
	s4.Surface.Color.B = 1
	s5 := ray.NewSphere()
	s5.Scale(10, 10, 10)
	s5.Translate(20, 20, 0)
	s5.Surface.Color.R = 1
	s5.Surface.Color.G = 1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(s1, s2, s3, s4, s5)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
