package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 800)
	cam.RotateX(-math.Pi / 12).RotateZ(math.Pi / 4)
	//cam.Rotate(-math.Pi/12, 0, 0)
	cam.Translate(5, 0, 70)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(-3, -7, 25)

	p1 := ray.NewPlane()
	p1.Surface = ray.Diffuse
	p1.Surface.Color.B = 1
	p1.RotateX(math.Pi / 4)
	p1.Scale(30, 30, 30)

	s1 := ray.NewSphere()
	s1.Scale(10, 10, 10)
	s1.Surface.Color.R = 1
	s2 := ray.NewSphere()
	s2.Translate(-3, -7, 25)
	s2.Surface.Color.R = 1
	s2.Surface.Color.G = 1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .0, G: .0, B: .0}
	s.AddLights(l1)
	s.AddObjects(s1, s2)
	log.Printf("%v", s1.Surface.Color)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
