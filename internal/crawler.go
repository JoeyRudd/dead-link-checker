package internal

import (
	"net/url"
)

func CrawlSite(startURL string, maxDepth int) []string {
	visited := make(map[string]bool)
	var allLinks []string

	// Call recursive function
	crawlRecursive(startURL, startURL, 0, maxDepth, visited, &allLinks)

	return allLinks

}

func crawlRecursive(currentURL, baseURL string, depth, maxDepth int, visited map[string]bool, allLinks *[]string) {
	// Stop if too deep
	if depth > maxDepth {
		return
	}

	// Stop if URL has already been visited
	if visited[currentURL] {
		return
	}

	// Mark url as visited
	visited[currentURL] = true

	// Use scrape function to get HTML string
	htmlContent, err := Scrape(currentURL)
	if err != nil {
		return // Skip page if it can't be scraped
	}

	// Parse the links from web
	links, err := ParseLinks(htmlContent)
	if err != nil {
		return
	}

	for _, link := range links {
		absoluteURL := resolveURL(link, currentURL)
		if absoluteURL != "" {
			*allLinks = append(*allLinks, absoluteURL)

			if isInternalLink(absoluteURL, baseURL) {
				crawlRecursive(absoluteURL, baseURL, depth+1, maxDepth, visited, allLinks)
			}
		}
	}
}

func resolveURL(link, baseURL string) string {
	// Parse the current page url
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// Parse the link found on the page
	rel, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// Combine base URL with the relative link
	return base.ResolveReference(rel).String()
}

func isInternalLink(link, baseURL string) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	return linkURL.Hostname() == baseURLParsed.Hostname()
}
