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
	// Custom server that closes connection immediately
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer ts.Close()

	_, err := Scrape(ts.URL)
	if err == nil || !strings.Contains(err.Error(), "reading body") {
		t.Errorf("expected body read error, got %v", err)
	}
}
