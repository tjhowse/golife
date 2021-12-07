package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"

	gg "github.com/fogleman/gg"
	"github.com/tjhowse/tjgo"
)

const WORLD_HEIGHT = 50
const WORLD_WIDTH = 50
const GRID_TO_PIXEL = 10

type World struct {
	crits  []*Crit
	frames int // The number of frames that have passed
	dc     *gg.Context
}

// Add a crit to the world.
func (w *World) AddCrit(c *Crit) {
	w.crits = append(w.crits, c)
}

// Call Tick() on all crits in the world, "ticks" times.
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

// Draw the world to the provided buffer as PNG.
func (w *World) Draw(o *bytes.Buffer) {
	// dc := gg.NewContext(WORLD_WIDTH*GRID_TO_PIXEL, WORLD_HEIGHT*GRID_TO_PIXEL)
	w.dc.SetRGB(1, 1, 1)
	w.dc.Clear()

	for i := 0; i < len(w.crits); i++ {
		w.crits[i].Draw(w.dc)
	}
	// dc.SetRGB(0, 0, 0)
	// dc.Fill()
	w.dc.EncodePNG(o)
}

// Handles serving the world as a png to the webserver
func (w *World) ImgHandler(response http.ResponseWriter, request *http.Request) {

	buff := bytes.Buffer{}
	w.Draw(&buff)
	response.Header().Set("Content-Length", fmt.Sprint(buff.Len()))
	response.Header().Set("Content-Type", "image/png")
	response.Write(buff.Bytes())
}

// Get a random position on the world that is not occupied by a crit.
func (w *World) GetRandomValidPosition() (x, y int) {
	for x, y = -1, -1; !w.CheckMove(x, y); {
		x = rand.Intn(WORLD_WIDTH)
		y = rand.Intn(WORLD_HEIGHT)
	}
	return x, y
}

// Cull some of the crits in the world based on cruel, arbitrary rules.
func (w *World) CullCrits(left, right int) {
	living := make([]*Crit, 0)
	for _, c := range w.crits {
		if c.pos.X < left || c.pos.X > right {
			living = append(living, c)
		}
	}
	w.crits = living
}

// Adds this many crits to the world with random connectomes.
func (w *World) AddRandomCrits(count int) {
	for i := 0; i < count; i++ {
		x, y := w.GetRandomValidPosition()
		var connectome Connectome
		connectome.Randomise()
		w.AddCrit(NewCrit(x, y, GRID_TO_PIXEL, w, connectome))
	}
}

// Randomise the position of all crits.
func (w *World) RandomiseCritPositions() {
	for i := 0; i < len(w.crits); i++ {
		w.crits[i].pos.X, w.crits[i].pos.X = w.GetRandomValidPosition()
	}
}

// This tops up the world to the provided number using mutated brains from
// surviving crits.
func (w *World) RefillCritsWithMutatedConnectomes(count int, mutations int) {

	critsToMake := count - len(w.crits)
	for i := 0; i < critsToMake; i++ {
		x, y := w.GetRandomValidPosition()
		var connectome Connectome
		if len(w.crits) > 0 {
			connectome.CopyFrom(&w.crits[rand.Intn(len(w.crits))].b.connectome)
			connectome.Mutate(mutations)
		} else {
			// Everyone died :(
			connectome.Randomise()
		}
		w.AddCrit(NewCrit(x, y, GRID_TO_PIXEL, w, connectome))
	}
}

func (w *World) GetCritIndexClosestToImgClick(x, y int) (int, error) {
	var closest int
	var closestDist = float64(1000000)
	if len(w.crits) == 0 {
		return -1, fmt.Errorf("no crits in world")
	}

	for i, c := range w.crits {
		critPos := c.pos
		critPos.Scale(10)
		dist := critPos.Distance(tjgo.Vec2{X: x, Y: y})
		if dist < closestDist {
			closest = i
			closestDist = dist
		}
	}
	return closest, nil
}
