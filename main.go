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

func main() {
	rand.Seed(time.Now().UnixNano())
	w := World{}
	x := 0
	y := 0
	for i := 0; i < 100; i++ {
		for x, y = -1, -1; !w.CheckMove(x, y); {
			x = rand.Intn(WORLD_WIDTH/GRID_TO_PIXEL) * GRID_TO_PIXEL
			y = rand.Intn(WORLD_WIDTH/GRID_TO_PIXEL) * GRID_TO_PIXEL
		}
		connectome := [32]byte{}
		rand.Read(connectome[:])
		w.AddCrit(NewCrit(x, y, GRID_TO_PIXEL, color.Black, &w, connectome))
	}

	// w.AddCrit(NewCrit(WORLD_HEIGHT/3, WORLD_WIDTH/3, 10, color.Black, &w))
	http.HandleFunc("/img", w.ImgHandler)
	http.Handle("/", http.FileServer(http.FS(content)))
	go http.ListenAndServe("192.168.1.50:8082", nil)
	for {
		w.Tick()
		time.Sleep(time.Second)
	}
}
