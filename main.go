package main

import (
	"embed"
	"image/color"
	"math/rand"
	"net/http"
	"time"
)

// For the record I do not approve of this kind of funky metaprogramming nonsense.
//go:embed static/index.html
//go:embed static/script.js
var content embed.FS

type Sim struct {
	w World
}

func (s *Sim) BrainImgHandler(response http.ResponseWriter, request *http.Request) {
	s.w.crits[0].b.ImgHandler(response, request)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	sim := Sim{}
	// w := World{}
	sim.w.AddRandomCrits(100)
	sim.w.crits[0].c = color.RGBA{255, 0, 0, 255}

	http.HandleFunc("/world", sim.w.ImgHandler)
	http.HandleFunc("/brain", sim.BrainImgHandler)
	http.Handle("/", http.FileServer(http.FS(content)))
	go http.ListenAndServe("192.168.1.50:8082", nil)
	for {
		sim.w.Tick(1)
		time.Sleep(time.Second)
		// sim.w.CullCrits()
		// println("Living: ", len(sim.w.crits))
		// sim.w.RefillCritsWithMutatedConnectomes(100)
	}
}
