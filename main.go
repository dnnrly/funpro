package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	flag.Parse()
	fmt.Printf("Starting funpro server...\n")
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request %v", r)
	w.WriteHeader(200)
}
