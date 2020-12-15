package ray

import (
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
)

// MaxDepth is the max tracing recursion level
var MaxDepth = 7

// Scene contains the objects, lights, camera and method to render them
type Scene struct {
	MaxDepth     int
	cam          *Camera
	lights       []Light
	objects      []Object
	Ambiant      FloatColor
	raysPerDepth []int
	traceChan    chan []pixel
	drawChan     chan []pixel
	num          int
	lasty        int
	Preview      bool
}

type pixel struct {
	x int
	y int
	w int
	h int
	c FloatColor
}

// NewScene instantiates a scene with a Camera
func NewScene(cam *Camera) *Scene {
	s := &Scene{
		MaxDepth:     MaxDepth,
		cam:          cam,
		lights:       make([]Light, 0),
		objects:      make([]Object, 0),
		raysPerDepth: make([]int, MaxDepth+1),
		Preview:      true,
	}
	return s
}

// AddLights adds the lights to the scene
func (s *Scene) AddLights(list ...Light) {
	s.lights = append(s.lights, list...)
}

// AddObjects adds the objects to the scene
func (s *Scene) AddObjects(list ...Object) {
	s.objects = append(s.objects, list...)
}

func (s *Scene) linearTracing() {
	pb := newPixelBatch(s.traceChan, 16)
	for y := 0; y < s.cam.Height; y++ {
		for x := 0; x < s.cam.Width; x++ {
			pb.add(pixel{x: x, y: y, w: 1, h: 1})
		}
	}
	pb.flush()
}

type pixelBatch struct {
	b    []pixel
	c    chan []pixel
	size int
}

func newPixelBatch(c chan []pixel, size int) *pixelBatch {
	return &pixelBatch{c: c, size: size}
}

func (pb *pixelBatch) add(p pixel) {
	pb.b = append(pb.b, p)
	if len(pb.b) == pb.size {
		pb.flush()
	}
}

func (pb *pixelBatch) flush() {
	pb.c <- pb.b
	pb.b = nil
}

func (s *Scene) progressiveTracingBatch() {
	max := s.cam.Width
	if s.cam.Height > s.cam.Width {
		max = s.cam.Height
	}

	pow := 1
	for t := max; t > 1; t /= 2 {
		pow *= 2
	}

	pb := newPixelBatch(s.traceChan, 16)
	pb.add(pixel{x: 0, y: 0, w: s.cam.Width, h: s.cam.Height})
	for mod := pow; mod > 0; mod /= 2 {
		for y := 0; y < s.cam.Height; y += mod {
			for x := 0; x < s.cam.Width; x += mod {
				if x%(2*mod) == 0 && y%(2*mod) == 0 {
					continue
				}
				pb.add(pixel{x: x, y: y, w: mod, h: mod})
			}
		}
	}
	pb.flush()
}

// Raytrace ...
func (s *Scene) Raytrace() {
	s.Preview = true
	s.traceChan = make(chan []pixel, 1000)
	s.drawChan = make(chan []pixel, 1000)

	var wg sync.WaitGroup
	// start trace workers
	for i := 0; i < runtime.NumCPU()/4; i++ {
		go s.traceWorker(&wg)
		wg.Add(1)
	}

	/*
		runtime.SetBlockProfileRate(100)
		f, err := os.Create("block")
		if err != nil {
			// Error handling
		}

		p := pprof.Lookup("block")
		defer func() {
			err := p.WriteTo(f, 0)
			if err != nil {
				log.Printf("Error writing block profile: %v", err)
			}
		}()
	*/
	pv := newPreview(s)
	// start draw worker
	go s.drawWorker(pv)
	// send screen coords to workers
	go func() {
		pv.waitSetup()
		s.progressiveTracingBatch()
		close(s.traceChan)
		wg.Wait()
		log.Printf("rays per depth: %v", s.raysPerDepth)
		close(s.drawChan)
	}()
	pv.run()
}

func (s *Scene) traceWorker(wg *sync.WaitGroup) {
	pb := newPixelBatch(s.drawChan, 16)
	for b := range s.traceChan {
		for _, p := range b {
			r := s.cam.BuildRay(p.x, p.y)
			r.x, r.y = p.x, p.y
			r.Normalize()
			p.c = s.trace(r, 0)
			pb.add(p)
		}
	}
	pb.flush()
	wg.Done()
}

func (s *Scene) drawWorker(pv *Preview) {
	for b := range s.drawChan {
		for _, p := range b {
			s.num++
			s.cam.Image.SetRGBA(p.x, p.y, p.c.Color())
			s.lasty = p.y
		}
		pv.drawPixels(b)
	}
	pv.endRender()
}

// Background ...
func (s *Scene) Background() FloatColor {
	return FloatColor{0.1, 0.1, 0.1}
}

// trace a ray
// r: normalized ray in global space
func (s *Scene) trace(r Ray, depth int) FloatColor {
	if depth > s.MaxDepth {
		return s.Background()
	}
	//s.raysPerDepth[depth]++
	hit := s.findIntersection(r)
	if hit == nil {
		return s.Background()
	}
	//log.Printf("scene: %#v %#v\n", obj, sd)
	c := s.whitted(r, hit)
	//log.Printf("--- %d,%d=%v", x, y, c)
	if hit.Surface.Ks > 0 {
		refl := s.reflection(hit, depth)
		//log.Printf("TRACE --- %d,%d=%v %#v %d", x, y, refl, hit, depth)
		c = c.Add(refl)
	}
	return c
}

