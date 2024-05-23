package dcrs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Config struct to hold the configuration details
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Function to load the configuration from a file
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

func apiCall(url, method string, jsonBody map[string]interface{}, config *Config) (string, http.Header, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if jsonBody != nil {
		jsonData, err := json.Marshal(jsonBody)
		if err != nil {
			return "", nil, fmt.Errorf("could not marshal json: %v", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", nil, fmt.Errorf("could not create request: %v", err)
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return "", nil, fmt.Errorf("could not create request: %v", err)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(config.Username, config.Password)

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("could not read response body: %v", err)
	}

	return string(body), resp.Header, nil
}

// Function to get the simulation ID
func GetSimId(graphID string, config *Config) (string, error) {
	url := fmt.Sprintf("https://repository.dcrgraphs.net/api/graphs/%s/sims", graphID)
	method := "POST"

	_, headers, err := apiCall(url, method, nil, config)
	if err != nil {
		return "", err
	}

	// Extracting the 'x-dcr-simulation-id' from the header
	simulationID := headers.Get("x-dcr-simulation-id")
	if simulationID == "" {
		return "", fmt.Errorf("x-dcr-simulation-id header not found in the response")
	}

	return simulationID, nil
}

// Function to get the relations using the simulation ID
func GetRelations(graphID, simID string, config *Config) (string, error) {
	url := fmt.Sprintf("https://repository.dcrgraphs.net/api/graphs/%s/sims/%s/relations", graphID, simID)
	method := "GET"

	responseText, _, err := apiCall(url, method, nil, config)
	if err != nil {
		return "", err
	}

	return responseText, nil
}
