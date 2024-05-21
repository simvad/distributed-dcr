package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// SubscriptionData represents the JSON data structure for subscription requests
type SubscriptionData struct {
	GraphID string `json:"graphid"`
	SimID   string `json:"simid"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dcr-cli {subscribe|getSubs|unsubscribe} [options]")
		os.Exit(1)
	}

	command := os.Args[1]
	flags := flag.NewFlagSet("flags", flag.ExitOnError)
	graphID := flags.String("g", "", "The graph ID")
	simID := flags.String("s", "", "The simulation ID")
	baseURL := flags.String("url", "http://localhost:8080", "The base URL of the server")

	flags.Parse(os.Args[2:])

	switch command {
	case "subscribe":
		if *graphID == "" || *simID == "" {
			fmt.Println("Both graphID (-g) and simID (-s) must be provided for subscription.")
			os.Exit(1)
		}
		subscribe(*baseURL, *graphID, *simID)
	case "getSubs":
		getSubs(*baseURL)
	case "unsubscribe":
		if *graphID == "" || *simID == "" {
			fmt.Println("Both graphID (-g) and simID (-s) must be provided for unsubscription.")
			os.Exit(1)
		}
		unsubscribe(*baseURL, *graphID, *simID)
	default:
		fmt.Println("Invalid command. Please use one of: subscribe, getSubs, or unsubscribe.")
		os.Exit(1)
	}
}

func subscribe(baseURL, graphID, simID string) {
	url := fmt.Sprintf("%s/subscribe", baseURL)
	data := SubscriptionData{GraphID: graphID, SimID: simID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("Subscription successful:", resp.Status)
}

func getSubs(baseURL string) {
	url := fmt.Sprintf("%s/getSubs", baseURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	fmt.Println("Get subscriptions successful:")
	fmt.Println(string(body))
}

func unsubscribe(baseURL, graphID, simID string) {
	url := fmt.Sprintf("%s/unsubscribe", baseURL)
	data := SubscriptionData{GraphID: graphID, SimID: simID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("Unsubscription successful:", resp.Status)
}
