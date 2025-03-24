package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func getWatchlist(username string) []string {
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
		fmt.Println("Found title:", title) // Debug output
		movies = append(movies, title)
	})

	// Debug: Print "Older" button detection & follow next page
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

func main() {
	username := "mc22k" // Change to your Letterboxd username
	watchlist := getWatchlist(username)

	fmt.Println("\nFull Watchlist for", username)
	for _, movie := range watchlist {
		fmt.Println("-", movie)
	}
}
