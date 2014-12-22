package airinput

// #include "screen.h"
// #include "common.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"os/exec"
	"regexp"
	"sync"
)

var screenOnce = sync.Once{}

const DEV_FB0 = "/dev/graphics/fb0"

// TakeSnapshot of android phone (by read /dev/fb0)
// Only ok with few phones, a lot of phone will got blank image.
func TakeSnapshot2() *image.RGBA {
	var pict C.struct_picture
	C.TakeScreenshot(C.CString(DEV_FB0), &pict)
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
	//fmt.Println(width, height, format)
	img = image.NewRGBA(image.Rectangle{image.ZP, image.Point{int(width), int(height)}})
	_, err = bf.Read(img.Pix)
	return
}

// Refrerence code of python-adbviewclient
func ScreenSize() (width int, height int, err error) {
	out, err := exec.Command("dumpsys", "window").Output()
	if err != nil {
		return
	}
	rsRE := regexp.MustCompile(`\s*mRestrictedScreen=\(\d+,\d+\) (?P<w>\d+)x(?P<h>\d+)`)
	matches := rsRE.FindStringSubmatch(string(out))
	if len(matches) == 0 {
		err = errors.New("get shape(width,height) from device error")
		return
	}
	return atoi(matches[1]), atoi(matches[2]), nil
}

// Use ioctl
// Got error some times
func ScreenSize2() (w, h int) {
	screenOnce.Do(func() {
		C.screen_init()
	})
	C.screen_init()
	w = int(C.width())
	h = int(C.height())
	return
}
