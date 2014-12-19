// Package main provides ...
package main

import (
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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
						setTimeout(function(){location.reload(true)}, 2000);
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
	screenFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		img, _ := airinput.TakeSnapshot()
		png.Encode(w, img)
	}
	http.HandleFunc("/screen.png", screenFunc)
	http.HandleFunc("/screen/", screenFunc)
	http.ListenAndServe(addr, nil)
}
