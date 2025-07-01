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
		{"example.com", "http://example.com"},
		{"/relative/path", ""},
		{"../up/one", ""},
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
		{makeAnchor("example.com"), "http://example.com"},
		{makeAnchor("mailto:test@example.com"), ""},
		{makeAnchor("/relative/path"), ""},
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
