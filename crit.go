package main

import (
	"crypto/rand"
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

func NewCrit(x, y, size int, c color.Color, w *World) *Crit {
	connectome := [32]byte{}
	rand.Read(connectome[:])
	b := NewBrain(connectome)
	s := make([]float64, len(b.inputNeurons))
	crit := Crit{tj.Vec2{X: x, Y: y}, size, c, w, b, s}
	b.outputNeurons[0].function = crit.MoveUp
	b.outputNeurons[1].function = crit.MoveDown
	b.outputNeurons[2].function = crit.MoveLeft
	b.outputNeurons[3].function = crit.MoveRight
	return &crit
}

func (c *Crit) MoveUp() {
	c.pos.Y -= 10
}

func (c *Crit) MoveDown() {
	c.pos.Y += 10
}

func (c *Crit) MoveLeft() {
	c.pos.X -= 10
}

func (c *Crit) MoveRight() {
	c.pos.X += 10
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
