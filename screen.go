package main

// #include "screen.h"
import "C"
import "sync"

var screenOnce = sync.Once{}

func ScreenSize() (w, h int) {
	screenOnce.Do(func() {
		C.screen_init()
	})
	w = int(C.width())
	h = int(C.height())
	return
}
