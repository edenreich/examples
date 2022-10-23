package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./public")))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // extermly long timeout for demo purposes
		WriteTimeout: 10 * time.Second, // extermly long timeout for demo purposes
	}

	log.Println("Server is listening on port :8080")
	log.Fatal(server.ListenAndServe())
}
