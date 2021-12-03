package main

type ActivationFunction func()

type Neuron struct {
	activation float64
	threshold  float64
	function   ActivationFunction
}

func (n *Neuron) Tick() {
	if n.function == nil {
		return
	}
	if n.activation >= n.threshold {
		n.function()
	}
}

type Synapse struct {
	weight   float64
	to, from *Neuron
}

type Brain struct {
	inputNeurons    []Neuron
	internalNeurons []Neuron
	outputNeurons   []Neuron
	synapses        []Synapse
}

func (b *Brain) Tick(senses []float64) {
	// Process input neurons
	if len(senses) > len(b.inputNeurons) {
		panic("Too many senses")
	}
	for i, s := range senses {
		b.inputNeurons[i].activation = s
	}
	for i := 0; i < len(b.synapses); i++ {
		if b.synapses[i].from == nil || b.synapses[i].to == nil {
			continue
		}
		b.synapses[i].to.activation += b.synapses[i].from.activation * b.synapses[i].weight
	}
	// Process output neurons
	for i := 0; i < len(b.outputNeurons); i++ {
		b.outputNeurons[i].Tick()
	}
}

// func (b *Brain)

// const INPUT_NEURON_BITS = 2
// const INTERNAL_NEURON_BITS = 1
// const OUTPUT_NEURON_BITS = 2
// const NEURON_WEIGHT_BITS = 8 - INPUT_NEURON_BITS - INTERNAL_NEURON_BITS - OUTPUT_NEURON_BITS

const INPUT_NEURON_COUNT = 6
const INTERNAL_NEURON_COUNT = 2
const OUTPUT_NEURONS_COUNT = 6
const SYNAPSE_COUNT = 10

func NewBrain(connectome [32]byte) *Brain {

	// Do some basic sanity-checking
	if SYNAPSE_COUNT*3 > len(connectome) {
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
	b.inputNeurons = make([]Neuron, INPUT_NEURON_COUNT)
	b.internalNeurons = make([]Neuron, INTERNAL_NEURON_COUNT)
	b.outputNeurons = make([]Neuron, OUTPUT_NEURONS_COUNT)
	for i := 0; i < len(b.inputNeurons); i++ {
		b.outputNeurons[i].threshold = 0.5
	}
	b.synapses = make([]Synapse, SYNAPSE_COUNT)
	// b.synapses[0].from = &b.inputNeurons[0]
	// b.synapses[0].to = &b.outputNeurons[0]
	// b.synapses[0].weight = 1
	valid_from_neurons := make([]*Neuron, 0)
	valid_to_neurons := make([]*Neuron, 0)
	// We can only wire from input neurons
	for i := 0; i < INPUT_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.inputNeurons[i])
	}
	// Internal neurons can have connections to and from them
	for i := 0; i < INTERNAL_NEURON_COUNT; i++ {
		valid_from_neurons = append(valid_from_neurons, &b.internalNeurons[i])
		valid_to_neurons = append(valid_to_neurons, &b.inputNeurons[i])
	}
	// Output neurons can only have connections to them
	for i := 0; i < OUTPUT_NEURONS_COUNT; i++ {
		valid_to_neurons = append(valid_to_neurons, &b.outputNeurons[i])
	}
	for i := 0; i < len(b.synapses); i++ {
		b.synapses[i].from = valid_from_neurons[int(connectome[i*3])%len(valid_from_neurons)]
		b.synapses[i].to = valid_to_neurons[int(connectome[i*3+1])%len(valid_to_neurons)]
		b.synapses[i].weight = float64(connectome[i*3+2]) / float64(255)
	}
	return &b
}
