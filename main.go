package main

import (
	"embed"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	gg "github.com/fogleman/gg"
	"github.com/tjhowse/tjgo"
)

// For the record I do not approve of this kind of funky metaprogramming nonsense.
//go:embed static/index.html
//go:embed static/script.js
var content embed.FS

type Sim struct {
	w World
}

const CRIT_COUNT = 100

func main() {
	rand.Seed(time.Now().UnixNano())
	sim := Sim{}
	displayedBrain := 0
	// w := World{}
	sim.w.dc = gg.NewContext(WORLD_WIDTH*GRID_TO_PIXEL, WORLD_HEIGHT*GRID_TO_PIXEL)
	sim.w.AddRandomCrits(CRIT_COUNT)

	http.HandleFunc("/world", sim.w.ImgHandler)
	http.HandleFunc("/brain", func(response http.ResponseWriter, request *http.Request) {
		// return
		if displayedBrain >= len(sim.w.crits) {
			return
		}
		sim.w.crits[displayedBrain].b.ImgHandler(response, request)
	})
	http.HandleFunc("/restart", func(response http.ResponseWriter, request *http.Request) {
		sim.w.CullCrits(0, WORLD_WIDTH)
		sim.w.RefillCritsWithMutatedConnectomes(CRIT_COUNT)
	})
	http.HandleFunc("/click", func(response http.ResponseWriter, request *http.Request) {
		u, err := url.Parse(request.URL.String())
		if err != nil {
			fmt.Errorf("Error parsing URL: %v", err)
			return
		}
		m, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			fmt.Errorf("Error parsing query: %v", err)
			return
		}
		if m["id"][0] == "world" {
			i, err := sim.w.GetCritIndexClosestToImgClick(tjgo.Str2int(m["x"][0]), tjgo.Str2int(m["y"][0]))
			if err != nil {
				fmt.Errorf("Error getting crit index: %v", err)
				return
			}
			displayedBrain = i
		}
		// fmt.Printf(("%v\n"), m)
		// println(x, y)

	})
	http.HandleFunc("/tick", func(response http.ResponseWriter, request *http.Request) {
		u, err := url.Parse(request.URL.String())
		if err != nil {
			fmt.Errorf("Error parsing URL: %v", err)
			return
		}
		m, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			fmt.Errorf("Error parsing query: %v", err)
			return
		}
		count := tjgo.Str2int(m["ticks"][0])
		fmt.Println("Ticking ", count, " ticks")
		sim.w.Tick(count)

	})
	http.HandleFunc("/cull", func(response http.ResponseWriter, request *http.Request) {
		sim.w.CullCrits(CRIT_COUNT/WORLD_HEIGHT, WORLD_WIDTH)
		displayedBrain = 0
	})
	http.HandleFunc("/topUp", func(response http.ResponseWriter, request *http.Request) {
		sim.w.RandomiseCritPositions()
		sim.w.RefillCritsWithMutatedConnectomes(CRIT_COUNT)
	})
	http.HandleFunc("/cycles", func(response http.ResponseWriter, request *http.Request) {
		displayedBrain = 0
		for i := 0; i < 10; i++ {
			sim.w.Tick(1000)
			sim.w.CullCrits(CRIT_COUNT/WORLD_HEIGHT, WORLD_WIDTH)
			sim.w.RandomiseCritPositions()
			sim.w.RefillCritsWithMutatedConnectomes(CRIT_COUNT)
		}
	})
	http.Handle("/", http.FileServer(http.FS(content)))
	go http.ListenAndServe("192.168.1.50:8082", nil)
	for {
		sim.w.Tick(1)
		time.Sleep(time.Second / 2)
		// sim.w.CullCrits(CRIT_COUNT/WORLD_HEIGHT, WORLD_WIDTH)
		// println("Living: ", len(sim.w.crits))
		// sim.w.RandomiseCritPositions()
		// sim.w.RefillCritsWithMutatedConnectomes(CRIT_COUNT)
	}
}
