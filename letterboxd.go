package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// getLetterboxdWatchlist scrapes the Letterboxd watchlist for a given username.
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

	// Extract movie titles.
	c.OnHTML(".poster-container img", func(e *colly.HTMLElement) {
		title := e.Attr("alt")
		fmt.Println("Found title:", title)
		movies = append(movies, title)
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
		fmt.Println("Error visiting start URL, check your config-file:", err)
	}

	// Wait for all pages to be scraped.
	wg.Wait()

	return movies
}
