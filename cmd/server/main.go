package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// SubscribeRequest represents the structure of a subscription request
type SubscribeRequest struct {
	GraphID string `json:"graphid"`
	SimID   string `json:"simid"`
}

// Subscriptions encapsulates the subscription list and its mutex
type Subscriptions struct {
	mu            sync.Mutex
	subscriptions []SubscribeRequest
	wg            sync.WaitGroup // Wait group to ensure all requests are processed before retrieving subscriptions
}

// NewSubscriptions creates a new Subscriptions instance
func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		subscriptions: []SubscribeRequest{},
	}
}

var subscriptions = NewSubscriptions()

// Add adds a new subscription to the list in a thread-safe manner
func (s *Subscriptions) Add(req SubscribeRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscriptions = append(s.subscriptions, req)
	log.Printf("Subscription added: %+v\n", req)
	s.wg.Done() // Notify that request processing is complete
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Increment wait group counter before processing request
	subscriptions.wg.Add(1)

	// Process the request asynchronously
	go processSubscription(req)

	// Respond immediately to the client
	w.WriteHeader(http.StatusAccepted)               // Setting the status code to 202 Accepted
	w.Write([]byte("Subscription request accepted")) // Writing the response body
}

func processSubscription(req SubscribeRequest) {
	subscriptions.Add(req)
}

type unsubscribeRequest struct {
	GraphID string `json:"graphid"`
	SimID   string `json:"simid"`
}

func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req unsubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Process the request asynchronously
	go processUnsubscription(req)

	// Respond immediately to the client
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Unsubscription request accepted"))
}

func processUnsubscription(req unsubscribeRequest) {
	subscriptions.Remove(req)
}

// define subscriptions remove
func (s *Subscriptions) Remove(req unsubscribeRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, sub := range s.subscriptions {
		if sub.GraphID == req.GraphID && sub.SimID == req.SimID {
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)
			log.Printf("Subscription removed: %+v\n", req)
			return
		}
	}
	log.Printf("Subscription not found: %+v\n", req)
}

func GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Wait for all pending subscription requests to complete
	subscriptions.wg.Wait()

	subs := subscriptions.Get()
	json.NewEncoder(w).Encode(subs)
}

func (s *Subscriptions) Get() []SubscribeRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.subscriptions
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/subscribe", SubscribeHandler).Methods("POST")
	r.HandleFunc("/unsubscribe", UnsubscribeHandler).Methods("POST")
	r.HandleFunc("/getSubs", GetSubscriptionsHandler).Methods("GET")

	http.Handle("/", r)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
