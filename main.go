package main

import (
	"fmt"
)

func main() {
	// Load configuration if available
	config, err := LoadConfig() // Call LoadConfig from config.go
	if err == nil {
		fmt.Println("Configuration loaded from config.json:")
		fmt.Printf("Letterboxd Username: %s\n", config.LetterBoxdUsername) // Fixed case
		fmt.Printf("Kodi IP Address: %s\n", config.KodiIP)
		fmt.Printf("Kodi Port: %s\n", config.KodiPort)
		fmt.Printf("Kodi Username: %s\n", config.KodiUsername)
		fmt.Println("Kodi Password: [hidden]") // Security measure

		// Fetch and display Letterboxd watchlist
		movies := getLetterboxdWatchlist(config.LetterBoxdUsername) // Fixed case
		fmt.Println("\nYour Letterboxd Watchlist:")
		for _, movie := range movies {
			fmt.Println("-", movie)
		}

		// Fetch and display Kodi movies
		kodiMovies, err := getKodiMovies() // Pass config to function
		if err != nil {
			fmt.Println("Error fetching Kodi movies:", err)
			return
		}

		fmt.Println("\nYour Kodi Library Movies:")
		for _, movie := range kodiMovies {
			fmt.Printf("- %s (%d)\n", movie.Title, movie.Year)
		}

		return
	}

	// If config is missing or corrupted, prompt the user to enter details
	var newConfig Config
	fmt.Println("No valid config file found. Please enter your details.")

	fmt.Print("Letterboxd Username: ")
	fmt.Scanln(&newConfig.LetterBoxdUsername) // Fixed case

	fmt.Print("Kodi IP Address: ")
	fmt.Scanln(&newConfig.KodiIP)

	fmt.Print("Kodi Port: ")
	fmt.Scanln(&newConfig.KodiPort)

	fmt.Print("Kodi Username (leave blank if not applicable): ")
	fmt.Scanln(&newConfig.KodiUsername)

	fmt.Print("Kodi Password: ")
	fmt.Scanln(&newConfig.KodiPassword)

	// Save the newly entered configuration
	if err := SaveConfig(newConfig); err != nil { // Call SaveConfig from config.go
		fmt.Println("Error saving configuration:", err)
		return
	}

	fmt.Println("Configuration saved to config.json")
}
