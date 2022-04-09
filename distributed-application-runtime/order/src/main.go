package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
)

const MAX_LINE_ITEMS = 10

// todo: replace this with dapr state store
var InMemoryDB = []*Order{}

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

func NewLineItem() *LineItem {
	rand.Seed(time.Now().UnixNano())
	li := LineItem{}
	li.Id = uuid.New()
	li.Description = "Item " + strconv.Itoa(rand.Intn(100))
	li.UnitPrice = float64(rand.Intn(100)) / 100
	li.Quantity = uint8(rand.Intn(100))
	return &li
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

func NewOrder() *Order {
	order := Order{}
	order.Id = uuid.New()
	order.Number = "Order " + strconv.Itoa(rand.Intn(100))
	order.Customer.Id = uuid.New()
	order.Customer.FirstName = "John"
	order.Customer.LastName = "Doe"
	order.Customer.Email = "john.doe@gmail.com"
	order.Shipping.Id = uuid.New()
	order.Shipping.Address = "123 Main St, San Francisco, CA"
	order.Billing.Id = uuid.New()
	order.Billing.Address = "123 Main St, San Francisco, CA"
	for i := 1; i < rand.Intn(MAX_LINE_ITEMS); i++ {
		order.LineItems = append(order.LineItems, *NewLineItem())
	}
	return &order
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/bills", billsHandler)

	log.Println("Server started listening on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		randomOrder := NewOrder()
		InMemoryDB = append(InMemoryDB, randomOrder)
		log.Println("Created order: " + randomOrder.Number)

		event := cloudevents.New()
		event.SetID(randomOrder.Id.String())
		event.SetType("order.created")
		event.SetSource("http://localhost:8080/orders")
		event.SetData(randomOrder)

		payload, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		outputs := []string{}
		topics := strings.Split(os.Getenv("TOPICS"), ",")
		for _, topic := range topics {
			// https://docs.dapr.io/reference/api/pubsub_api/#publish-a-message-to-a-given-topic
			// Here is where we publish the message to the topic via dapr proxy server
			// dapr could be integrated with a pub/sub system like redis, kafka, rabbitmq, google pubsub, azure eventhubs, amazon SQS, kinesis etc..
			// without us needing to change anything on the application.
			endpoint := fmt.Sprintf("http://localhost:3500/%s/publish/%s/%s", os.Getenv("DAPR_VERSION"), os.Getenv("PUBSUB_NAME"), topic)
			response, err := http.Post(endpoint, "application/cloudevents+json", strings.NewReader(string(payload)))
			if err != nil {
				panic(err)
			}
			output := fmt.Sprintf("Publishing a message on %s on topic %s. Payload - %s", os.Getenv("PUBSUB_NAME"), topic, payload)
			log.Println(output)

			if response.StatusCode != http.StatusNoContent {
				panic(fmt.Sprintf("Error publishing message on topic %s. Status - %s", topic, response.Status))
			}

			outputs = append(outputs, output)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strings.Join(outputs, "\n")))
		return
	case "GET":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(InMemoryDB)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func billsHandler(w http.ResponseWriter, r *http.Request) {
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

	// If the order has been paid, we can publish a cloud event to the all of the subscribers
	// to orders topic with event of type "order.paid".
	if order.PaidAt.IsZero() {
		return
	}

	event.SetType("order.paid")
	event.SetSource("http://localhost:8080/orders")
	event.SetData(order)

	p, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	topics := strings.Split(os.Getenv("TOPICS"), ",")
	for _, topic := range topics {
		endpoint := fmt.Sprintf("http://localhost:3500/%s/publish/%s/%s", os.Getenv("DAPR_VERSION"), os.Getenv("PUBSUB_NAME"), topic)
		response, err := http.Post(endpoint, "application/cloudevents+json", strings.NewReader(string(p)))
		if err != nil {
			panic(err)
		}

		log.Printf("Publishing a message on %s on topic %s. Payload - %s\n", os.Getenv("PUBSUB_NAME"), topic, p)
		if response.StatusCode != http.StatusNoContent {
			panic(fmt.Sprintf("Error publishing message on topic %s. Status - %s", topic, response.Status))
		}
	}

	w.WriteHeader(http.StatusOK)
}
