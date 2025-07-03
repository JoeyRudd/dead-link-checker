package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDeadLinks(t *testing.T) {

	// create test servers
	server200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted) // Good server
	}))
	defer server200.Close()

	server301 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMovedPermanently)
	}))
	defer server301.Close()

	server404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server404.Close()

	server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server500.Close()

	// test urls
	testURLs := []string{
		server200.URL,
		server301.URL,
		server404.URL,
		server500.URL,
	}

	// expected result
	expectedURls := []string{
		server301.URL,
		server404.URL,
		server500.URL,
	}

	resultURls := GetDeadLinks(&testURLs)

	// check they're of equal lenth
	if len(resultURls) != 3 {
		t.Errorf("expected 3 dead links, got %d", len(resultURls))
	}

	// compare result to expected result
	for i, resultURL := range resultURls {
		if resultURL != expectedURls[i] {
			t.Errorf("expected URL '%s', got '%s'", expectedURls[i], resultURL)
		}
	}

}
