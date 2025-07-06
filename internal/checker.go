package internal

import (
	"fmt"
	"net/http"
	"time"
)

// GetDeadLinks checks for dead urls
func GetDeadLinks(urls *[]string) []string {
	// Holds dead links
	var deadLinks []string

	// Client with timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for _, url := range *urls {
		// Create a head request to url
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.printf("Error creating request for URL %s: %s\n", url, err)
			deadLinks = append(deadLinks, url) // Network error
			continue
		}

		// Add user agent header
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		// Execute request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Errorf("Error getting URL %s: %s\n", url, err)
			deadLinks = append(deadLinks, url)
			continue
		}
		resp.Body.Close() // close response body

		// After following a redirect, only treat 4xx or 5xx as dead
		if resp.StatusCode >= 400 {
			deadLinks = append(deadLinks, url) // Non 2xx code
		}

	}
	return deadLinks
}
