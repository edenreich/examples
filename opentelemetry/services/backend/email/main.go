package main

import (
	"log"
	"net/http"
	"time"
)

func emailHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Sending email to customer..")

	time.Sleep(300 * time.Millisecond)

	log.Println("Email sent!")
	w.Write([]byte("Email sent!"))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/email", emailHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // extermly long timeout for demo purposes
		WriteTimeout: 10 * time.Second, // extermly long timeout for demo purposes
	}

	log.Println("Server is listening on port :8080")
	log.Fatal(server.ListenAndServe())
}
