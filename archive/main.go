package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting Letterboxd Watchlist vs Kodi comparison...")

	// Get username from command line or use default from letterboxd.go
	username := letterboxdUsername
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	// Load Letterboxd watchlist
	fmt.Printf("Fetching watchlist for user %s from Letterboxd...\n", username)
	letterboxdMovies := getLetterboxdWatchlist(username)
	fmt.Printf("Found %d movies in Letterboxd watchlist\n", len(letterboxdMovies))

	// Load Kodi library
	fmt.Println("Fetching movies from Kodi library...")
	kodiMovies, err := getKodiMovies()
	if err != nil {
		fmt.Println("Error fetching Kodi library:", err)
		return
	}
	fmt.Printf("Found %d movies in Kodi library\n", len(kodiMovies))

	// Create normalized map of Kodi movies for faster lookup
	kodiMovieMap := make(map[string]KodiMovie)
	for _, movie := range kodiMovies {
		normalizedTitle := normalizeTitle(movie.Title)
		kodiMovieMap[normalizedTitle] = movie
	}

	// Compare the watchlist with Kodi library
	var inBoth, onlyInWatchlist []string

	for _, lbMovieTitle := range letterboxdMovies {
		normalizedTitle := normalizeTitle(lbMovieTitle)
		if _, exists := kodiMovieMap[normalizedTitle]; exists {
			inBoth = append(inBoth, lbMovieTitle)
		} else {
			onlyInWatchlist = append(onlyInWatchlist, lbMovieTitle)
		}
	}

	// Print results
	fmt.Println("\n=== COMPARISON RESULTS ===")

	fmt.Printf("\nMovies in both watchlist and Kodi library (%d):\n", len(inBoth))
	for i, movie := range inBoth {
		fmt.Printf("  %d. %s\n", i+1, movie)
	}

	fmt.Printf("\nMovies in watchlist but not in Kodi library (%d):\n", len(onlyInWatchlist))
	for i, movie := range onlyInWatchlist {
		fmt.Printf("  %d. %s\n", i+1, movie)
	}

	// Save movies not in Kodi to a file
	outputFile := "movies_to_add.txt"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
	} else {
		defer file.Close()
		for _, movie := range onlyInWatchlist {
			file.WriteString(movie + "\n")
		}
		fmt.Printf("\nList of movies to add has been saved to %s\n", outputFile)
	}
}
