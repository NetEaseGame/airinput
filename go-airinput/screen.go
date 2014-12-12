package airinput

// #include "screen.h"
// #include "common.h"
import "C"
import (
	"image"
	"sync"
)

var screenOnce = sync.Once{}

func ScreenSize() (w, h int) {
	screenOnce.Do(func() {
		C.screen_init()
	})
	w = int(C.width())
	h = int(C.height())
	return
}

// TakeSnapshot of android phone
// Only ok with few phones, a lot of phone will got blank image.
func TakeSnapshot() *image.RGBA {
	var pict C.struct_picture
	C.TakeScreenshot(C.CString("/dev/graphics/fb0"), &pict)
	w, h := int(pict.xres), int(pict.yres)
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	size := w * h * 4 // Assume bytes per pixel is 4 bytes

	img.Pix = []byte(C.GoStringN(pict.buffer, C.int(size)))
	return img
}
