package storage

import (
	"os"
	"testing"
)

var testShortener *Shortener

func TestMain(m *testing.M) {
	// Initialize the test Shortener with a Redis instance running on localhost
	testShortener = NewShortener("localhost:6379")

	// Run the tests
	exitVal := m.Run()

	// Clean up resources
	os.Exit(exitVal)
}

func TestShortener_SetAndGet(t *testing.T) {
	shortURL := "http://localhost:8080/2a3786"
	longURL := "https://www.infracloud.io/blogs/"
	err := testShortener.Set(shortURL, longURL)
	if err != nil {
		t.Errorf("Error setting URL: %v", err)
	}

	retrievedURL, err := testShortener.Get(shortURL)
	if err != nil {
		t.Errorf("Error getting URL: %v", err)
	}

	if retrievedURL != longURL {
		t.Errorf("Expected URL %s, got %s", longURL, retrievedURL)
	}
}

func TestShortener_IncrementCounter(t *testing.T) {
	domain := "infracloud.io"
	err := testShortener.IncrementCounter(domain)
	if err != nil {
		t.Errorf("Error incrementing counter: %v", err)
	}
}

func TestShortener_GetkeysByPattern(t *testing.T) {
	pattern := "domain-counter:"
	_, err := testShortener.GetkeysByPattern(pattern)
	if err != nil {
		t.Errorf("Error getting keys by pattern: %v", err)
	}
}

func TestShortener_GetKeysCounter(t *testing.T) {
	keys := []string{"domain-counter:infracloud.io", "domain-counter:redis.io"} // Assuming these keys exist in the Redis instance
	_ = testShortener.GetKeysCounter(keys)
}
