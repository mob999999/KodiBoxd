package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Kodi API Credentials
const (
	kodiHost     = "192.168.178.42" // Replace with your Kodi-Instances IP address
	kodiPort     = "8080"           // Default Kodi JSON-RPC port
	kodiUsername = "kodi"
	kodiPassword = "kodimb"
)

// Letterboxd username
const letterboxdUsername = "mc22k" // Your Letterboxd username

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
	Movies []KodiMovie `json:"movies"`
}

type KodiMovie struct {
	Title string `json:"title"`
	Year  int    `json:"year,omitempty"`
}

// Function to fetch movies from Kodi library
func getKodiMovies() ([]KodiMovie, error) {
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
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Build the URL
	kodiURL := fmt.Sprintf("http://%s:%s/jsonrpc", kodiHost, kodiPort)

	// Create HTTP request
	req, err := http.NewRequest("POST", kodiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(kodiUsername + ":" + kodiPassword))
	req.Header.Set("Authorization", "Basic "+auth)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error contacting Kodi: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Parse JSON response
	var result kodiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// Check for Kodi error
	if result.Error != nil {
		return nil, fmt.Errorf("Kodi API error (%d): %s", result.Error.Code, result.Error.Message)
	}

	// Parse movie list
	var movieList movieListResponse
	if err := json.Unmarshal(result.Result, &movieList); err != nil {
		return nil, fmt.Errorf("error parsing movie list: %w", err)
	}

	return movieList.Movies, nil
}

// Function to scrape Letterboxd watchlist
func getLetterboxdWatchlist(username string) []string {
	c := colly.NewCollector()

	// Set User-Agent to bypass anti-bot detection
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"

	// Add delay to avoid getting blocked
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*letterboxd.com*",
		Delay:       2 * time.Second, // Delay between requests
		RandomDelay: 1 * time.Second, // Randomize delay slightly
	})

	var movies []string

	// Log when visiting pages
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	// Extract movie titles
	c.OnHTML(".poster-container img", func(e *colly.HTMLElement) {
		title := e.Attr("alt")
		fmt.Println("Found title:", title)
		movies = append(movies, title)
	})

	// Follow next page
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Found 'Older' button, going to:", nextPage)
		c.Visit(nextPage)
	})

	// Start scraping from page 1
	startURL := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", username)
	c.Visit(startURL)

	return movies
}

// Function to extract year from movie title (if present)
// Example: "The Matrix (1999)" -> "The Matrix", "1999"
func extractYearFromTitle(title string) (string, string) {
	// Check if title ends with a year in parentheses
	if strings.HasSuffix(title, ")") {
		openParenIndex := strings.LastIndex(title, "(")
		if openParenIndex != -1 {
			year := title[openParenIndex+1 : len(title)-1]
			// Verify that the extracted part is a 4-digit number
			if len(year) == 4 && strings.Count(year, "0")+strings.Count(year, "1")+
				strings.Count(year, "2")+strings.Count(year, "3")+strings.Count(year, "4")+
				strings.Count(year, "5")+strings.Count(year, "6")+strings.Count(year, "7")+
				strings.Count(year, "8")+strings.Count(year, "9") == 4 {
				return title[:openParenIndex-1], year
			}
		}
	}
	return title, "" // No year found
}

// Function to normalize movie titles for comparison
func normalizeTitle(title string) string {
	// Extract year from title and remove it (if present)
	titleOnly, _ := extractYearFromTitle(title)

	// Convert to lowercase
	titleOnly = strings.ToLower(titleOnly)

	// Remove common prefixes like "The ", "A ", etc.
	prefixes := []string{"the ", "a ", "an "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(titleOnly, prefix) {
			titleOnly = titleOnly[len(prefix):]
		}
	}

	// Remove special characters and extra spaces
	specialChars := []string{":", "-", ",", ".", "'", "\"", "!", "?", "(", ")", "[", "]"}
	for _, char := range specialChars {
		titleOnly = strings.ReplaceAll(titleOnly, char, "")
	}

	// Replace multiple spaces with a single space
	for strings.Contains(titleOnly, "  ") {
		titleOnly = strings.ReplaceAll(titleOnly, "  ", " ")
	}

	return strings.TrimSpace(titleOnly)
}

func main() {
	fmt.Println("Starting Letterboxd Watchlist vs Kodi comparison...")

	// Get username from command line or use default
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
	fmt.Println("\n= = = = = = = = = = = = = = = = = = = = = = = = = = = = = =")
	fmt.Println("COMPARISON RESULTS")
	fmt.Println("= = = = = = = = = = = = = = = = = = = = = = = = = = = = = =")

	// Movies in both lists
	fmt.Printf("\nMovies in both watchlist and Kodi library (%d):\n", len(inBoth))
	for i, movie := range inBoth {
		fmt.Printf("  %d. %s\n", i+1, movie)
	}

	// Movies only in watchlist
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
