package helpers

import (
	"testing"
)

func TestGenerateUniqueIdentifier(t *testing.T) {
	originalURL := "https://www.infracloud.io/blogs/"
	expected := "2a3786"
	uniqueIdentifier := GenerateUniqueIdentifier(originalURL)
	if uniqueIdentifier != expected {
		t.Errorf("GenerateUniqueIdentifier(%s) = %s; want %s", originalURL, uniqueIdentifier, expected)
	}
}

func TestCodetoUrl(t *testing.T) {
	code := "2a3786"
	expected := "http://localhost:8080/2a3786"
	shortURL := CodetoUrl(code)
	if shortURL != expected {
		t.Errorf("CodetoUrl(%s) = %s; want %s", code, shortURL, expected)
	}
}

func TestExtractDomain(t *testing.T) {
	url := "https://www.infracloud.io/blogs/"
	expected := "infracloud.io"
	domain := ExtractDomain(url)
	if domain != expected {
		t.Errorf("ExtractDomain(%s) = %s; want %s", url, domain, expected)
	}
}
