package airinput

import (
	// #include "input.h"
	"C"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/codeskyblue/comtool"
)

var goDebug = false

func dprintf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if goDebug {
		fmt.Printf(format, args...)
	}
}

// Have to be call before other func
func Init() error {
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
		for _, possibleName := range []string{"ouchscreen$", "-tpd$", "atmel-maxtouch"} {
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
		return errors.New("cannot autodetect touchpad event")
	}
	C.input_init(C.CString(fmt.Sprintf("/dev/input/event%d", id)))
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
	msec := int(duration.Nanoseconds() / 1e6)
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
