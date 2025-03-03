package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/", RequestHandler)
	http.ListenAndServe(":3333", nil)
}
