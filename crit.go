package main

import (
	"image/color"

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

func NewCrit(x, y, size int, c color.Color, w *World, connectome [32]byte) *Crit {
	b := NewBrain(connectome)
	s := make([]float64, len(b.inputNeurons))
	crit := Crit{tj.Vec2{X: x, Y: y}, GRID_TO_PIXEL / 2, c, w, b, s}
	b.outputNeurons[0].function = crit.MoveUp
	b.outputNeurons[1].function = crit.MoveDown
	b.outputNeurons[2].function = crit.MoveLeft
	b.outputNeurons[3].function = crit.MoveRight
	return &crit
}

func (c *Crit) MoveUp() {
	if c.w.CheckMove(c.pos.X, c.pos.Y-1) {
		c.pos.Y -= GRID_TO_PIXEL
	}
}

func (c *Crit) MoveDown() {
	if c.w.CheckMove(c.pos.X, c.pos.Y+1) {
		c.pos.Y += GRID_TO_PIXEL
	}
}

func (c *Crit) MoveLeft() {
	if c.w.CheckMove(c.pos.X-1, c.pos.Y) {
		c.pos.X -= GRID_TO_PIXEL
	}
}

func (c *Crit) MoveRight() {
	if c.w.CheckMove(c.pos.X+1, c.pos.Y) {
		c.pos.X += GRID_TO_PIXEL
	}
}

func (c *Crit) Draw(dc *gg.Context) {
	dc.SetColor(c.c)
	dc.DrawCircle(float64(c.pos.X), float64(c.pos.Y), float64(c.size))
}

func (c *Crit) Tick() {
	// Senses should probably go here? The crit can map between the world and the brain's
	// input neurons.

	// Toggle sense 0 on and off every frame.
	if c.w.frames%2 == 0 {
		c.s[0] = 1
	} else {
		c.s[0] = 0
	}
	c.b.Tick(c.s)
}
