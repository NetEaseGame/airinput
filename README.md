airinput
=====================
Simulate touch,drag,pinch for Android phone.

### The lib of go-airinput
[![GoDoc](https://godoc.org/github.com/NetEase/airinput/go-airinput?status.svg)](https://godoc.org/github.com/NetEase/airinput/go-airinput)

### Distribution
I put a pre compiled file in **dist** folder

About the usage:

1. use `run.bat` to push file to your phone.
2. open browser `http://<phone ip addr>:21000`

![IMG](images/browser-airinput.png)

Support post js code to `http://<phone ip addr>:21000/runjs`

Example js code:

	exec("input", "keyevent", "ENTER")
	tap(20, 30, 100) // position: (20, 30), duration: 100ms
	drag(10, 12, 50, 60, 10, 100) // start(10, 12), end(50, 60), steps: 10, duration: 100ms
	// pinch(ax0, ay0, ax1, ay1, bx0, by0, bx1, by1, steps, duration)

Also support use adb to run js

	adb shell /data/local/tmp/air-native -i /dev/input/event1 -runjs='tap(400, 400, 2000)'

### Snapshot performance
* PNG has best quality and file is small, but compress use lot of time.
* BMP file is big, but donnot use compress time.
* JPG quality no good than PNG, but fast also very small.

Use `screencap -p` is very very slow, about 4~5s. but `screencap` is fast, only take about 500ms.

I develop the pngdiff lib <https://github.com/codeskyblue/pngdiff>. Diff is fast, only about 70ms. But encoding still take a lot of time(about 400ms). `^_|`

So the total time is about 1s.

### About
Still in develop, but the code is healthy. 

This code need Go1.4, follow [offical instruction](http://code.google.com/p/go/source/browse/README?repo=mobile) to setup environment.

A lot code is from <https://github.com/wlach/orangutan>, orangutan you are a great people.

use `sh build.sh` to build.

Sample code is in `main.go` now. 

Licence is under [MIT](LICENSE).
