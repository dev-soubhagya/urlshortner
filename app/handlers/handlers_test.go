package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dev-soubhagya/urlshortner/storage"
)

func TestShortenURL(t *testing.T) {
	// Initialize a test Shortener
	testShortener := storage.NewShortener("localhost:6379")
	h := NewHandler(testShortener)

	// Test case payload
	payload := map[string]string{"url": "https://www.infracloud.io/blogs/"}

	// Encode payload to JSON
	reqBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a test request
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	h.ShortenURL(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, rr.Code)
	}
}

func TestRedirect(t *testing.T) {
	// Initialize a test Shortener
	testShortener := storage.NewShortener("localhost:6379")
	h := NewHandler(testShortener)

	// Create a test request
	req, err := http.NewRequest("GET", "/shorturlcode", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	h.Redirect(rr, req)

	// Check the status code
	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status code %d; got %d", http.StatusMovedPermanently, rr.Code)
	}
}

func TestMetrics(t *testing.T) {
	// Initialize a test Shortener
	testShortener := storage.NewShortener("localhost:6379")
	h := NewHandler(testShortener)

	// Create a test request
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	h.Metrics(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, rr.Code)
	}
}
