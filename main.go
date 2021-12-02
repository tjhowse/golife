package main

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"net/http"
	"time"

	gg "github.com/fogleman/gg"
)

const WORLD_HEIGHT = 1000
const WORLD_WIDTH = 1000

// TODO work out why this doesn't work.
type ticker interface {
	Tick()
}

func TickAll(t []ticker) {
	for i := 0; i < len(t); i++ {
		t[i].Tick()
	}
}

func TickIt(t *ticker) {
	(*t).Tick()
}

// TickAll(w.crits)

type World struct {
	crits  []*Crit
	frames int // The number of frames that have passed
}

func (w *World) AddCrit(c *Crit) {
	w.crits = append(w.crits, c)
}

func (w *World) Tick() {
	for i := 0; i < len(w.crits); i++ {
		w.crits[i].Tick()
	}
	w.frames++
}

func (w *World) Draw(o *bytes.Buffer) {
	dc := gg.NewContext(WORLD_WIDTH, WORLD_HEIGHT)

	for i := 0; i < len(w.crits); i++ {
		w.crits[i].Draw(dc)
	}
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.EncodePNG(o)
}

func (w *World) ImgHandler(response http.ResponseWriter, request *http.Request) {

	buff := bytes.Buffer{}
	w.Draw(&buff)
	response.Header().Set("Content-Length", fmt.Sprint(buff.Len()))
	response.Header().Set("Content-Type", "image/png")
	response.Write(buff.Bytes())
}

// For the record I do not approve of this kind of funky metaprogramming nonsense.
//go:embed static/index.html
//go:embed static/script.js
var content embed.FS

func main() {
	w := World{}
	w.AddCrit(NewCrit(WORLD_HEIGHT/2, WORLD_WIDTH/2, 10, color.Black, &w))
	// w.AddCrit(NewCrit(WORLD_HEIGHT/3, WORLD_WIDTH/3, 10, color.Black, &w))
	http.HandleFunc("/img", w.ImgHandler)
	http.Handle("/", http.FileServer(http.FS(content)))
	go http.ListenAndServe("192.168.1.50:8082", nil)
	for {
		w.Tick()
		time.Sleep(time.Second)
		println("Tick")
	}
}
