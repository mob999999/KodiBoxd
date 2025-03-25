package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct to hold the configuration data
type Config struct {
	LetterBoxdUsername string `json:"LetterBoxdUsername"` // Fixed case
	KodiIP             string `json:"KodiIP"`
	KodiPort           string `json:"KodiPort"`
	KodiUsername       string `json:"KodiUsername,omitempty"`
	KodiPassword       string `json:"KodiPassword,omitempty"`
}

// LoadConfig reads the configuration from `config.json`
func LoadConfig() (*Config, error) {
	var config Config

	// Check if config.json exists
	file, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	// Decode JSON into Config struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &config, nil
}

// SaveConfig writes the given config to `config.json`
func SaveConfig(config Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}
