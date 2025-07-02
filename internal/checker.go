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
		// Make a head request to url
		resp, err := client.Head(url)
		if err != nil {
			fmt.Printf("Error fetching URL %s: %s\n", url, err)
			deadLinks = append(deadLinks, url) // Network error
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			deadLinks = append(deadLinks, url) // Non 2xx code
		}

	}
	return deadLinks
}
