package display

import (
	"github.com/kierdavis/anna/analyser"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const NumRecent = 4

type Display struct {
	window *sdl.Window
	renderer *sdl.Renderer
	recent [NumRecent][]float64
	nextRecent int
	average []float64
	points []sdl.Point
}

func New() (d *Display, err error) {
	d = &Display{}
	
	err = d.init()
	if err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

func (d *Display) init() (err error) {
	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}
	
	d.window, err = sdl.CreateWindow("anna", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1400, 500, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	
	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	
	return nil
}

func (d *Display) Close() {
	if d.renderer != nil {
		d.renderer.Destroy()
	}
	
	if d.window != nil {
		d.window.Destroy()
	}
	
	sdl.Quit()
}

func (d *Display) Draw(an analyser.Analysis) (err error) {
	// Fill black
	err = d.renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {return err}
	
	err = d.renderer.Clear()
	if err != nil {return err}
	
	// Prepare dataset
	start := 2
	end := 600
	
	ampls := d.getNextRecentDataset(end - start)
	j := 0
	
	for i := start; i < end; i++ {
		_, ampl, _ := an.Info(i)
		ampls[j] = ampl * 1.75
		j++
	}
	
	// Plot this dataset
	err = d.renderer.SetDrawColor(0x00, 0x66, 0x00, 0xFF)
	if err != nil {return err}
	
	err = d.plot(ampls)
	if err != nil {return err}
	
	// Plot recent average dataset
	d.calcAverage()
	
	err = d.renderer.SetDrawColor(0x00, 0xFF, 0x00, 0xFF)
	if err != nil {return err}
	
	err = d.plot(d.average)
	if err != nil {return err}
	
	d.renderer.Present()
	return nil
}

func (d *Display) plot(dataset []float64) (err error) {
	if d.points == nil {
		d.points = make([]sdl.Point, len(dataset))
	}
	
	width, height := d.window.GetSize()
	marginX := 20
	marginY := 20
	stepX := float64(width - 2*marginX) / float64(len(dataset) - 1)
	rangeY := float64(height - 2*marginY)
	
	for i, ampl := range dataset {
		if ampl > 1.0 {
			ampl = 1.0
		}
		x := marginX + round(float64(i) * stepX)
		y := marginY + round((1.0 - ampl) * rangeY)
		d.points[i] = sdl.Point{int32(x), int32(y)}
	}
	
	return d.renderer.DrawLines(d.points)
}

func (d *Display) getNextRecentDataset(size int) (dataset []float64) {
	if d.recent[d.nextRecent] == nil {
		dataset = make([]float64, size)
		d.recent[d.nextRecent] = dataset
	} else {
		dataset = d.recent[d.nextRecent]
	}
	d.nextRecent = (d.nextRecent + 1) % len(d.recent)
	return dataset
}

func (d *Display) calcAverage() {
	if d.average == nil {
		d.average = make([]float64, len(d.recent[0]))
	}
	n := 0
	for _, ds := range d.recent {
		if ds != nil {
			n++
			for i, x := range ds {
				d.average[i] += x
			}
		}
	}
	nn := float64(n)
	for i := range d.average {
		d.average[i] /= nn
	}
}

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}
