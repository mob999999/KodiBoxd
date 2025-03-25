package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LetterboxdUsername string `json:"letterBoxdUsername"`
	KodiIP             string `json:"kodiIP"`
	KodiPort           string `json:"kodiPort"`
	KodiUsername       string `json:"kodiUsername,omitempty"` // omitempty to skip if empty
	KodiPassword       string `json:"kodiPassword,omitempty"` // omitempty to skip if empty
}

func main() {
	var config Config

	// Check if config.json already exists
	if _, err := os.Stat("config.json"); err == nil {
		// File exists, read the existing configuration
		file, err := os.Open("config.json")
		if err != nil {
			fmt.Println("Error opening config file:", err)
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		fmt.Println("Configuration already exists:")
		fmt.Printf("Letterboxd Username: %s\n", config.LetterboxdUsername)
		fmt.Printf("Kodi IP Address: %s\n", config.KodiIP)
		fmt.Printf("Kodi Port: %s\n", config.KodiPort)
		fmt.Printf("Kodi Username: %s\n", config.KodiUsername)
		fmt.Println("Kodi Password: [hidden]") // Do not print the password for security reasons
		return
	}

	// If the config file does not exist, prompt for input
	fmt.Print("Letterboxd Username: ")
	fmt.Scanln(&config.LetterboxdUsername)

	fmt.Print("Kodi IP Address: ")
	fmt.Scanln(&config.KodiIP)

	fmt.Print("Kodi Port: ")
	fmt.Scanln(&config.KodiPort)

	fmt.Print("Kodi Username (leave blank if not applicable): ")
	fmt.Scanln(&config.KodiUsername)

	fmt.Print("Kodi Password: ")
	fmt.Scanln(&config.KodiPassword)

	// Create or open the config.json file
	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}
	defer file.Close()

	// Encode the config struct to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(config); err != nil {
		fmt.Println("Error writing to config file:", err)
		return
	}

	fmt.Println("Configuration saved to config.json")
}
