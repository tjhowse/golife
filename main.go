package main

import (
	"embed"
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
	w.AddRandomCrits(100)

	http.HandleFunc("/img", w.ImgHandler)
	http.Handle("/", http.FileServer(http.FS(content)))
	go http.ListenAndServe("192.168.1.50:8082", nil)
	for {
		w.Tick(1000)
		time.Sleep(time.Second)
		w.CullCrits()
		println("Living: ", len(w.crits))
		w.RefillCritsWithMutatedConnectomes(100)
	}
}
