package main

import (
	"fmt"
	"strings"
)

type Movie struct {
	Title string
	Year  int
}

// NormalizeTitle standardizes movie titles for comparison
func NormalizeTitle(title string) string {
	return strings.ToLower(strings.TrimSpace(title))
}

// CompareMovies checks which movies from the watchlist are in the Kodi library
func CompareMovies(watchlist []Movie, kodiMovies []Movie) {
	movieMap := make(map[string]bool)

	// Store Kodi library movies in a map for quick lookup
	for _, movie := range kodiMovies {
		key := fmt.Sprintf("%s (%d)", NormalizeTitle(movie.Title), movie.Year)
		movieMap[key] = true
	}

	// Compare watchlist with Kodi library
	fmt.Println("\nComparison Results:")
	found := false
	missingMovies := []Movie{}

	for _, movie := range watchlist {
		key := fmt.Sprintf("%s (%d)", NormalizeTitle(movie.Title), movie.Year)
		if movieMap[key] {
			fmt.Println("✔ Found in Kodi Library:", movie.Title, "(", movie.Year, ")")
			found = true
		} else {
			missingMovies = append(missingMovies, movie)
		}
	}

	if !found {
		fmt.Println("❌ No movies from your watchlist are in your Kodi library.")
	}

	// Print missing movies
	if len(missingMovies) > 0 {
		fmt.Println("\n❌ Movies in your watchlist that are NOT in your Kodi library:")
		for _, movie := range missingMovies {
			fmt.Printf("- %s (%d)\n", movie.Title, movie.Year)
		}
	}
}
