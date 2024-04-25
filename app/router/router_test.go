package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dev-soubhagya/urlshortner/handlers"
	"github.com/dev-soubhagya/urlshortner/storage"
)

func TestRouters(t *testing.T) {
	// Initialize a test Shortener
	testShortener := storage.NewShortener("localhost:6379")

	// Initialize the router with the test Shortener
	Routers(testShortener)

	// Define test cases for each route
	testCases := []struct {
		method     string
		url        string
		handler    http.HandlerFunc
		statusCode int
	}{
		//{method: "POST", url: "/shorten", handler: http.HandlerFunc(handlers.NewHandler(testShortener).ShortenURL), statusCode: http.StatusOK},
		{method: "GET", url: "/", handler: http.HandlerFunc(handlers.NewHandler(testShortener).Redirect), statusCode: http.StatusMovedPermanently},
		{method: "GET", url: "/metrics", handler: http.HandlerFunc(handlers.NewHandler(testShortener).Metrics), statusCode: http.StatusOK},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Create a test request
		req, err := http.NewRequest(tc.method, tc.url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request using the router
		http.DefaultServeMux.ServeHTTP(rr, req)

		// Check the status code
		if rr.Code != tc.statusCode {
			t.Errorf("Expected status code %d for %s %s; got %d", tc.statusCode, tc.method, tc.url, rr.Code)
		}

		// Check if the correct handler is called
		if rr.Body != nil && tc.handler != nil && rr.Body.String() != "" {
			expectedBody := httptest.NewRecorder()
			tc.handler.ServeHTTP(expectedBody, req)
			if rr.Body.String() != expectedBody.Body.String() {
				t.Errorf("Handler mismatch for %s %s", tc.method, tc.url)
			}
		}
	}
}
