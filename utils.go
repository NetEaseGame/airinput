// Package main provides ...
package main

import (
	"encoding/json"
	"net/http"

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

func MyIP() (ip string, err error) {
	ips, err := comtool.GetLocalIPs()
	if err != nil {
		return
	}
	ip = ""
	for _, i := range ips {
		ip = ip + i.String() + "\n"
	}
	return
}
