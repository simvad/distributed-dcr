package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/subscribe", subscribeHandler).Methods("POST")
	r.handleFunc("/unsubscribe", unsubscribeHandler).Methods("POST")
	r.HandleFunc("/subscriptions", subscriptionsHandler).Methods("POST")

	http.Handle("/", r)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
