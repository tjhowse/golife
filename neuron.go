package main

import (
	"fmt"
	"math"

	gg "github.com/fogleman/gg"
)

type ActivationFunction func()

type Neuron struct {
	activation  float64
	inputSum    float64
	bias        float64
	function    ActivationFunction
	x, y        int
	debugTop    string
	debugBottom string
	label       string
}

// Check if the neuron's bound function needs to be fired
// based on its activation level.
func (n *Neuron) Tick() {
	if n.function == nil {
		return
	}
	// Do the sigmoid calc here to turn inputSum into activation
	n.activation = 1 / (1 + math.Exp(-n.inputSum-n.bias))
	if n.activation >= SIGMOID_OUTPUT_ACTIVATION_THRESHOLD {
		n.function()
	}
	n.debugTop = fmt.Sprintf("%0.2f", n.activation)
	n.debugBottom = fmt.Sprintf("%0.2f", n.bias)
	n.inputSum = 0
}

func (n *Neuron) Draw(dc *gg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString(n.debugTop, float64(n.x)-15, float64(n.y)-20)
	dc.DrawString(n.debugBottom, float64(n.x)-15, float64(n.y)+25)
	dc.DrawString(n.label, float64(n.x)-15, float64(n.y)+50)
	dc.SetRGB(n.activation, 0, 0)
	dc.DrawCircle(float64(n.x), float64(n.y), float64(10))
	dc.Fill()
}
