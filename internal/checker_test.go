package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDeadLinks(t *testing.T) {

	// Accepted server
	server200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted) // Good server
	}))
	defer server200.Close()

	// Replace server301 with this
	server301 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server200.URL, http.StatusMovedPermanently)
	}))
	defer server301.Close()

	// Not found
	server404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server404.Close()

	// Server error
	server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server500.Close()

	// Add this after your existing servers
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server200.URL, http.StatusMovedPermanently)
	}))
	defer redirectServer.Close()

	// test urls
	testURLs := []string{
		server200.URL,
		server301.URL,
		server404.URL,
		server500.URL,
		redirectServer.URL,
	}

	// expected result - using map for order-independent comparison
	expectedDeadURLs := map[string]bool{
		server404.URL: true,
		server500.URL: true,
	}

	resultURLs := GetDeadLinks(&testURLs)

	// check they're of equal length
	if len(resultURLs) != 2 {
		t.Errorf("expected 2 dead links, got %d", len(resultURLs))
	}

	// compare result to expected result (order independent)
	for _, resultURL := range resultURLs {
		if !expectedDeadURLs[resultURL] {
			t.Errorf("unexpected dead URL '%s'", resultURL)
		}
	}

	// check that all expected URLs are present
	resultMap := make(map[string]bool)
	for _, url := range resultURLs {
		resultMap[url] = true
	}

	for expectedURL := range expectedDeadURLs {
		if !resultMap[expectedURL] {
			t.Errorf("expected dead URL '%s' not found in results", expectedURL)
		}
	}
}
