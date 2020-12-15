package main

import (
	"log"
	"math"
	"math/rand"
	"os"

	"runtime/pprof"

	"github.com/dlecorfec/ray"
)

func weirdObj() *ray.BoundingBox {
	s := ray.NewBoundingBox()
	rd := rand.New(rand.NewSource(2))
	rdc := rand.New(rand.NewSource(1))

	nsph := 40
	nbr := 12
	for j := 0; j < nbr; j++ {
		fj := float64(j)
		for i := 0; i < nsph; i++ {
			fi := float64(i)
			maxp := float64(nsph)
			rcx := float64(30)
			rcz := float64(20)
			sph := ray.NewSphere().Scale(1+rd.Float64(), 1+rd.Float64(), 1+rd.Float64())
			phase := fj * 2 * math.Pi / float64(nbr)
			sph.Translate((math.Pow(fi, 1.5)/maxp)*rcx*math.Cos(fi*math.Pi/maxp+phase), 0, (math.Pow(fi, 1.5)/maxp)*rcz*math.Sin(fi*math.Pi/maxp+phase))
			sph.Surface = ray.White1
			if rdc.Int31n(2) > 0 {
				sph.Surface.Color = ray.FloatColor{R: 2, G: 0.3, B: 0.1}
			}
			s.AddObjects(sph)
		}
	}

	nx, nz := 10, 10
	fx, fz := float64(nx), float64(nz)

	globe := ray.NewSphere().Scale(6, 6, 6)
	globe.Surface.Ks = 0.1
	globe.Surface.Color = ray.FloatColor{R: 2, G: .2, B: .2}
	bbpikes := ray.NewBoundingBox()
	bbpikes.AddObjects(globe)
	s.AddObjects(globe)
	rand.Seed(1)
	for i := 0; i < nx; i++ {
		fi := float64(i)
		for j := 0; j < nz; j++ {
			fj := float64(j)
			obj := ray.NewCube().Scale(.3, 10+rand.Float64(), .3)
			obj.RotateX(math.Pi*fi/fx + rand.Float64()/20 - 1/10)
			obj.RotateY(rand.Float64()/20 - 1/10)
			obj.RotateZ(math.Pi*fj/fz + rand.Float64()/20 - 1/10)
			obj.Surface = ray.White1
			obj.Surface.Ks = 0.1
			draw := true

			if (i == 4 || i == 6) && j%2 == 0 {
				draw = false
			}
			if i == 5 && j != 0 {
				draw = false
			}
			if draw {
				bbpikes.AddObjects(obj)
				//s.AddObjects(obj)
			}
		}
	}
	s.AddObjects(bbpikes)
	return s
}

func main() {
	f, err := os.Create("015-cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// camera/scene
	cam := ray.NewCamera(11, 8, 5, 1280)
	//cam := ray.NewCamera(10, 5, 4, 800)

	cam.RotateX(-math.Pi / 20).RotateY(math.Pi / 4.5).RotateZ(0)
	cam.Translate(100, 21, 110)
	s := ray.NewScene(cam)
	s.Ambiant = ray.FloatColor{R: .5, G: .5, B: 0.5}

	// lights
	l1 := ray.NewPointLight(ray.FloatColor{R: 1, G: 1, B: 1}).Translate(20, 40, 30)
	s.AddLights(l1)
	//dl := ray.NewSphere().Translate(11.5, 11.5, 31.5)
	//dl.Surface.Color = ray.FloatColor{R: 0, G: 1, B: 0}
	//s.AddObjects(dl)
	// objects
	groundsize := 200.0
	sol := ray.NewPlane().Scale(groundsize, groundsize, groundsize)
	sol.RotateX(math.Pi/2).Translate(0, 0, -40)
	//sol.Surface.Color = ray.FloatColor{R: 3, G: 3, B: 0.1}
	sol.Surface.Color = ray.FloatColor{R: 3, G: 2, B: 0.5}
	sol.Surface.Ks = 0
	//s.AddObjects(sol)

	w := weirdObj()
	w.RotateZ(-math.Pi / 6)
	w.RotateY(-math.Pi / 6)
	w.RotateX(math.Pi / 6)
	s.AddObjects(w)
	s.Raytrace()
	err = s.WritePNG("")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
