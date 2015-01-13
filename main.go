// Currently this is a sample
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/netease/airinput/go-airinput"
	"github.com/sevlyar/go-daemon"
)

var (
//	m = macaron.Classic()
)

const SAVEFILE = "/data/local/tmp/cache-airinput.txt"

var (
	addr     = flag.String("addr", ":21000", "listen address")
	debug    = flag.Bool("debug", false, "enable debug")
	isDaemon = flag.Bool("daemon", false, "run as daemon")
	fix      = flag.Bool("fix", false, "fix unexpected problem caused by airinput")
	tpevent  = flag.String("i", "", "touchpad event, eg: /dev/input/event1")
	remote   = flag.String("remote", "", "remote control center, eg: 10.0.0.1:9000")
	runjs    = flag.String("runjs", "", "javascript code to run")
	test     = flag.Bool("test", false, "just test for develop")
	quite    = flag.Bool("q", false, "donot show debug info")
)

func lprintf(format string, v ...interface{}) {
	if !*quite {
		log.Printf(format, v...)
	}
}

func main() {
	flag.Parse()
	airinput.Debug(*debug)

	if *fix {
		airinput.Release()
		return
	}

	if *tpevent == "" {
		var err error
		var data []byte
		if data, err = ioutil.ReadFile(SAVEFILE); err == nil {
			if strings.HasPrefix(string(data), "/dev/input/") {
				*tpevent = strings.TrimSpace(string(data))
			}
		} else {
			if *tpevent, err = airinput.GuessTouchpad(); err != nil {
				log.Fatal(err)
			}
			if fd, err := os.Create(SAVEFILE); err == nil {
				defer fd.Close()
				fd.Write([]byte(*tpevent))
			}
		}
		lprintf("Use tpd event: %s\n", *tpevent)
	}

	// initial
	if err := airinput.Init(*tpevent); err != nil {
		lprintf("initial\n")
		log.Fatal(err)
	}

	if *test {
		lprintf("rotation: %d", airinput.Rotation())
		return
	}

	if *runjs != "" {
		RunJS(*runjs)
		return
	}

	if *remote != "" {
		r, err := http.Get("http://" + *remote + "/connect?serialno=" + url.QueryEscape(SerialNo()))
		if err == nil {
			lprintf("Remote connected\n")
			r.Body.Close()
		}
	}

	ipinfo, _ := MyIP()
	lprintf("IP: %v\n", ipinfo)
	lprintf("Listen on: %v\n", *addr)

	if *isDaemon {
		context := new(daemon.Context)
		child, _ := context.Reborn()
		if child != nil {
			println("daemon started")
		} else {
			defer context.Release()
			ServeWeb(*addr)
		}
	} else {
		ServeWeb(*addr)
		return
	}
}
