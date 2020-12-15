package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(12, 5, 4, 800)
	cam.Translate(0, 3, 30)
	cam.RotateX(-math.Pi / 8).RotateY(3 * math.Pi / 10).RotateZ(0)

	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 0, B: 0}).Translate(0.6, 8, 0)
	l2 := ray.NewPointLight(ray.FloatColor{R: 0, G: 1, B: 0}).Translate(-0.6, 8, 0)
	l3 := ray.NewPointLight(ray.FloatColor{R: 0, G: 0, B: 1}).Translate(0, 8, 0.6)

	roomsize := 3.5
	groundsize := 21.0

	murG := ray.NewPlane().Scale(roomsize, roomsize, roomsize)
	murG.RotateZ(math.Pi / 2)
	murG.Translate(-roomsize, roomsize, 0)
	murG.Surface = ray.Mirror

	murD := ray.NewPlane().Scale(roomsize, roomsize, roomsize)
	murD.RotateX(math.Pi / 2)
	murD.Translate(0, roomsize, -roomsize)
	murD.Surface = ray.Mirror

	sol := ray.NewPlane().Scale(groundsize, groundsize, groundsize)
	sol.Surface = ray.Ocher2

	s1 := ray.NewSphere().Translate(0, 1, 0)
	s1.Surface = ray.White1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: 0, G: 0, B: 0}
	s.AddLights(l1, l2, l3)
	s.AddObjects(s1, murG, murD, sol)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
