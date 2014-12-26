// Currently this is a sample
package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/netease/airinput/go-airinput"
	"github.com/sevlyar/go-daemon"
)

var (
//	m = macaron.Classic()
)

var (
	addr     = flag.String("addr", ":21000", "listen address")
	debug    = flag.Bool("debug", false, "enable debug")
	isDaemon = flag.Bool("daemon", false, "run as daemon")
	fix      = flag.Bool("fix", false, "fix unexpected problem caused by airinput")
	tpevent  = flag.String("i", "", "touchpad event, eg: /dev/input/event1")
	remote   = flag.String("remote", "", "remote control center, eg: 10.0.0.1:9000")
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
		log.Printf("Use tpd event: %s\n", *tpevent)
	}

	// initial
	if err := airinput.Init(*tpevent); err != nil {
		log.Println("initial")
		log.Fatal(err)
	}

	if *runjs != "" {
		RunJS(*runjs)
		return
	}

	if *remote != "" {
		r, err := http.Get("http://" + *remote + "/connect?serialno=" + url.QueryEscape(SerialNo()))
		if err == nil {
			log.Printf("Remote connected\n")
			r.Body.Close()
		}
	}

	ipinfo, _ := MyIP()
	log.Printf("IP: %v\n", ipinfo)
	log.Printf("Listen on: %v\n", *addr)

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
