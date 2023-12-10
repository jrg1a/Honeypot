package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	//TODO: implement API for honey pot server (fetch logs, send logs to a SIEM, etc.)
	http.ListenAndServe(":8080", nil)
}
