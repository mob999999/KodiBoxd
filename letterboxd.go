package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// getLetterboxdWatchlist scrapes the Letterboxd watchlist for a given username.
// structing release year from supplemental Letterboxd JSON
type FilmData struct {
	ReleaseYear int `json:"releaseYear"`
}

func getLetterboxdWatchlist(letterboxdUsername string) []string {
	c := colly.NewCollector()
	var wg sync.WaitGroup

	// Set User-Agent to bypass anti-bot detection.
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"

	// Add delay to avoid getting blocked.
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*letterboxd.com*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	var movies []string

	// Log when visiting pages.
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
		// Increment the WaitGroup counter when a new request starts.
		wg.Add(1)
	})

	c.OnHTML(".poster-container div[data-production-data-endpoint]", func(e *colly.HTMLElement) {
		// Extract the title from the nested <img> tag
		title := e.ChildAttr("img", "alt")
		if title == "" {
			fmt.Println("No title found")
			return
		}

		// Extract the JSON endpoint
		endpoint := e.Attr("data-production-data-endpoint")
		if endpoint == "" {
			fmt.Println("No JSON endpoint found")
			return
		}
		jsonURL := "https://letterboxd.com" + endpoint

		// Fetch JSON data
		resp, err := http.Get(jsonURL)
		if err != nil {
			fmt.Printf("Error fetching JSON: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("HTTP error: %d\n", resp.StatusCode)
			return
		}

		// Parse the JSON
		var filmData FilmData
		if err := json.NewDecoder(resp.Body).Decode(&filmData); err != nil {
			fmt.Printf("JSON decode error: %v\n", err)
			return
		}

		// Combine title and year
		movieWithYear := fmt.Sprintf("%s (%d)", title, filmData.ReleaseYear)
		movies = append(movies, movieWithYear)
	})

	// Log if no titles are found.
	c.OnScraped(func(r *colly.Response) {
		// Decrement the WaitGroup counter when a request finishes.
		wg.Done()
		if len(movies) == 0 {
			fmt.Println("No movie titles found on the page.")
		}
	})

	// Follow the next page.
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Found 'Older' button, going to:", nextPage)
		c.Visit(nextPage)
	})

	// Start scraping from page 1.
	startURL := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", letterboxdUsername)
	err := c.Visit(startURL)
	if err != nil {
		fmt.Println("Error visiting start URL, check your username in config.json", err)
	}

	// Wait for all pages to be scraped.
	wg.Wait()

	return movies
}
