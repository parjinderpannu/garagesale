package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Entry point of program
func main() {

	// Convert the Echo function to a type that implements http.Handler
	h := http.HandlerFunc(Echo)

	log.Println("listening on localhost:8888")
	if err := http.ListenAndServe("localhost:8888", h); err != nil {
		log.Fatal(err)
	}
}

// Echo just tells you about the request you made
func Echo(w http.ResponseWriter, r *http.Request) {

	id := rand.Intn(1000)

	fmt.Println("starting", id)

	time.Sleep(3 * time.Second)

	fmt.Fprintln(w, "You asked to", r.Method, r.URL.Path)

	fmt.Println("ending", id)
}

//https://golang.org/pkg/net/http/
//https://golang.org/pkg/net/http/#ListenAndServe
//https://golang.org/pkg/net/http/#Handler
//https://golang.org/pkg/net/http/#HandlerFunc
