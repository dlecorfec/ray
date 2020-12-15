package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 1500)
	cam.Translate(5, 0, 70)
	cam.RotateX(-math.Pi / 12).RotateZ(math.Pi / 4)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1})
	l1.Translate(-3, -7, 25)

	p1 := ray.NewPlane().RotateX(math.Pi/4).Scale(300, 1, 300)
	p1.Surface = ray.Diffuse
	p1.Surface.Color.B = 1

	p2 := ray.NewPlane().Scale(100, 1, 300).
		RotateX(math.Pi/3).RotateY(math.Pi/12).RotateZ(math.Pi/12).Translate(0, 0, -50)

	s1 := ray.NewSphere().Scale(10, 10, 10)
	s1.Surface.Color.R = 1
	s2 := ray.NewSphere().Scale(10, 10, 10).Translate(-20, -20, 0)
	s2.Surface.Color.G = 1
	s3 := ray.NewSphere().Scale(10, 10, 10).Translate(-20, -50, 0)
	s3.Surface.Color.B = 1
	s4 := ray.NewSphere().Scale(10, 10, 10).Translate(20, -50, 0)
	s4.Surface.Color.R = 1
	s4.Surface.Color.B = 1
	s5 := ray.NewSphere().Scale(10, 10, 10).Translate(20, 0, 0)
	s5.Surface.Color.R = 1
	s5.Surface.Color.G = 1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(s1, p2)
	log.Printf("%v", s1.Surface.Color)
	log.Printf("%v", s2.Surface.Color)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
