package airinput

// #include "screen.h"
// #include "common.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"os/exec"
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

// TakeSnapshot of android phone (by read /dev/fb0)
// Only ok with few phones, a lot of phone will got blank image.
func TakeSnapshot2() *image.RGBA {
	var pict C.struct_picture
	C.TakeScreenshot(C.CString("/dev/graphics/fb0"), &pict)
	w, h := int(pict.xres), int(pict.yres)
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	size := w * h * 4 // Assume bytes per pixel is 4 bytes

	img.Pix = []byte(C.GoStringN(pict.buffer, C.int(size)))
	return img
}

// TakeSnapshot by cmd: /system/bin/screencap
func TakeSnapshot() (img *image.RGBA, err error) {
	bf := bytes.NewBuffer(nil)
	cmd := exec.Command("/system/bin/screencap")
	cmd.Stdout = bf
	if err = cmd.Run(); err != nil {
		return
	}
	var width, height, format int32
	binary.Read(bf, binary.LittleEndian, &width)
	binary.Read(bf, binary.LittleEndian, &height)
	err = binary.Read(bf, binary.LittleEndian, &format)
	if err != nil {
		return
	}
	fmt.Println(width, height, format)
	img = image.NewRGBA(image.Rectangle{image.ZP, image.Point{int(width), int(height)}})
	_, err = bf.Read(img.Pix)
	return
}
