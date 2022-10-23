package main

import (
	"log"
	"net/http"
	"time"
)

func orderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Confirming the order..")

	// Sleep the application for 2 seconds to simulate a long running process (in real world this could be a database call perhaps)
	time.Sleep(2 * time.Second)

	log.Println("Order confirmed!")
	w.Write([]byte("Order confirmed!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/order", orderHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // extermly long timeout for demo purposes
		WriteTimeout: 10 * time.Second, // extermly long timeout for demo purposes
	}

	log.Println("Server is listening on port :8080")
	log.Fatal(server.ListenAndServe())
}
