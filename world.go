package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math/rand"
	"net/http"

	gg "github.com/fogleman/gg"
)

const WORLD_HEIGHT = 1000
const WORLD_WIDTH = 1000
const GRID_TO_PIXEL = 10

type World struct {
	crits  []*Crit
	frames int // The number of frames that have passed
}

func (w *World) AddCrit(c *Crit) {
	w.crits = append(w.crits, c)
}

func (w *World) Tick(ticks int) {
	for j := 0; j < ticks; j++ {
		for i := 0; i < len(w.crits); i++ {
			w.crits[i].Tick()
		}
		w.frames++
	}
}

// These "Check.*" methods are how a crit checks whether it can do what it wants.
func (w *World) CheckMove(x, y int) bool {
	if x < 0 || x >= WORLD_WIDTH || y < 0 || y >= WORLD_HEIGHT {
		return false
	}
	for _, c := range w.crits {
		if c.pos.X == x && c.pos.Y == y {
			return false
		}
	}
	return true
}

func (w *World) Draw(o *bytes.Buffer) {
	dc := gg.NewContext(WORLD_WIDTH, WORLD_HEIGHT)

	for i := 0; i < len(w.crits); i++ {
		w.crits[i].Draw(dc)
	}
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.EncodePNG(o)
}

func (w *World) ImgHandler(response http.ResponseWriter, request *http.Request) {

	buff := bytes.Buffer{}
	w.Draw(&buff)
	response.Header().Set("Content-Length", fmt.Sprint(buff.Len()))
	response.Header().Set("Content-Type", "image/png")
	response.Write(buff.Bytes())
}

func (w *World) GetRandomValidPosition() (x, y int) {
	for x, y = -1, -1; !w.CheckMove(x, y); {
		x = rand.Intn(WORLD_WIDTH/GRID_TO_PIXEL) * GRID_TO_PIXEL
		y = rand.Intn(WORLD_WIDTH/GRID_TO_PIXEL) * GRID_TO_PIXEL
	}
	return x, y
}

func (w *World) CullCrits() {
	living := make([]*Crit, 0)
	for _, c := range w.crits {
		if c.pos.X > WORLD_WIDTH/2 {
			living = append(living, c)
		}
	}
	w.crits = living
}

func (w *World) AddRandomCrits(count int) {
	for i := 0; i < count; i++ {
		x, y := w.GetRandomValidPosition()
		var connectome Connectome
		connectome.Randomise()
		w.AddCrit(NewCrit(x, y, GRID_TO_PIXEL, color.Black, w, connectome))
	}
}

// This tops up the world to the provided number using mutated brains from
// surviving crits.
func (w *World) RefillCritsWithMutatedConnectomes(count int) {
	critsToMake := count - len(w.crits)
	for i := 0; i < critsToMake; i++ {
		x, y := w.GetRandomValidPosition()
		var connectome Connectome
		connectome.CopyFrom(&w.crits[rand.Intn(len(w.crits))].b.connectome)
		connectome.Mutate(10)
		w.AddCrit(NewCrit(x, y, GRID_TO_PIXEL, color.Black, w, connectome))
	}
}
