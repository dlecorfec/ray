package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(10, 16, 9, 800).RotateY(math.Pi/12).Translate(0, 0, 70)

	l1 := ray.NewPointLight(ray.FloatColor{R: 2, G: 2, B: 2})
	//l1.Translate(50, 100, 50)
	l1.Translate(10, 10, 100)
	sa := make([]ray.Object, 0, 20)
	for i := -80; i < 90; i += 20 {
		for j := -80; j < 90; j += 20 {
			s := ray.NewSphere()
			s.Scale(10, 10, 10)
			//s.Scale(2+8*rand.Float64(), 2+8*rand.Float64(), 2+8*rand.Float64())
			//s.Rotate(rand.Float64()*math.Pi, rand.Float64()*math.Pi, rand.Float64()*math.Pi)
			s.Translate(float64(i), float64(j), 0)
			s.Surface.Color.R = float64(i+80) / 170
			s.Surface.Color.G = float64(j+80) / 340
			s.Surface.Color.B = rand.Float64()
			s.Surface.Ks = 0.3 + rand.Float64()/6
			s.Surface.Nphong = float64(2 + (rand.Intn(30)/2)*2)
			sa = append(sa, s)
		}
	}

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: .5}
	s.AddLights(l1)
	s.AddObjects(sa...)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
