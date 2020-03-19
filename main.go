package main

import (
	"fmt"
	"log"
	"net/http"
)

// Entry point of program
func main() {

	h := http.HandlerFunc(Echo)

	log.Println("listening on localhost:8888")
	if err := http.ListenAndServe("localhost:8888", h); err != nil {
		log.Fatal(err)
	}
}

// Echo just tells you about the request you made
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You asked to", r.Method, r.URL.Path)
}

//https://golang.org/pkg/net/http/
//https://golang.org/pkg/net/http/#ListenAndServe
//https://golang.org/pkg/net/http/#Handler
//https://golang.org/pkg/net/http/#HandlerFunc
