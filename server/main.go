// Package main provides ...
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		serialNo := r.URL.Query().Get("serialno")
		fmt.Println(strings.Split(r.RemoteAddr, ":")[0], serialNo)
		io.WriteString(w, "connected")
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}
