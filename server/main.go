// Package main provides ...
package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/snappy-go/snappy"
	"github.com/codeskyblue/pngdiff"
	"golang.org/x/image/bmp"
)

func main() {
	var md5sum string
	var rip = "10.242.87.153"
	var lastImg *image.RGBA
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		serialNo := r.URL.Query().Get("serialno")
		fmt.Println(strings.Split(r.RemoteAddr, ":")[0], serialNo)
		rip = strings.Split(r.RemoteAddr, ":")[0]

		io.WriteString(w, "connected")
	})
	http.HandleFunc("/screen.png", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(fmt.Sprintf("http://%s:21000/patch.snappy?md5sum="+md5sum, rip))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotModified {
			png.Encode(w, lastImg)
			return
		}
		data, _ := ioutil.ReadAll(resp.Body)
		rawPixes, _ := snappy.Decode(nil, data)
		_ = rawPixes
		fmt.Println(resp.Header)

		hdr := resp.Header
		md5sum = hdr.Get("X-Md5sum")
		var isPatch = hdr.Get("X-Patch") == "true"
		var width, height int
		fmt.Sscanf(hdr.Get("X-Width")+" "+hdr.Get("X-Height"), "%d %d", &width, &height)

		img := image.NewRGBA(image.Rectangle{image.ZP, image.Point{width, height}})
		img.Pix = rawPixes
		if isPatch {
			lastImg, _ = pngdiff.Patch(lastImg, img)
			start := time.Now()
			bmp.Encode(w, lastImg)
			//png.Encode(w, lastImg)
			fmt.Println("patch:", time.Now().Sub(start))
			return
		} else {
			lastImg = img
			//png.Encode(w, img)
			bmp.Encode(w, lastImg)
		}
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}
