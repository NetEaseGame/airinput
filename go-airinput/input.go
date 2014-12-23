package airinput

import (
	// #include "input.h"
	"C"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/codeskyblue/comtool"
)

var goDebug = false
var rawWidth, rawHeight int

func dprintf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if goDebug {
		fmt.Printf(format, args...)
	}
}

func atoi(a string) int {
	var i int
	_, err := fmt.Sscanf(a, "%d", &i)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func real2raw(x, y int) (rx, ry int) {
	w, h, _ := ScreenSize()
	if w > h {
		w, h = h, w
	}
	x = x * rawWidth / w
	y = y * rawHeight / h
	return x, y
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

func GuessTouchpad() (p string, err error) {
	for i := 0; i < 10; i++ {
		event := fmt.Sprintf("/dev/input/event%d", i)
		if _, _, err := GetRawSize(event); err != nil {
			continue
		} else {
			return event, nil
		}
	}
	return "", errors.New("touchpad event not found")
}

func GuessTouchpadByName() (p string, err error) {
	var id int = -1
	for i := 0; i < 10; i++ {
		namePath := fmt.Sprintf("/sys/class/input/event%d/device/name", i)
		if !comtool.Exists(namePath) {
			continue
		}
		data, err := ioutil.ReadFile(namePath)
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(data))
		dprintf("event%d: name: %s", i, name)
		// $name may have Touchscreen and touchscreen
		// atmel-maxtouch: Xiaomi2
		for _, possibleName := range []string{"ouchscreen$", "synaptics-rmi-ts", "ist30xx_ts_input", "mtk-tpd$", "atmel-maxtouch"} {
			re := regexp.MustCompile(possibleName)
			if re.MatchString(name) {
				id = i
				break
			}
		}
		if id != -1 {
			break
		}
	}
	dprintf("eventid: %d", id)
	if id == -1 {
		return "", errors.New("cannot autodetect touchpad event")
	}
	return fmt.Sprintf("/dev/input/event%d", id), nil
}

// Have to be call before other func
// if tpdEvent == "", it will auto guess
func Init(tpdEvent string) (err error) {
	if tpdEvent == "" {
		tpdEvent, err = GuessTouchpad()
		if err != nil {
			return
		}
	}
	rawWidth, rawHeight, err = GetRawSize(tpdEvent)
	if err != nil {
		return err
	}
	C.input_init(C.CString(tpdEvent))
	return nil
}

// Wether show debug info
func Debug(state bool) {
	flag := C.int(0)
	if state {
		flag = C.int(1)
	}
	goDebug = state
	C.set_debug(flag)
}

// Tap position for some time
func Tap(x, y int, duration time.Duration) {
	x, y = real2raw(x, y)
	msec := int(duration.Nanoseconds() / 1e6)
	fmt.Println(x, y)
	C.tap(C.int(x), C.int(y), C.int(msec))
}

// Drag form A to B
func Drag(startX, startY int, endX, endY int, steps int, duration time.Duration) {
	msec := int(duration.Nanoseconds() / 1e6)
	C.drag(C.int(startX), C.int(startY), C.int(endX), C.int(endY), C.int(steps), C.int(msec))
}

// This is like two drag
// Can implements like shrink and magnify
func Pinch(Ax0, Ay0, Ax1, Ay1 int,
	Bx0, By0, Bx1, By1 int, steps int, duration time.Duration) {
	msec := int(duration.Nanoseconds() / 1e6)
	C.pinch(C.int(Ax0), C.int(Ay0), C.int(Ax1), C.int(Ay1),
		C.int(Bx0), C.int(By0), C.int(Bx1), C.int(By1), C.int(steps), C.int(msec))
}

func Release() {
	C.release()
}
