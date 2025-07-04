package internal

import (
	"golang.org/x/net/html"
	"testing"
)

func TestCleanURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" http://example.com ", "http://example.com"},
		{"https://example.com", "https://example.com"},
		{"example.com", "example.com"},       // Allow relative paths without protocol
		{"/relative/path", "/relative/path"}, // Allow relative paths for site traversal
		{"../up/one", "../up/one"},           // Allow parent directory navigation
		{"mailto:test@example.com", ""},
		{"javascript:alert('x')", ""},
		{"tel:1234567890", ""},
		{"", ""},
	}
	for _, tt := range tests {
		result := cleanURL(tt.input)
		if result != tt.expected {
			t.Errorf("cleanURL(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractHref(t *testing.T) {
	// Helper to create an anchor node with given href
	makeAnchor := func(href string) *html.Node {
		return &html.Node{
			Type: html.ElementNode,
			Data: "a",
			Attr: []html.Attribute{{Key: "href", Val: href}},
		}
	}

	tests := []struct {
		node     *html.Node
		expected string
	}{
		{makeAnchor("http://example.com"), "http://example.com"},
		{makeAnchor("example.com"), "example.com"}, // Allow relative paths without protocol
		{makeAnchor("mailto:test@example.com"), ""},
		{makeAnchor("/relative/path"), "/relative/path"}, // Allow relative paths for site traversal
		{makeAnchor(""), ""},
		{&html.Node{Type: html.ElementNode, Data: "div"}, ""}, // not an anchor
	}

	for _, tt := range tests {
		result := extractHref(tt.node)
		if result != tt.expected {
			t.Errorf("extractHref() = %q; want %q", result, tt.expected)
		}
	}

}

func TestParseLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "single link",
			html:     `<a href="http://example.com">Link</a>`,
			expected: []string{"http://example.com"},
		},
		{
			name:     "multiple links",
			html:     `<a href="http://example.com">Link1</a><a href="https://test.com">Link2</a>`,
			expected: []string{"http://example.com", "https://test.com"},
		},
		{
			name:     "nested links in divs",
			html:     `<div><a href="example.com">Link</a></div><p><a href="https://test.com">Test</a></p>`,
			expected: []string{"example.com", "https://test.com"}, // Allow relative paths
		},
		{
			name:     "no links",
			html:     `<div>No links here</div>`,
			expected: []string{},
		},
		{
			name:     "invalid links filtered out",
			html:     `<a href="mailto:test@example.com">Email</a><a href="/relative">Relative</a><a href="http://valid.com">Valid</a>`,
			expected: []string{"/relative", "http://valid.com"}, // Include relative paths for site traversal
		},
		{
			name:     "empty html",
			html:     "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLinks(tt.html)
			if err != nil {
				t.Fatalf("ParseLinks() error = %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("ParseLinks() got %d links, want %d", len(result), len(tt.expected))
			}
			for i, link := range result {
				if i >= len(tt.expected) || link != tt.expected[i] {
					t.Errorf("ParseLinks() = %v, want %v", result, tt.expected)
					break
				}
			}
		})
	}
}

func TestTraverseNodes(t *testing.T) {
	// Helper to create HTML structure
	createNode := func(nodeType html.NodeType, data string, attrs []html.Attribute) *html.Node {
		return &html.Node{
			Type: nodeType,
			Data: data,
			Attr: attrs,
		}
	}

	// Create a simple HTML structure: <div><a href="http://example.com">Link</a></div>
	anchor := createNode(html.ElementNode, "a", []html.Attribute{{Key: "href", Val: "http://example.com"}})
	div := createNode(html.ElementNode, "div", nil)
	div.FirstChild = anchor
	anchor.Parent = div

	var links []string
	traverseNodes(div, &links)

	expected := []string{"http://example.com"}
	if len(links) != len(expected) {
		t.Errorf("traverseNodes() got %d links, want %d", len(links), len(expected))
	}
	if len(links) > 0 && links[0] != expected[0] {
		t.Errorf("traverseNodes() = %v, want %v", links, expected)
	}
}
