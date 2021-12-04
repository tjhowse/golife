package main

import (
	"image/color"
	"math/rand"

	gg "github.com/fogleman/gg"
	tj "github.com/tjhowse/tjgo"
)

type Crit struct {
	pos  tj.Vec2
	size int
	c    color.Color
	w    *World
	b    *Brain
	s    []float64
}

// A factory function for creating crits.
func NewCrit(x, y, size int, w *World, connectome Connectome) *Crit {
	b := NewBrain(connectome)
	s := make([]float64, len(b.inputNeurons))
	// c := color.Black
	// Derive a colour from a connectome

	var red uint8
	var green uint8
	var blue uint8
	for i := 0; i < len(connectome.c)-3; i += 3 {
		red += uint8(connectome.c[i])
		green += uint8(connectome.c[i+1])
		blue += uint8(connectome.c[i+2])
	}
	c := color.RGBA{red, green, blue, 255}
	crit := Crit{tj.Vec2{X: x, Y: y}, GRID_TO_PIXEL / 2, c, w, b, s}
	b.outputNeurons[0].function = crit.MoveUp
	b.outputNeurons[1].function = crit.MoveDown
	b.outputNeurons[2].function = crit.MoveLeft
	b.outputNeurons[3].function = crit.MoveRight
	b.outputNeurons[4].function = crit.MoveRandomly
	return &crit
}

// Movement functions to be bound to output neurons.
func (c *Crit) MoveUp() {
	if c.w.CheckMove(c.pos.X, c.pos.Y-1) {
		c.pos.Y -= 1
	}
}

func (c *Crit) MoveDown() {
	if c.w.CheckMove(c.pos.X, c.pos.Y+1) {
		c.pos.Y += 1
	}
}

func (c *Crit) MoveLeft() {
	if c.w.CheckMove(c.pos.X-1, c.pos.Y) {
		c.pos.X -= 1
	}
}

func (c *Crit) MoveRight() {
	if c.w.CheckMove(c.pos.X+1, c.pos.Y) {
		c.pos.X += 1
	}
}

func (c *Crit) MoveRandomly() {
	switch rand.Intn(4) {
	case 0:
		c.MoveUp()
	case 1:
		c.MoveDown()
	case 2:
		c.MoveLeft()
	case 3:
		c.MoveRight()
	}
}

// Draw the crit on the given context.
func (c *Crit) Draw(dc *gg.Context) {
	dc.SetColor(c.c)
	dc.DrawCircle(float64(c.pos.X*GRID_TO_PIXEL), float64(c.pos.Y*GRID_TO_PIXEL), float64(c.size))
	dc.Fill()
}

// Update the crit's body and brain.
func (c *Crit) Tick() {
	// Toggle sense 0 on and off every frame.
	if c.w.frames%2 == 0 {
		c.s[0] = 1
	} else {
		c.s[0] = 0
	}
	// TODO Add more senses
	c.b.Tick(c.s)
}
