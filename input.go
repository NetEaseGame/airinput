/*
click down

sendevent /dev/input/event1 1 330 1
sendevent /dev/input/event1 3 53 539
sendevent /dev/input/event1 3 54 959
sendevent /dev/input/event1 0 0 0

click up

sendevent /dev/input/event1 3 57 243
sendevent /dev/input/event1 1 330 0
sendevent /dev/input/event1 0 0 0
*/

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type DeviceTouchScreen struct {
	InputEvent string `json:"input_event"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	RawWidth   int    `json:"raw_width"`
	RawHeight  int    `json:"raw_height"`
}

type InputDevices struct {
	TouchScreen DeviceTouchScreen `json:"touchscreen"`
}

type Event struct {
	_time uint64
	Type  uint16
	Code  uint16
	Value int32
}

var (
	devicefd *os.File
	iptdevs  *InputDevices
)

func initInput() (err error) {
	iptdevs, err = getinputevent()
	if err != nil {
		return
	}
	tscreen := iptdevs.TouchScreen
	devicefd, err = os.OpenFile(tscreen.InputEvent, os.O_RDWR, 0644)
	return err
}

func deferInput() {
	if devicefd != nil {
		devicefd.Close()
	}
}

func sendevent(fd io.Writer, type_, code string, value int32) (err error) {
	event := new(Event)

	fmt.Sscanf(type_, "%x", &event.Type)
	fmt.Sscanf(code, "%x", &event.Code)
	event.Value = value

	buffer := bytes.NewBuffer(nil)
	binary.Write(buffer, binary.LittleEndian, event._time)
	binary.Write(buffer, binary.LittleEndian, event.Type)
	binary.Write(buffer, binary.LittleEndian, event.Code)
	binary.Write(buffer, binary.LittleEndian, event.Value)
	_, err = io.Copy(fd, buffer)
	return err
}

func getinputevent() (inputdevs *InputDevices, err error) {
	var curpwd = filepath.Dir(os.Args[0])
	jsonfile := filepath.Join(curpwd, "devices.json")
	fd, er := os.Open(jsonfile)
	if er != nil {
		err = er
		return
	}
	defer fd.Close()
	inputdevs = new(InputDevices)
	err = json.NewDecoder(fd).Decode(inputdevs)
	return
}

func xy2rawxy(x, y int) (int32, int32) {
	w, h := iptdevs.TouchScreen.Width, iptdevs.TouchScreen.Height
	rw, rh := iptdevs.TouchScreen.RawWidth, iptdevs.TouchScreen.RawHeight
	rx := int32(float32(x) / float32(w) * float32(rw))
	ry := int32(float32(y) / float32(h) * float32(rh))
	return rx, ry
}

func clickDown(x, y int) {
	fd := devicefd
	rx, ry := xy2rawxy(x, y)

	sendevent(fd, "3", "0039", 0x0ffffff4) // tracking id
	sendevent(fd, "3", "0030", 5)          // major ?
	sendevent(fd, "1", "014a", 1)          // btn-touch down
	sendevent(fd, "3", "0035", rx)         // abs-mt-position x
	sendevent(fd, "3", "0036", ry)         // abs-mt-position y
	sendevent(fd, "3", "003a", 37)         // pressure
	sendevent(fd, "3", "0032", 4)          // ?
	sendevent(fd, "0", "0000", 0)          // sync-report
}
func clickUp() {
	fd := devicefd
	sendevent(fd, "3", "0039", -1) // tracking id
	sendevent(fd, "1", "014a", 0)  // btn-touch up
	sendevent(fd, "0", "0000", 0)  // sync-report
}

var (
	ErrArguments = errors.New("error arguments parsed")
)

func cmdInputTap(x, y int, duration time.Duration) (err error) {
	fmt.Printf("airinput tap %d %d\n", x, y)
	clickDown(x, y)
	time.Sleep(duration)
	clickUp()
	return nil
}

func cmdSwipe(x1, y1, x2, y2 int, duration time.Duration) (err error) {
	fmt.Printf("swipe %d %d   %d %d %s\n", x1, y1, x2, y2, duration)

	fd := devicefd
	move := func(x, y int, duration time.Duration) {
		rx, ry := xy2rawxy(x, y)
		sendevent(fd, "3", "0035", rx) // abs-mt-position x
		sendevent(fd, "3", "0036", ry) // abs-mt-position y
		sendevent(fd, "0", "0000", 0)  // sync-report
		if duration > 0 {
			time.Sleep(duration)
		}
	}
	mx, my := (x1+x2)/2, (y1+y2)/2
	start := time.Now()
	gap := duration / 4
	clickDown(x1, y1)                        // p1
	move((mx+x1)/2, (my+y1)/2, gap)          // p2
	move(mx, my, gap)                        // p3(middle)
	move((mx+x2)/2, (my+y2)/2, gap)          // p4
	move(x2, y2, duration-time.Since(start)) // p5
	clickUp()
	return nil
}
