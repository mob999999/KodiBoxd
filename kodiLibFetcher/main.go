package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Kodi API Credentials
const (
	kodiHost     = "192.168.178.42" // Replace with your Xiaomi box's IP address
	kodiPort     = "8080"           // Default Kodi JSON-RPC port
	kodiUsername = "kodi"
	kodiPassword = "kodimb"
)

type kodiRequest struct {
	JsonRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	ID      string                 `json:"id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type kodiResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type movieListResponse struct {
	Movies []struct {
		Title string `json:"title"`
		Year  int    `json:"year,omitempty"`
	} `json:"movies"`
}

// Main function - entry point of program
func main() {
	fmt.Println("Starting Kodi Library Fetcher...")

	// Create request for movies
	reqBody := kodiRequest{
		JsonRPC: "2.0",
		Method:  "VideoLibrary.GetMovies",
		ID:      "1",
		Params:  map[string]interface{}{"properties": []string{"title", "year"}},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}

	// Build the URL
	kodiURL := fmt.Sprintf("http://%s:%s/jsonrpc", kodiHost, kodiPort)

	// Create HTTP request
	req, err := http.NewRequest("POST", kodiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(kodiUsername + ":" + kodiPassword))
	req.Header.Set("Authorization", "Basic "+auth)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error contacting Kodi:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse JSON response
	var result kodiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Check for Kodi error
	if result.Error != nil {
		fmt.Printf("Kodi API error (%d): %s\n", result.Error.Code, result.Error.Message)
		return
	}

	// Parse movie list
	var movieList movieListResponse
	if err := json.Unmarshal(result.Result, &movieList); err != nil {
		fmt.Println("Error parsing movie list:", err)
		return
	}

	// Print movie titles
	fmt.Printf("Found %d movies in Kodi library:\n", len(movieList.Movies))
	for i, movie := range movieList.Movies {
		if i < 10 { // Just show the first 10 for example
			fmt.Printf("%d. %s", i+1, movie.Title)
			if movie.Year > 0 {
				fmt.Printf(" (%d)", movie.Year)
			}
			fmt.Println()
		}
	}
}
