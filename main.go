// Currently this is a sample
package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/netease/airinput/go-airinput"
)

var (
//	m = macaron.Classic()
)

func ServeWeb(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		img, _ := airinput.TakeSnapshot()
		png.Encode(w, img)
	})
	http.ListenAndServe(addr, nil)
}

func main() {
	ServeWeb(":21000")
	airinput.Debug(false)
	img, err := airinput.TakeSnapshot()
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Create("/data/local/tmp/air.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	png.Encode(fd, img)
	return

	w, h := airinput.ScreenSize()
	fmt.Printf("width: %d, height: %d\n", w, h)

	lx, ly := w/6, 300
	mx, my := w/2, ly
	rx, ry := w/6*5, ly

	// initial
	if err := airinput.Init(); err != nil {
		log.Fatal(err)
	}
	airinput.Pinch(lx, ly, mx, my,
		rx, ry, mx, my, 10, time.Second)

	time.Sleep(time.Second * 1)

	airinput.Pinch(mx, my, lx, ly,
		mx, my, rx, ry, 10, time.Second)
}
