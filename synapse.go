package main

import gg "github.com/fogleman/gg"

// Used to connect neurons on the brain.
type Synapse struct {
	weight   float64
	to, from *Neuron
}

// http://neuralnetworksanddeeplearning.com/chap2.html
// Re-read this and employ the wisdom within.

func (s *Synapse) Tick() {
	if s.from == nil || s.to == nil {
		return
	}
	s.to.inputSum += s.weight * s.from.activation
}

func (s *Synapse) Draw(dc *gg.Context) {
	if s.from == nil || s.to == nil {
		return
	}
	if s.weight > 0 {
		dc.SetRGB(0, 0, s.weight)
	} else {
		dc.SetRGB(s.weight, 0, 0)
	}

	dc.SetLineWidth(2)
	dc.DrawLine(float64(s.from.x), float64(s.from.y), float64(s.to.x), float64(s.to.y))
	dc.Stroke()
}
