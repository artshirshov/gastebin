package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	log.Println("Starting server :8080")
	log.Fatal(s.ListenAndServe())
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!")
}
