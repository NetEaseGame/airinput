// Package main provides ...
package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"code.google.com/p/snappy-go/snappy"

	"github.com/codeskyblue/pngdiff"
	"github.com/netease/airinput/go-airinput"
)

func ServeWeb(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `<html>
		<head>
			<meta charset="UTF-8">
			<title>Native-Airinput</title>
		</head>
		<body>
			<h2>Native Airinput 
				<small><button id="refresh">刷新</button></small> <input type="checkbox" id="autoreload" checked />
			</h2> 
			<div><span id="message"></span></div>
		<div><img id="screen" src="/screen.png" height="500px"/></div>
		<textarea id="jscode" style="height:100px; width:500px"></textarea>
		<button id="btn-run">RUN</button>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.1/jquery.js"></script>
		<script>
		$(function(){
			/* common functions */
			String.prototype.format = function() {
				var formatted = this;
				for (var arg in arguments) {
					formatted = formatted.replace("{" + arg + "}", arguments[arg]);
				}
				return formatted;
			};
			$("#btn-run").click(function(){
				$.ajax('/runjs', {type:'POST', processData: false, data: $("#jscode").val()});
			});
			var startTime = 0, startPoint = [0, 0];
			var steps = 0;
			var $screen = $("#screen");
			var loadscreen = function(e){
				// 
				$screen.attr("src", "/screen/"+new Date().getTime());
				//location.reload(true);
			};
			//setInterval(function() { loadscreen() }, 4000);
			$("#refresh").click(function(){loadscreen()});

			/* click event */
			var mousepoint = function(e){
				var offset = $screen.offset();
				var scale = $screen[0].naturalHeight*1.0/$screen.height();
				var x = parseInt((e.pageX - offset.left)*scale);
				var y = parseInt((e.pageY - offset.top) *scale);
				return [x, y];
			};
			$("#screen").mousedown(function(e){
				startPoint = mousepoint(e);
				startTime = new Date().getMilliseconds();
				steps = 0;
				console.log(startPoint);
			});
			$("#screen").mousemove(function(e){
				steps += 1;
			});
			$("#screen").mouseup(function(e){
				var ep = mousepoint(e);
				var msec = new Date().getMilliseconds() - startTime;
				var jscode = 'tap({0}, {1}, {2})'.format(ep[0], ep[1], msec);
				console.log(jscode);
				$.ajax('/runjs', {type:'POST', processData: false, data: jscode, success: function(){
					if ($('#autoreload')[0].checked){
						setTimeout(function(){loadscreen()}, 2000);
					}
				}});
				console.log(startPoint, ep, steps, msec);
			});
		});
		</script>
		</body></html>`)
	})
	http.HandleFunc("/runjs", func(w http.ResponseWriter, r *http.Request) {
		code, _ := ioutil.ReadAll(r.Body)
		ret, _ := RunJS(string(code))
		io.WriteString(w, ret.String())
	})
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		w, h, _ := airinput.ScreenSize()
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
	cache := NewRGBACache(2) // cache size = 2
	screenFunc := func(w http.ResponseWriter, r *http.Request) {
		img, _ := airinput.TakeSnapshot()
		md5sum := fmt.Sprintf("%x", md5.Sum(img.Pix))
		cache.Put(md5sum, img)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("X-Md5sum", md5sum)
		png.Encode(w, img)
	}
	http.HandleFunc("/screen.png", screenFunc)
	http.HandleFunc("/screen/", screenFunc)
	http.HandleFunc("/patch.hijack", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("new conn")
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		cimg, _ := airinput.TakeSnapshot()
		var first bool = true
		for {
			time.Sleep(time.Millisecond * 1000)
			//_, er := bufrw.WriteString("Now we're speaking raw TCP. Say hi: \n")
			fmt.Println(first)
			if !first {
				img, _ := airinput.TakeSnapshot()
				patch, _ := pngdiff.Diff(cimg, img)
				bytes, _ := snappy.Encode(nil, patch.Pix)
				cimg = img
				fmt.Println("IN-LEN", len(bytes))
				binary.Write(bufrw, binary.LittleEndian, uint32(len(bytes)))
				_, er := bufrw.Write(bytes)
				if er != nil {
					break
				}
				bufrw.Flush()
				continue
			}
			first = false

			bytes, _ := snappy.Encode(nil, cimg.Pix)
			fmt.Println("LEN", len(bytes))
			binary.Write(bufrw, binary.LittleEndian, uint32(len(bytes)))
			_, er := bufrw.Write(bytes)
			if er != nil {
				break
			}
			bufrw.Flush()
		}
		fmt.Println("END")
	})
	http.HandleFunc("/patch.snappy", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		img, _ := airinput.TakeSnapshot()
		fmt.Println("snapshot", time.Now().Sub(start))
		md5old := r.FormValue("md5sum")
		fmt.Println(md5old)
		//fmt.Println(cache.Get(md5old))
		start = time.Now()
		md5new := fmt.Sprintf("%x", md5.Sum(img.Pix))
		fmt.Println("md5", time.Now().Sub(start))
		cache.Put(md5new, img)
		if md5new == md5old {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("X-Md5sum", md5new)
		w.Header().Set("X-Width", strconv.Itoa(img.Bounds().Max.X))
		w.Header().Set("X-Height", strconv.Itoa(img.Bounds().Max.Y))

		if imold, exists := cache.Get(md5old); exists {
			w.Header().Set("X-Patch", "true")
			start = time.Now()
			patch, _ := pngdiff.Diff(imold, img)
			fmt.Println("diff", time.Now().Sub(start))
			start = time.Now()
			bytes, err := snappy.Encode(nil, patch.Pix)
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
			//png.Encode(w, patch)
			fmt.Println("patch encode", time.Now().Sub(start))
		} else {
			w.Header().Set("X-Patch", "false")
			start = time.Now()
			bytes, err := snappy.Encode(nil, img.Pix)
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
			//png.Encode(w, img)
			fmt.Println("raw encode", time.Now().Sub(start))
		}
	})
	http.ListenAndServe(addr, nil)
}
