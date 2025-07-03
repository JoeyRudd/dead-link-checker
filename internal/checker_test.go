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

	// Redirect
	server301 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMovedPermanently)
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

	// expected result
	expectedURls := []string{
		server404.URL,
		server500.URL,
		redirectServer.URL,
	}

	resultURls := GetDeadLinks(&testURLs)

	// check they're of equal lenth
	if len(resultURls) != 2 {
		t.Errorf("expected 2 dead links, got %d", len(resultURls))
	}

	// compare result to expected result
	for i, resultURL := range resultURls {
		if resultURL != expectedURls[i] {
			t.Errorf("expected URL '%s', got '%s'", expectedURls[i], resultURL)
		}
	}

}
