package main

// #include "input.h"
import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)
import "github.com/codeskyblue/comtool"

func init() {
	var id int = -1
	for i := 0; i < 10; i++ {
		namePath := fmt.Sprintf("/sys/class/input/event%d/device/name", i)
		if !comtool.Exists(namePath) {
			break
		}
		data, err := ioutil.ReadFile(namePath)
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(data))
		fmt.Printf("event%d: name: %s\n", i, name)
		for _, possibleName := range []string{"touchscreen$", "-tpd$"} {
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
	fmt.Println("eventid:", id)
	if id == -1 {
		log.Fatal("can autodetect touchpad event")
	}
	C.input_init(C.CString(fmt.Sprintf("/dev/input/event%d", id)))
}

func main() {
	//C.tap(C.int(154), C.int(1082), 5000) // 1000 = 1s
	fmt.Println(ScreenSize())
	w, _ := ScreenSize()
	fmt.Println(w / 6 * 5)
	x, y := C.int(w/6), C.int(300)
	//C.drag(x, y, C.int(945), y, 10, 5000)
	mx := C.int(w / 2)
	C.pinch(x, y, mx, y, C.int(w/6*5), y, mx, y, 5, 1000)
	time.Sleep(time.Second * 1)
	C.pinch(mx, y, x, y, mx, y, C.int(w/6*5), y, 5, 1000)
}
