package internal

import (
	"testing"
)

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		baseURL  string
		expected string
	}{
		{
			name:     "relative path",
			link:     "/about",
			baseURL:  "https://example.com",
			expected: "https://example.com/about",
		},
		{
			name:     "relative path with subdirectory",
			link:     "contact",
			baseURL:  "https://example.com/products/",
			expected: "https://example.com/products/contact",
		},
		{
			name:     "absolute URL",
			link:     "https://other.com/page",
			baseURL:  "https://example.com",
			expected: "https://other.com/page",
		},
		{
			name:     "fragment only",
			link:     "#section",
			baseURL:  "https://example.com/page",
			expected: "https://example.com/page#section",
		},
		{
			name:     "query parameter",
			link:     "?search=test",
			baseURL:  "https://example.com/page",
			expected: "https://example.com/page?search=test",
		},
		{
			name:     "parent directory",
			link:     "../parent",
			baseURL:  "https://example.com/sub/page",
			expected: "https://example.com/parent",
		},
		{
			name:     "invalid base URL",
			link:     "/about",
			baseURL:  "ht tp://invalid url with spaces",
			expected: "",
		},
		{
			name:     "invalid link",
			link:     "ht tp://invalid-url",
			baseURL:  "https://example.com",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveURL(tt.link, tt.baseURL)
			if result != tt.expected {
				t.Errorf("resolveURL(%q, %q) = %q, expected %q", tt.link, tt.baseURL, result, tt.expected)
			}
		})
	}
}

func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		baseURL  string
		expected bool
	}{
		{
			name:     "same domain",
			link:     "https://example.com/about",
			baseURL:  "https://example.com",
			expected: true,
		},
		{
			name:     "same domain different path",
			link:     "https://example.com/products/item",
			baseURL:  "https://example.com/home",
			expected: true,
		},
		{
			name:     "same domain with www",
			link:     "https://www.example.com/about",
			baseURL:  "https://www.example.com",
			expected: true,
		},
		{
			name:     "different domain",
			link:     "https://other.com/page",
			baseURL:  "https://example.com",
			expected: false,
		},
		{
			name:     "subdomain",
			link:     "https://sub.example.com/page",
			baseURL:  "https://example.com",
			expected: false,
		},
		{
			name:     "different protocol same domain",
			link:     "http://example.com/page",
			baseURL:  "https://example.com",
			expected: true,
		},
		{
			name:     "different port same domain",
			link:     "https://example.com:8080/page",
			baseURL:  "https://example.com",
			expected: true,
		},
		{
			name:     "invalid link URL",
			link:     "ht tp://invalid-url",
			baseURL:  "https://example.com",
			expected: false,
		},
		{
			name:     "invalid base URL",
			link:     "https://example.com/page",
			baseURL:  "not-a-valid-url",
			expected: false,
		},
		{
			name:     "both invalid URLs",
			link:     "ht tp://invalid",
			baseURL:  "not-a-valid-url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInternalLink(tt.link, tt.baseURL)
			if result != tt.expected {
				t.Errorf("isInternalLink(%q, %q) = %v, expected %v", tt.link, tt.baseURL, result, tt.expected)
			}
		})
	}
}

func TestCrawlRecursive(t *testing.T) {
	// Note: These tests would require mocking the Scrape and ParseLinks functions
	// For now, we'll test the basic structure and edge cases

	t.Run("max depth exceeded", func(t *testing.T) {
		visited := make(map[string]bool)
		var allLinks []string

		// Call with depth > maxDepth
		crawlRecursive("https://example.com", "https://example.com", 2, 1, visited, &allLinks)

		// Should not add any links since depth exceeds maxDepth
		if len(allLinks) != 0 {
			t.Errorf("Expected no links when depth exceeds maxDepth, got %d links", len(allLinks))
		}

		// Should not mark URL as visited
		if visited["https://example.com"] {
			t.Errorf("Expected URL not to be marked as visited when depth exceeds maxDepth")
		}
	})

	t.Run("already visited URL", func(t *testing.T) {
		visited := make(map[string]bool)
		visited["https://example.com"] = true
		var allLinks []string

		crawlRecursive("https://example.com", "https://example.com", 0, 2, visited, &allLinks)

		// Should not process already visited URL
		if len(allLinks) != 0 {
			t.Errorf("Expected no links when URL already visited, got %d links", len(allLinks))
		}
	})
}

func TestCrawlSite(t *testing.T) {
	t.Run("basic structure", func(t *testing.T) {
		// This test verifies the basic structure of CrawlSite
		// In a real scenario, you'd want to mock the HTTP calls

		startURL := "https://example.com"
		maxDepth := 1

		result := CrawlSite(startURL, maxDepth)

		// Result should be a slice (even if empty due to network issues in test)
		if result == nil {
			t.Errorf("Expected non-nil result from CrawlSite")
		}
	})

	t.Run("invalid start URL", func(t *testing.T) {
		startURL := "not-a-valid-url"
		maxDepth := 1

		result := CrawlSite(startURL, maxDepth)

		// Should return empty slice for invalid URL
		if len(result) != 0 {
			t.Errorf("Expected empty result for invalid start URL, got %d links", len(result))
		}
	})

	t.Run("zero max depth", func(t *testing.T) {
		startURL := "https://example.com"
		maxDepth := 0

		result := CrawlSite(startURL, maxDepth)

		// Should still process the start URL at depth 0
		if result == nil {
			t.Errorf("Expected non-nil result even with maxDepth 0")
		}
	})

	t.Run("negative max depth", func(t *testing.T) {
		startURL := "https://example.com"
		maxDepth := -1

		result := CrawlSite(startURL, maxDepth)

		// Should return empty slice since depth 0 > maxDepth -1
		if len(result) != 0 {
			t.Errorf("Expected empty result for negative maxDepth, got %d links", len(result))
		}
	})
}

// Benchmark tests
func BenchmarkResolveURL(b *testing.B) {
	link := "/about/contact"
	baseURL := "https://example.com/products/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveURL(link, baseURL)
	}
}

func BenchmarkIsInternalLink(b *testing.B) {
	link := "https://example.com/about"
	baseURL := "https://example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isInternalLink(link, baseURL)
	}
}
