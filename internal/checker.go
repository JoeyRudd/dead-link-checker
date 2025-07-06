package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// GetDeadLinks checks for dead urls
func GetDeadLinks(urls *[]string) []string {
	// Holds dead links
	var deadLinks []string
	// Protext deadLinks from race conditions
	var mu sync.Mutex
	// Waits for all goroutines to finihs
	var wg sync.WaitGroup

	// Client with timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for _, url := range *urls {
		// Tell waitgroup a goroutine is starting
		wg.Add(1)

		go func(url string) {
			// Tell waitgroup this curren goroutine is complete
			defer wg.Done()

			// Create a head request to url
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request for URL %s: %s\n", url, err)
				mu.Lock()
				deadLinks = append(deadLinks, url) // Network error
				mu.Unlock()
				return
			}

			// Add user agent header
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

			// Execute request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error getting URL %s: %s\n", url, err)
				mu.Lock()
				deadLinks = append(deadLinks, url)
				mu.Unlock()
				return
			}
			resp.Body.Close() // close response body

			// After following a redirect, only treat 4xx or 5xx as dead
			if resp.StatusCode >= 400 {
				mu.Lock()
				deadLinks = append(deadLinks, url) // Non 2xx code
				mu.Unlock()
			}

		}(url)

	}

	// Wait for all the goroutines to be done then return
	wg.Wait()
	return deadLinks
}
