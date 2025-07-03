package internal

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// ParseLinks parses the HTML content and extracts all links.
func ParseLinks(htmlContent string) ([]string, error) {
	var links []string
	// Convert htmlString into a strings.reader
	htmlReader := strings.NewReader(htmlContent)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("HTML Error: %w", err)
	}
	// Call traverse nodes
	traverseNodes(doc, &links)
	return links, nil
}

// traverseNodes recursively traverses the HTML nodes to find links.
func traverseNodes(n *html.Node, links *[]string) {
	// Check if the current node is an anchor tage
	if n.Type == html.ElementNode && n.Data == "a" {
		// Extract the href
		href := extractHref(n)
		// If it is valid (string will not be empty), then append to array
		if href != "" {
			*links = append(*links, href)
		}
	}
	// Recurse on child and sibling nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNodes(c, links)
	}

}

// extractHref extracts the href attribute from an anchor tag.
func extractHref(node *html.Node) string {
	// Check if node is an HTML node and that it is an anchor tag
	if node.Type == html.ElementNode && node.Data == "a" {
		// Iterate through the attributes of the node to find href
		for _, a := range node.Attr {
			// If the attribute is href, return its value
			if a.Key == "href" {
				// clean the URL to ensure it has a scheme
				return cleanURL(a.Val)
			}
		}
	}
	// If no href attribute is found, return an empty string
	return ""

}

func cleanURL(url string) string {
	// Trim whitespace from the URL
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}

	// Filter out invalid schemes
	if strings.HasPrefix(url, "javascript:") ||
		strings.HasPrefix(url, "mailto:") ||
		strings.HasPrefix(url, "tel:") {
		return ""
	}

	// Ensure the URL starts with a valid scheme
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") &&
		!strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "../") {
		url = "http://" + url
	}

	return url
}
