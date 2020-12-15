package main

import (
	"log"
	"math"

	"github.com/dlecorfec/ray"
)

func main() {
	cam := ray.NewCamera(22, 16, 9, 1980)
	cam.Translate(0, 2, 30)
	cam.RotateX(-math.Pi / 8).RotateY(3 * math.Pi / 10).RotateZ(0)

	//l1 := ray.NewPointLight(ray.Red).Translate(0.6, 8, 0)
	//l2 := ray.NewPointLight(ray.Green).Translate(-0.6, 8, 0)
	//l3 := ray.NewPointLight(ray.Blue).Translate(0, 8, 0.6)
	l4 := ray.NewPointLight(ray.White).Translate(0, 8, 0.6)
	l5 := ray.NewPointLight(ray.FloatColor{R: 1, G: .2, B: .1}).Translate(-6, 7, 10)
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

	c1 := ray.NewCube().RotateX(-math.Pi/12).RotateY(-math.Pi/12).Translate(2, 1.5, -1.5)
	c1.Surface = ray.White1

	s1 := ray.NewSphere().Translate(-2, 1, 0.5)
	s1.Surface = ray.White1

	s2 := ray.NewSphere().Scale(2, 2, 2).Translate(-2, 2, 5)
	s1.Surface = ray.White1

	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: 0.2, G: 0.2, B: 0.2}
	s.AddLights(l4, l5)
	s.AddObjects(c1, s1, s2, murG, murD, sol)
	s.Raytrace()
	err := s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
