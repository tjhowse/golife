package main

import (
	"bytes"
	"fmt"
	"net/http"

	gg "github.com/fogleman/gg"
)

const BRAIN_IMG_WIDTH = 500
const BRAIN_IMG_HEIGHT = 500

const SIGMOID_OUTPUT_ACTIVATION_THRESHOLD = 0.75

const INPUT_NEURON_COUNT = 6
const INTERNAL_NEURON_COUNT = 5
const OUTPUT_NEURONS_COUNT = 6
const SYNAPSE_COUNT = 10

type Brain struct {
	inputNeurons    []Neuron
	internalNeurons []Neuron
	outputNeurons   []Neuron
	synapses        []Synapse
	connectome      Connectome
	dc              *gg.Context
}

// Draw the brain to the provided buffer as PNG.
func (b *Brain) Draw(o *bytes.Buffer) {
	// dc := gg.NewContext(BRAIN_IMG_WIDTH, BRAIN_IMG_HEIGHT)
	b.dc.SetRGB(1, 1, 1)
	b.dc.Clear()
	b.dc.SetRGB(0, 0, 0)
	b.dc.DrawString("In:", 0, BRAIN_IMG_HEIGHT/4)
	b.dc.DrawString("Internal:", 0, 2*BRAIN_IMG_HEIGHT/4)
	b.dc.DrawString("Out:", 0, 3*BRAIN_IMG_HEIGHT/4)
	for i := 0; i < len(b.inputNeurons); i++ {
		b.inputNeurons[i].Draw(b.dc)
	}
	for i := 0; i < len(b.internalNeurons); i++ {
		b.internalNeurons[i].Draw(b.dc)
	}
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].Draw(b.dc)
	}
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].Draw(b.dc)
	}
	b.dc.EncodePNG(o)
}

// Set the positions of the neurons in the drawn image.
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

	// Tick all the synapses to propagate the input sums to prepare for ticking
	// the neurons.
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].Tick()
	}

	// Process internal neurons
	for i := 0; i < len(b.internalNeurons); i++ {
		b.internalNeurons[i].Tick()
	}

	// Process output neurons
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].Tick()
	}
}

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
	b.dc = gg.NewContext(BRAIN_IMG_WIDTH, BRAIN_IMG_HEIGHT)

	// Create the neurons in the brain.
	b.inputNeurons = make([]Neuron, INPUT_NEURON_COUNT)
	b.internalNeurons = make([]Neuron, INTERNAL_NEURON_COUNT)
	b.outputNeurons = make([]Neuron, OUTPUT_NEURONS_COUNT)

	// Set where the neurons are displayed in the brain image.
	b.SetNeuronPositions()

	b.synapses = make([]Synapse, SYNAPSE_COUNT)

	// Make some helper slices to help with allocating synapses.
	valid_from_neurons := make([]*Neuron, 0)
	valid_to_neurons := make([]*Neuron, 0)
	all_neurons := make([]*Neuron, 0)

	// We can only wire from input neurons
	for i := 0; i < INPUT_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.inputNeurons[i])
		all_neurons = append(all_neurons, &b.inputNeurons[i])
	}
	// Internal neurons can have connections to and from them
	for i := 0; i < INTERNAL_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.internalNeurons[i])
		valid_to_neurons = append(valid_to_neurons, &b.internalNeurons[i])
		all_neurons = append(all_neurons, &b.internalNeurons[i])
	}
	// Output neurons can only have connections to them
	for i := 0; i < OUTPUT_NEURONS_COUNT; i++ {
		valid_to_neurons = append(valid_to_neurons, &b.outputNeurons[i])
		all_neurons = append(all_neurons, &b.outputNeurons[i])
	}
	for i := 0; i < len(all_neurons); i++ {
		all_neurons[i].bias = float64(int8(b.connectome.GetByte())) / 128
	}

	// Create the synapses
	// For each synapse, assign a from and to neuron, and a weight, using three
	// bytes from the connectome.
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].from = valid_from_neurons[int(b.connectome.GetByte())%len(valid_from_neurons)]
		b.synapses[i].to = valid_to_neurons[int(b.connectome.GetByte())%len(valid_to_neurons)]
		// Interpret the third byte of the synapse as a signed 8-bit integer
		b.synapses[i].weight = float64(int8(b.connectome.GetByte())) / 128
	}
	return &b
}
