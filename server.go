// Package main provides ...
package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
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

func init() {
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
			$screen.attr("src", "/screen.jpg?random="+new Date().getTime());

			var loadscreen = function(e){
				$screen.attr("src", "/screen.jpg?random="+new Date().getTime());
			};

			$screen
				.unbind()
				.load(function(){
					console.log("image loaded");
					loadscreen();
				});

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
				msec = 200;
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
	// Run Javascript code
	http.HandleFunc("/runjs", func(w http.ResponseWriter, r *http.Request) {
		code, _ := ioutil.ReadAll(r.Body)
		ret, _ := RunJS(string(code))
		io.WriteString(w, ret.String())
	})

	// Exit the server after 0.5s
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			time.Sleep(500 * time.Microsecond)
			os.Exit(0)
		}()
		io.WriteString(w, "Server exit after 0.5s")
	})

	http.HandleFunc("/screen.png", func(w http.ResponseWriter, r *http.Request) {
		img, _ := airinput.Snapshot()
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // for ajax
		if ckname := r.URL.Query().Get("cookiename"); ckname != "" {
			w.Header().Set("Set-Cookie", fmt.Sprintf("%s=%d", ckname, time.Now().UnixNano()/1000))
		}
		png.Encode(w, img)
	})
	var JPEG_QUALITY = &jpeg.Options{60}
	http.HandleFunc("/screen.jpg", func(w http.ResponseWriter, r *http.Request) {
		img, _ := airinput.Snapshot()
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // for ajax
		if ckname := r.URL.Query().Get("cookiename"); ckname != "" {
			w.Header().Set("Set-Cookie", fmt.Sprintf("%s=%d", ckname, time.Now().UnixNano()/1000))
		}
		jpeg.Encode(w, img, JPEG_QUALITY)
	})
	// patch part
	var lastImage *image.RGBA
	http.HandleFunc("/patch.jpg", func(w http.ResponseWriter, r *http.Request) {
		img, _ := airinput.Snapshot()
		if lastImage == nil {
			jpeg.Encode(w, img, JPEG_QUALITY)
		} else {
			patch, _ := pngdiff.Diff(lastImage, img)
			jpeg.Encode(w, patch, JPEG_QUALITY)
		}
		lastImage = img
	})

	// Send full image first time
	// Then send patch until connection closed
	http.HandleFunc("/patch.hijack", func(w http.ResponseWriter, r *http.Request) {
		log.Println("new conn")
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

		defer func() {
			if e := recover(); e != nil {
				log.Println("panic", e)
			}
		}()
		bwrite := func(wr io.Writer, a interface{}) {
			if err := binary.Write(wr, binary.LittleEndian, a); err != nil {
				panic(err.Error())
			}
		}
		cimg, _ := airinput.Snapshot()
		var first bool = true
		for {
			var bytes []byte
			if !first {
				img, _ := airinput.Snapshot()
				patch, _ := pngdiff.Diff(cimg, img)
				bytes, _ = snappy.Encode(nil, patch.Pix)
				cimg = img
				log.Println("IN-LEN", len(bytes))
			} else {
				first = false
				// Write width and height
				bwrite(bufrw, uint32(cimg.Rect.Max.X))
				bwrite(bufrw, uint32(cimg.Rect.Max.Y))
				bytes, _ = snappy.Encode(nil, cimg.Pix)
				log.Println("LEN", len(bytes))
			}
			err := binary.Write(bufrw, binary.LittleEndian, uint32(len(bytes)))
			if err != nil {
				break
			}
			if _, err := bufrw.Write(bytes); err != nil {
				break
			}
			bufrw.Flush()
			time.Sleep(time.Millisecond * 100)
		}
		log.Println("END")
	})
	// Request need md5sum in query
	// eg: http GET /patch.snappy?md5sum=xklj21901294123912
	// Will be a patch file when header['X-Patch'] == true
	cache := NewRGBACache(2) // cache size = 2
	http.HandleFunc("/patch.snappy", func(w http.ResponseWriter, r *http.Request) {
		img, _ := airinput.Snapshot()
		md5old := r.FormValue("md5sum")
		md5new := fmt.Sprintf("%x", md5.Sum(img.Pix))
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
			patch, _ := pngdiff.Diff(imold, img)
			bytes, err := snappy.Encode(nil, patch.Pix)
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
			//png.Encode(w, patch)
		} else {
			w.Header().Set("X-Patch", "false")
			bytes, err := snappy.Encode(nil, img.Pix)
			if err != nil {
				log.Println(err)
			}
			w.Write(bytes)
			//png.Encode(w, img)
		}
	})
}

func ServeWeb(addr string) {
	http.ListenAndServe(addr, nil)
}
