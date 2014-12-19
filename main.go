// Currently this is a sample
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/netease/airinput/go-airinput"
	"github.com/sevlyar/go-daemon"
)

var (
//	m = macaron.Classic()
)

var (
	addr     = flag.String("addr", ":21000", "server listen address")
	debug    = flag.Bool("debug", false, "enable debug")
	isDaemon = flag.Bool("daemon", false, "run as daemon")
	fix      = flag.Bool("fix", false, "fix unexpected problem caused by airinput")
	tpevent  = flag.String("i", "", "touchpad event, eg: /dev/input/event1")
	runjs    = flag.String("runjs", "", "javascript code to run")
)

func main() {
	flag.Parse()
	airinput.Debug(*debug)

	if *fix {
		airinput.Release()
		return
	}

	if *tpevent == "" {
		var err error
		*tpevent, err = airinput.GuessTouchpad()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Use tpd event: %s\n", *tpevent)
	}

	// initial
	if err := airinput.Init(*tpevent); err != nil {
		log.Fatal(err)
	}

	if *runjs != "" {
		RunJS(*runjs)
		return
	}

	ipinfo, _ := MyIP()
	fmt.Printf("IP: %v\n", ipinfo)
	fmt.Printf("Listen on: %v\n", *addr)

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
	// useless

	/*
		img, err := airinput.TakeSnapshot()
		if err != nil {
			log.Fatal(err)
			return
		}
		fd, err := os.Create("/data/local/tmp/air.png")
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		png.Encode(fd, img)
		return
	*/
}
