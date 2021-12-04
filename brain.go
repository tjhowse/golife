package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"

	gg "github.com/fogleman/gg"
)

const BRAIN_IMG_WIDTH = 500
const BRAIN_IMG_HEIGHT = 500

type ActivationFunction func()

type Neuron struct {
	activation float64
	threshold  float64
	function   ActivationFunction
	x, y       int
}

// Check if the neuron's bound function needs to be fired
// based on its activation level.
func (n *Neuron) Tick() {
	if n.function == nil {
		return
	}
	if n.activation >= n.threshold {
		n.function()
	}
}

func (n *Neuron) Draw(dc *gg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString(fmt.Sprintf("%0.2f", n.activation), float64(n.x)-15, float64(n.y)-20)
	dc.SetRGB(n.activation, 0, 0)
	dc.DrawCircle(float64(n.x), float64(n.y), float64(10))
	dc.Fill()
}

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
	s.to.activation += s.weight * s.from.activation

	if s.to.activation > 0 {
		s.to.activation = 1
	}
	if s.to.activation < 0 {
		s.to.activation = 0
	}

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

	dc.SetLineWidth(3)
	dc.DrawLine(float64(s.from.x), float64(s.from.y), float64(s.to.x), float64(s.to.y))
	// dc.SetStrokeStyle(gg.NewSolidPattern(color.RGB(0, 0, 0)))
	dc.Stroke()
}

type Brain struct {
	inputNeurons    []Neuron
	internalNeurons []Neuron
	outputNeurons   []Neuron
	synapses        []Synapse
	connectome      Connectome
}

// Draw the brain to the provided buffer as PNG.
func (b *Brain) Draw(o *bytes.Buffer) {
	dc := gg.NewContext(BRAIN_IMG_WIDTH, BRAIN_IMG_HEIGHT)
	for i := 0; i < len(b.inputNeurons); i++ {
		b.inputNeurons[i].Draw(dc)
	}
	for i := 0; i < len(b.internalNeurons); i++ {
		b.internalNeurons[i].Draw(dc)
	}
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].Draw(dc)
	}
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].Draw(dc)
	}
	dc.EncodePNG(o)
}

// Draw the brain to the provided buffer as PNG.
func (b *Brain) SetNeuronPositions() {
	y_spacing := BRAIN_IMG_HEIGHT / 4
	y := y_spacing
	x_spacing := BRAIN_IMG_WIDTH / (len(b.inputNeurons) + 1)
	// x := x_spacing
	for i := 0; i < len(b.inputNeurons); i++ {
		b.inputNeurons[i].x = x_spacing * (i + 1)
		b.inputNeurons[i].y = y
	}
	y += y_spacing
	x_spacing = BRAIN_IMG_WIDTH / (len(b.internalNeurons) + 1)
	for i := 0; i < len(b.internalNeurons); i++ {
		b.internalNeurons[i].x = x_spacing * (i + 1)
		b.internalNeurons[i].y = y
	}
	y += y_spacing
	x_spacing = BRAIN_IMG_WIDTH / (len(b.outputNeurons) + 1)
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].x = x_spacing * (i + 1)
		b.outputNeurons[i].y = y
	}
}

// Handles serving the world as a png to the webserver
func (b *Brain) ImgHandler(response http.ResponseWriter, request *http.Request) {
	buff := bytes.Buffer{}
	b.Draw(&buff)
	response.Header().Set("Content-Length", fmt.Sprint(buff.Len()))
	response.Header().Set("Content-Type", "image/png")
	response.Write(buff.Bytes())
}

type Connectome struct {
	c [32]byte
}

// Randomise this connectome.
func (c *Connectome) Randomise() {
	rand.Read(c.c[:])
}

// Flip a number of bits at random in this connectome.
func (c *Connectome) Mutate(bitsToFlip int) {
	for i := 0; i < bitsToFlip; i++ {
		byteToFlip := rand.Intn(len(c.c))
		bitToFlip := rand.Intn(8)
		if c.c[byteToFlip]&(1<<bitToFlip) != 0 {
			c.c[byteToFlip] ^= 1 << bitToFlip
		} else {
			c.c[byteToFlip] |= 1 << bitToFlip
		}
	}
}

// Copy from the provided connectome into this one.
func (c *Connectome) CopyFrom(t *Connectome) {
	copy(c.c[:], t.c[:])
}

// Tick the brain.
func (b *Brain) Tick(senses []float64) {
	// Process input neurons
	if len(senses) > len(b.inputNeurons) {
		panic("Too many senses")
	}

	// Iterate over the senses and assign them to the input neurons.
	for i, s := range senses {
		b.inputNeurons[i].activation = s
	}

	// Tick all the synapses to propagate the activation levels.
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].Tick()
	}

	// Process output neurons
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].Tick()
	}
}

const INPUT_NEURON_COUNT = 6
const INTERNAL_NEURON_COUNT = 2
const OUTPUT_NEURONS_COUNT = 6
const SYNAPSE_COUNT = 10

func NewBrain(connectome Connectome) *Brain {

	// Do some basic sanity-checking
	if SYNAPSE_COUNT*3 > len(connectome.c) {
		panic("Connectome is too small")
	}
	if INPUT_NEURON_COUNT > 255 {
		panic("Too many input neurons")
	}
	if INTERNAL_NEURON_COUNT > 255 {
		panic("Too many internal neurons")
	}
	if OUTPUT_NEURONS_COUNT > 255 {
		panic("Too many output neurons")
	}

	b := Brain{}
	b.connectome = connectome

	// Create the neurons in the brain.
	b.inputNeurons = make([]Neuron, INPUT_NEURON_COUNT)
	b.internalNeurons = make([]Neuron, INTERNAL_NEURON_COUNT)
	b.outputNeurons = make([]Neuron, OUTPUT_NEURONS_COUNT)

	// Set where the neurons are displayed in the brain image.
	b.SetNeuronPositions()

	// Set a starting threshold of 0.5 for all neurons.
	for i := 0; i < len(b.inputNeurons); i++ {
		b.outputNeurons[i].threshold = 0.5
	}
	b.synapses = make([]Synapse, SYNAPSE_COUNT)

	// Make some helper slices to help with allocating synapses.
	valid_from_neurons := make([]*Neuron, 0)
	valid_to_neurons := make([]*Neuron, 0)

	// We can only wire from input neurons
	for i := 0; i < INPUT_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.inputNeurons[i])
	}
	// Internal neurons can have connections to and from them
	for i := 0; i < INTERNAL_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.internalNeurons[i])
		valid_to_neurons = append(valid_to_neurons, &b.internalNeurons[i])
	}
	// Output neurons can only have connections to them
	for i := 0; i < OUTPUT_NEURONS_COUNT; i++ {
		valid_to_neurons = append(valid_to_neurons, &b.outputNeurons[i])
	}

	// For each synapse, assign a from and to neuron, and a weight, using three
	// bytes from the connectome.
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].from = valid_from_neurons[int(connectome.c[i*3])%len(valid_from_neurons)]
		b.synapses[i].to = valid_to_neurons[int(connectome.c[i*3+1])%len(valid_to_neurons)]
		// Interpret the third byte of the synapse as a signed 8-bit integer
		b.synapses[i].weight = float64(int8(connectome.c[i*3+2])) / 128
	}
	return &b
}
