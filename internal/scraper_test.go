package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScrape_Success(t *testing.T) {
	expected := "<html><body>Hello, world!</body></html>"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer ts.Close()

	body, err := Scrape(ts.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if body != expected {
		t.Errorf("expected %q, got %q", expected, body)
	}
}

func TestScrape_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := Scrape(ts.URL)
	if err == nil || !strings.Contains(err.Error(), "status code") {
		t.Errorf("expected HTTP status code error, got %v", err)
	}
}

func TestScrape_BadURL(t *testing.T) {
	_, err := Scrape(":bad-url:")
	if err == nil {
		t.Error("expected error for bad URL, got nil")
	}
}

func TestScrape_ReadBodyError(t *testing.T) {
	// Custom server that sends headers but then closes connection during body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100") // Claim we have 100 bytes
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("partial")) // Write only partial content
		// Connection will be closed when handler returns, causing read error
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		// Hijack and close to simulate connection drop during body read
		if hijacker, ok := w.(http.Hijacker); ok {
			conn, _, _ := hijacker.Hijack()
			conn.Close()
		}
	}))
	defer ts.Close()

	_, err := Scrape(ts.URL)
	if err == nil {
		t.Error("expected error for connection drop during body read, got nil")
	}
	// The error might be either a body reading error or HTTP error depending on timing
	// Accept either as both are valid for this scenario
	if !strings.Contains(err.Error(), "reading body") && !strings.Contains(err.Error(), "HTTP Error") {
		t.Errorf("expected body read error or HTTP error, got %v", err)
	}
}
