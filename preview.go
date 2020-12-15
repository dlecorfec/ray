package ray

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Preview ...
type Preview struct {
	scene      *Scene
	preview    *ebiten.Image
	windowChan chan []pixel
	t          int
	lasty      int
	start      bool
	end        bool
	mode       int
}

// Update ...
func (pv *Preview) Update() error {
	if pv.end {
		//return fmt.Errorf("end")
	}
	switch {
	case inpututil.IsKeyJustReleased(ebiten.KeyEscape):
		return fmt.Errorf("esc")
	case inpututil.IsKeyJustReleased(ebiten.KeySpace):
		pv.mode ^= 1
	}
	return nil
}

// Draw ...
func (pv *Preview) Draw(screen *ebiten.Image) {
	pv.start = true
	pv.lasty = 0
	switch pv.mode {
	case 0:
		pv.drawRendering(screen, 0, 0, screen.Bounds().Max.X, screen.Bounds().Max.Y)
	case 1:
		pv.drawRendering(screen, screen.Bounds().Max.X/2, 0, screen.Bounds().Max.X, screen.Bounds().Max.Y/2)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f TPS: %f", ebiten.CurrentFPS(), ebiten.CurrentTPS()))
}

func (pv *Preview) drawRendering(screen *ebiten.Image, x, y, x2, y2 int) {
	stopdraw := false
	ratio := float64(x2-x) / float64(screen.Bounds().Max.X)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(ratio, ratio)
	opts.GeoM.Translate(float64(x), float64(y))

	if pv.end {
		screen.DrawImage(pv.preview, opts)
		return
	}
	c := time.After(15 * time.Millisecond)
	for {
		select {
		case <-c:
			stopdraw = true
			break
		case b := <-pv.windowChan:
			//log.Printf("drawing")
			for _, p := range b {
				select {
				case <-c:
					stopdraw = true
					break
				default:
				}
				ebitenutil.DrawRect(pv.preview, float64(p.x), float64(p.y), float64(p.w), float64(p.h), p.c)
				pv.lasty = p.y
			}
		default:
		}
		if stopdraw {
			break
		}
	}
	screen.DrawImage(pv.preview, opts)
	ebitenutil.DrawRect(screen, float64(x), float64(y)+float64(pv.lasty+1)*ratio, float64(screen.Bounds().Max.X), float64(1), color.RGBA{R: 0xff, A: 0xff})
}

// Layout ...
func (pv *Preview) Layout(outsideWidth, outsideHeight int) (int, int) {
	return pv.scene.cam.Width, pv.scene.cam.Height
}

func (pv *Preview) drawPixels(b []pixel) {
	if pv.preview != nil {
		pv.windowChan <- b
	}
}

func (pv *Preview) waitSetup() {
	if pv.preview == nil {
		return
	}
	for {
		if pv.start {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
}

func (pv *Preview) run() {
	if pv.preview == nil {
		for {
			if pv.end {
				break
			}
			time.Sleep(15 * time.Millisecond)
		}
		return
	}
	ebiten.SetWindowSize(pv.scene.cam.Width, pv.scene.cam.Height)
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(pv); err != nil {
		log.Printf("%s", err)
	}
}

func (pv *Preview) endRender() {
	close(pv.windowChan)
	pv.end = true
}

// newPreview ...
func newPreview(s *Scene) *Preview {
	var preview *ebiten.Image

	if s.Preview {
		preview = ebiten.NewImage(s.cam.Width, s.cam.Height)
	}
	return &Preview{scene: s, preview: preview, windowChan: make(chan []pixel, 2)}
}
