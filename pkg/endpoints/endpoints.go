package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type SubscribeRequest struct {
	GraphID string `json:"graphid"`
	SimID   string `json:"simid"`
}

type Server struct {
	Mu            sync.Mutex
	Subscriptions []SubscribeRequest
	Wg            sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		Subscriptions: []SubscribeRequest{}}
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Increment wait group counter before processing request
	s.Wg.Add(1)

	// Process the request asynchronously
	go s.processSubscription(req)

	// Respond immediately to the client
	w.WriteHeader(http.StatusAccepted)               // Setting the status code to 202 Accepted
	w.Write([]byte("Subscription request accepted")) // Writing the response body
}

func (s *Server) processSubscription(req SubscribeRequest) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Subscriptions = append(s.Subscriptions, req)
	log.Printf("Subscription added: %+v\n", req)
	s.Wg.Done() // Notify that request processing is complete
}

func (s *Server) GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Wait for all pending subscription requests to complete
	s.Wg.Wait()

	s.Mu.Lock()
	defer s.Mu.Unlock()
	subs := s.Subscriptions
	json.NewEncoder(w).Encode(subs)
}

func (s *Server) UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Process the request asynchronously
	go s.processUnsubscription(req)

	// Respond immediately to the client
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Unsubscription request accepted"))
}

func (s *Server) processUnsubscription(req SubscribeRequest) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i, sub := range s.Subscriptions {
		if sub.GraphID == req.GraphID && sub.SimID == req.SimID {
			s.Subscriptions = append(s.Subscriptions[:i], s.Subscriptions[i+1:]...)
			log.Printf("Subscription removed: %+v\n", req)
			return
		}
	}
	log.Printf("Subscription not found: %+v\n", req)
}
