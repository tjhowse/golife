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

func NewBrain(connectome [32]byte) *Brain {
	b := Brain{}
	b.inputNeurons = make([]Neuron, 4)
	b.internalNeurons = make([]Neuron, 2)
	b.outputNeurons = make([]Neuron, 4)
	for i := 0; i < len(b.inputNeurons); i++ {
		b.outputNeurons[i].threshold = 0.5
	}
	b.synapses = make([]Synapse, 10)
	b.synapses[0].from = &b.inputNeurons[0]
	b.synapses[0].to = &b.outputNeurons[0]
	b.synapses[0].weight = 1
	// TODO map connectome to synapses
	// for _, c := range connectome {
	// }
	return &b
}
