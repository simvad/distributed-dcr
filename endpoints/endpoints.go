package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
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
}

// NewSubscriptions creates a new Subscriptions instance
func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		subscriptions: []SubscribeRequest{},
	}
}

// Add adds a new subscription to the list in a thread-safe manner
func (s *Subscriptions) Add(req SubscribeRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscriptions = append(s.subscriptions, req)
	log.Printf("Subscription added: %+v\n", req)
}

var subscriptions = NewSubscriptions()

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

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

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
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
