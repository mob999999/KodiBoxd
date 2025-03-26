package main

import (
	"fmt"
	"strings"
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
			fmt.Println("Error fetching Kodi movies. Is your Kodi Device running and reachable? Is your config.json properly set up? ", err)
			return
		}

		fmt.Println("\nYour Kodi Library Movies:")
		for _, movie := range kodiMovies {
			fmt.Printf("- %s (%d)\n", movie.Title, movie.Year)
		}
		// Convert []KodiMovie to []Movie
		convertedKodiMovies := make([]Movie, len(kodiMovies))
		for i, km := range kodiMovies {
			convertedKodiMovies[i] = Movie{Title: km.Title, Year: km.Year}
		}
		// Convert Letterboxd entries to Movie structs
		convertedMovies := make([]Movie, 0, len(movies))
		for _, entry := range movies {
			// Split into title and year using parenthesis
			idx := strings.LastIndex(entry, " (")
			if idx == -1 || !strings.HasSuffix(entry, ")") {
				fmt.Printf("Failed to parse: %s\n", entry)
				continue
			}

			title := strings.TrimSpace(entry[:idx])
			yearStr := entry[idx+2 : len(entry)-1] // Extract "YYYY"

			var year int
			if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
				fmt.Printf("Failed to parse year in: %s\n", entry)
				continue
			}

			convertedMovies = append(convertedMovies, Movie{
				Title: title,
				Year:  year,
			})
		}
		// Call the comparison function with the converted list
		CompareMovies(convertedMovies, convertedKodiMovies)
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

	fmt.Println("Configuration saved to config.json. Please exit and rerun KodiBoxd for changes to take effect.")

	// Keep terminal open
	fmt.Print("Press Enter to exit...")
	var dummy string
	fmt.Scanln(&dummy) // Wait for user input
}
