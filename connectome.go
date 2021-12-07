package main

import "math/rand"

// Think of this as a blob of entropy used to hook up a brain.
type Connectome struct {
	c [64]byte
	i uint8
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

func (c *Connectome) GetByte() byte {
	b := c.c[c.i]
	c.i++
	if c.i >= uint8(len(c.c)) {
		panic("Ran out of bytes in the connectome")
	}
	return b
}
