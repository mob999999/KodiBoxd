package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// kodiRequest defines the structure for Kodi API requests.
type kodiRequest struct {
	JsonRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	ID      string                 `json:"id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// kodiResponse defines the structure for Kodi API responses.
type kodiResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// movieListResponse defines the expected structure of the movies list.
type movieListResponse struct {
	Movies []KodiMovie `json:"movies"`
}

// KodiMovie represents a movie in Kodi's library.
type KodiMovie struct {
	Title string `json:"title"`
	Year  int    `json:"year,omitempty"`
}

// getKodiMovies fetches movies from the Kodi library.
func getKodiMovies() ([]KodiMovie, error) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Create request for movies.
	reqBody := kodiRequest{
		JsonRPC: "2.0",
		Method:  "VideoLibrary.GetMovies",
		ID:      "1",
		Params:  map[string]interface{}{"properties": []string{"title", "year"}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	kodiURL := fmt.Sprintf("http://%s:%s/jsonrpc", config.KodiIP, config.KodiPort)
	req, err := http.NewRequest("POST", kodiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication only if credentials are provided
	if config.KodiUsername != "" && config.KodiPassword != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(config.KodiUsername + ":" + config.KodiPassword))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error contacting Kodi: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result kodiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("Kodi API error (%d): %s", result.Error.Code, result.Error.Message)
	}

	var movieList movieListResponse
	if err := json.Unmarshal(result.Result, &movieList); err != nil {
		return nil, fmt.Errorf("error parsing movie list: %w", err)
	}

	return movieList.Movies, nil
}
