package main

import (
	"distributed-dcr/pkg/endpoints"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	server := endpoints.NewServer()

	r := mux.NewRouter()

	r.Handle("/subscribe", http.HandlerFunc(server.SubscribeHandler)).Methods("POST")
	r.Handle("/subscriptions", http.HandlerFunc(server.GetSubscriptionsHandler)).Methods("GET")
	r.Handle("/unsubscribe", http.HandlerFunc(server.UnsubscribeHandler)).Methods("POST")

	http.Handle("/", r)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
