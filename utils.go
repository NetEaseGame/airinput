// Package main provides ...
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/codeskyblue/comtool"
)

const IFCONFIG_ME = "http://ifconfig.mt.nie.netease.com"

type IfconfigMe struct {
	RealIps []string `json:"X-Real-Ip"`
}

// Get my IP
// ... this is a local function
func MyIP2() (ip string, err error) {
	resp, err := http.Get(IFCONFIG_ME + "/all.json")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var ifm IfconfigMe
	err = json.NewDecoder(resp.Body).Decode(&ifm)
	if err != nil {
		return
	}
	return ifm.RealIps[0], nil
}

func MyIP() (ip []string, err error) {
	ips, err := comtool.GetLocalIPs()
	if err != nil {
		return
	}
	ip = make([]string, 0)
	for _, i := range ips {
		if strings.Contains(i.String(), ".") {
			ip = append(ip, i.String())
		}
	}
	return
}

func SerialNo() string {
	fd, err := os.Open("/sys/class/android_usb/android0/iSerial")
	if err != nil {
		return ""
	}
	data, _ := ioutil.ReadAll(fd)
	return string(data)
}
