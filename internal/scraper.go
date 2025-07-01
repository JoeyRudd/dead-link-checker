package internal

import (
	"fmt"
	"io"
	"net/http"
)

// Scrape performs a web scraping operation on the given URL.
func Scrape(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP Error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP Error status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("HTTP Error reading body: %w", err)
	}
	return string(body), nil
}
