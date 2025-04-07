package main

import (
	"net/http"
	"os"
)

func main() {
	_, err := os.Stat("./.cache")
	if err != nil {
		os.Mkdir("./.cache", 0700)
	}
	http.HandleFunc("/", RequestHandler)
	http.ListenAndServe(":3333", nil)
}
