package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	PaidAt    time.Time  `json:"paid_at"`
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

	log.Printf("Billing Order %s of Customer %s ...\n", order.Id, order.Customer.Email)
	// Do the billing, simulate an API call to a payment provider to charge the customer, assuming the was charged successfully
	log.Printf("Order %s of Customer %s billed successfully\n", order.Id, order.Customer.Email)

	order.PaidAt = time.Now()

	event.SetID(order.Id.String())
	event.SetType("order.paid")
	event.SetSource("http://localhost:8080/bills")
	event.SetData(order)

	p, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	// Publish an event order paid
	topics := strings.Split(os.Getenv("TOPICS"), ",")
	for _, topic := range topics {
		endpoint := fmt.Sprintf("http://localhost:3500/%s/publish/%s/%s", os.Getenv("DAPR_VERSION"), os.Getenv("PUBSUB_NAME"), topic)
		response, err := http.Post(endpoint, "application/cloudevents+json", strings.NewReader(string(p)))
		if err != nil {
			panic(err)
		}
		log.Printf("Publishing a message on %s on topic %s. Payload - %s", os.Getenv("PUBSUB_NAME"), topic, p)
		if response.StatusCode != http.StatusNoContent {
			panic(fmt.Sprintf("Error publishing message on topic %s. Status - %s", topic, response.Status))
		}
	}

	w.WriteHeader(http.StatusOK)
}
