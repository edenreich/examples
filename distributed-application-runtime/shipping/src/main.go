package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
)

type LineItem struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Quantity    uint8     `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
}

type Billing struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Address string    `json:"address,omitempty"`
}

type Shipping struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Address string    `json:"address,omitempty"`
}

type Customer struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type Order struct {
	Id        uuid.UUID  `json:"id"`
	Number    string     `json:"number"`
	LineItems []LineItem `json:"line_items"`
	Billing   Billing    `json:"billing"`
	Shipping  Shipping   `json:"shipping"`
	Customer  Customer   `json:"customer"`
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/orders", eventHandler)
	log.Println("Starting the server on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	event := cloudevents.New()
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Println("Error decoding event: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = event.Validate()
	if err != nil {
		log.Println("Error validating event: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := event.DataBytes()
	if err != nil {
		log.Println("Error getting event data as bytes: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order Order
	err = json.Unmarshal(payload, &order)
	if err != nil {
		log.Println("Error unmarshalling event data from bytes: ", order)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Shipping Order %s to Customer %s ...\n", order.Id, order.Customer.Email)

	w.WriteHeader(http.StatusOK)
}
