// Currently this is a sample
package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/netease/airinput/go-airinput"
)

func main() {
	airinput.Debug(false)
	img := airinput.TakeSnapshot()
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