func (s *Scene) reflection(h *Hit, depth int) FloatColor {
	if depth+1 > s.MaxDepth {
		return FloatColor{}
	}

	cosNI := -(h.globNorm.dir.Dot(h.globRay.dir))
	//log.Printf("cosNIAt(%d,%d,%v,%v)=%f", x, y, h.globNorm.dir, h.globRay.dir, cosNI)
	newRay := Ray{
		pt:  h.globNorm.pt,
		dir: h.globRay.dir.Add(h.globNorm.dir.Mult(2 * cosNI)),
	}
	newRay.Normalize()
	rc := s.trace(newRay, depth+1)
	c := rc.MulF(h.Surface.Ks)
	return c
}

func (s *Scene) whitted(r Ray, h *Hit) FloatColor {
	c := h.Surface.ColorAt(h)
	// ambiant
	a := c.MulC(s.Ambiant)
	//log.Printf("%v", h.Surface.Ka)
	a = a.MulF(h.Surface.Ka)
	wc := a
	for _, li := range s.lights {
		//log.Printf("norm=%v", h.globNorm)
		rl := li.RayToLight(h.globNorm.pt)
		dist := rl.dir.Norm()
		if dist < Epsilon {
			if s.debug(r) {
				//log.Printf("dist < Epsilon")
			}
			continue
		}
		rl.Normalize()
		vl := rl.dir
		cosNL := h.globNorm.dir.Dot(vl)
		//log.Printf("vl=%v norm=%v cosNL=%f", vl, h.globNorm.dir, cosNL)
		if cosNL < Epsilon {
			if s.debug(r) {
				//log.Printf("cosNL < Epsilon")
			}
			continue
		}
		//cosNL = 1
		// shadow?
		if s.debug(r) {
			log.Printf("vl=%v norm=%v cosNL=%f", vl, h.globNorm.dir, cosNL)
			//log.Printf("hidden %v %f", rl, dist)
		}
		rl.x, rl.y = r.x, r.y
		if s.isHidden(rl, dist) {
			if s.debug(rl) {
				//log.Printf("hidden")
			}
			continue
		}

		// diffuse term
		fatt := math.Exp(-.01 * dist)
		if li.Sun() {
			fatt = 1
		}
		//fatt := 1.0
		diffuse := li.Color(rl).MulC(c).MulF(h.Surface.Kd).MulF(cosNL).MulF(fatt)
		_ = diffuse
		//log.Printf("a=%v diffuse=%v liR=%v cR=%v Kd=%v cos=%v", a, diffuse, li.Color(rl).R, c.R, h.Surface.Kd, cosNL)
		wc.Add(diffuse)

		// specular term (phong)
		vr := h.globNorm.dir.Mult(2 * cosNL).Sub(vl)
		cosRO := vr.Dot(r.dir)
		if s.debug(r) {
			//log.Printf("vr=%v r=%v cosRO=%f cosNL=%f normale=%v vl=%v", vr, r.dir, cosRO, cosNL, h.globNorm.dir, vl)
		}
		if cosRO > 0 {
			continue
		}
		specular := li.Color(rl).MulF(h.Surface.Ks).MulF(math.Pow(cosRO, float64(h.Surface.Nphong))).MulF(fatt)
		_ = specular
		//log.Printf("r=%v cosRO=%f pow=%f", r.dir, cosRO, math.Pow(cosRO, float64(h.Surface.Nphong)))
		wc.Add(specular)
	}

	return wc
}

func (*Scene) debug(r Ray) bool {
	if r.x == 360 && (r.y == 200 || r.y == 250) {
		return true
	}
	if r.x == 4000 && r.y == 225 {
		return true
	}
	return false
}

// does rl intersect an object closer than dist?
func (s *Scene) isHidden(rl Ray, dist float64) bool {
	for _, obj := range s.objects {
		h := obj.Intersect(rl)
		if h == nil {
			continue
		}
		v := Vector3{h.globRay.pt[X] - rl.pt[X], h.globRay.pt[Y] - rl.pt[Y], h.globRay.pt[Z] - rl.pt[Z]}
		if s.debug(rl) {
			log.Printf("isHidden by %s, inter=%v rl=%v v=%v, |v|=%f", obj.Name(), h.globRay.pt, rl.pt, v, v.Norm())
		}

		if v.Norm() < dist {
			return true
		}
	}
	return false
}

// findIntersection finds the closest intersecting object, if any.
func (s *Scene) findIntersection(r Ray) *Hit {
	minDist := math.MaxFloat64
	var h *Hit
	for _, o := range s.objects {
		s := o.Intersect(r)
		if s != nil {
			//dist := SquareDist(r.pt, gn.pt)
			dist := r.pt.SquareDist(s.globNorm.pt)
			if dist < minDist {
				minDist = dist
				h = s
			}
		}
	}
	return h
}

// WriteJPG ...
func (s *Scene) WriteJPG(name string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, s.cam.Image, nil)
}

// WritePNG ...
func (s *Scene) WritePNG(name string) error {
	if name == "" {
		if len(os.Args) > 1 {
			name = os.Args[1]
		}
	}
	var out io.Writer
	if name != "" {
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
		defer os.Stdout.Sync()
	}
	return png.Encode(out, s.cam.Image)
}
