// Package main provides ...
package main

import (
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/netease/airinput/go-airinput"
)

func ServeWeb(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `<html>
		<head><title>Native-Airinput</title></head>
		<body><h2>Native Airinput</h2>
		<div><img src="/screen.png" height="500px"/></div>
		<a href="/test">pinch test</test>
		</body></html>`)
	})
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		w, h := airinput.ScreenSize()
		fmt.Printf("width: %d, height: %d\n", w, h)

		lx, ly := w/6, 300
		mx, my := w/2, ly
		rx, ry := w/6*5, ly
		airinput.Pinch(lx, ly, mx, my,
			rx, ry, mx, my, 10, time.Second)

		time.Sleep(time.Second * 1)

		airinput.Pinch(mx, my, lx, ly,
			mx, my, rx, ry, 10, time.Second)
		io.WriteString(rw, "pinch run finish")
	})
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			time.Sleep(500 * time.Microsecond)
			os.Exit(0)
		}()
		io.WriteString(w, "Server exit after 0.5s")
	})
	http.HandleFunc("/screen.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		img, _ := airinput.TakeSnapshot()
		png.Encode(w, img)
	})
	http.ListenAndServe(addr, nil)
}
