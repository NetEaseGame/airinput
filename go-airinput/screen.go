package airinput

// #include "screen.h"
// #include "common.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"log"
	"os/exec"
	"regexp"
	"strconv"
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
	size := w * h * 4 // Assume bytes per pixel is 4 bytes
	pixes := []byte(C.GoStringN(pict.buffer, C.int(size)))
	img := &image.RGBA{pixes, 4 * w, image.Rect(0, 0, w, h)}
	//img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	return img
}

var SCRBUFLEN int

// TakeSnapshot by cmd: /system/bin/screencap
func TakeSnapshot() (img *image.RGBA, err error) {
	var scrbf *bytes.Buffer
	if SCRBUFLEN == 0 {
		scrbf = bytes.NewBuffer(nil)
	} else {
		scrbf = bytes.NewBuffer(make([]byte, 0, SCRBUFLEN))
	}
	cmd := exec.Command("screencap")
	cmd.Stdout = scrbf
	if err = cmd.Run(); err != nil {
		return
	}
	var width, height, format int32
	binary.Read(scrbf, binary.LittleEndian, &width)
	binary.Read(scrbf, binary.LittleEndian, &height)
	SCRBUFLEN = int(width * height * 4)
	err = binary.Read(scrbf, binary.LittleEndian, &format)
	if err != nil {
		return
	}
	//fmt.Println(width, height, format)
	w, h := int(width), int(height)
	img = &image.RGBA{scrbf.Bytes(), 4 * w, image.Rect(0, 0, w, h)}
	return
	//img = image.NewRGBA(image.Rectangle{image.ZP, image.Point{int(width), int(height)}})
	//_, err = bf.Read(img.Pix)
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

const (
	ROTATION_0 = iota
	ROTATION_90
	ROTATION_180
	ROTATION_270
	ROTATION_UNKNOWN
)

// return 0,1,2,3
func Rotation() int {
	re := regexp.MustCompile(`SurfaceOrientation:\s+(\d+)`)
	dump, err := exec.Command("dumpsys", "input").Output()
	if err != nil {
		log.Println(err)
		return ROTATION_UNKNOWN
	}
	res := re.FindStringSubmatch(string(dump))
	if len(res) != 2 {
		return ROTATION_UNKNOWN
	}
	rotation, _ := strconv.Atoi(res[1])
	return rotation
}

// About: 0.03s
// Rotate coordinate to rotation_0
func CoordRotate(x, y int) (nx, ny int) {
	w, h, er := ScreenSize()
	if w > h {
		w, h = h, w
	}
	if er != nil {
		log.Println("screen size get failed:", er)
		return x, y
	}
	switch Rotation() {
	case ROTATION_90:
		return w - y, x
	case ROTATION_180:
		return w - x, h - y
	case ROTATION_270:
		return y, h - x
	case ROTATION_0:
		fallthrough
	default:
		return x, y
	}
}

func GetRawSize(event string) (width, height int, err error) {
	mxptn := regexp.MustCompile(`0035.*max (\d+)`)
	myptn := regexp.MustCompile(`0036.*max (\d+)`)
	out, err := exec.Command("getevent", "-p", event).Output()
	if err != nil {
		return
	}
	err = errors.New("touchpad event not recognized")
	mxs := mxptn.FindStringSubmatch(string(out))
	if len(mxs) == 0 {
		return
	}
	mys := myptn.FindStringSubmatch(string(out))
	if len(mys) == 0 {
		return
	}
	return atoi(mxs[1]), atoi(mys[1]), nil
}
