package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func GenerateUniqueIdentifier(originalURL string) string {
	// Create a new SHA-256 hash instance
	hasher := sha256.New()

	// Write the original URL bytes to the hasher
	hasher.Write([]byte(originalURL))

	// Calculate the hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash sum to a hexadecimal string
	hashString := hex.EncodeToString(hashSum)

	// Return the hexadecimal string as the unique identifier
	return hashString[:6]
}

func CodetoUrl(code string) string {
	shortCode := fmt.Sprintf("/%s", code)
	shortURL := "http://" + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT") + shortCode
	return shortURL
}

// extractDomain extracts the domain from a given URL.
func ExtractDomain(urlLikeString string) string {

	urlLikeString = strings.TrimSpace(urlLikeString)

	if regexp.MustCompile(`^https?`).MatchString(urlLikeString) {
		read, _ := url.Parse(urlLikeString)
		urlLikeString = read.Host
	}

	if regexp.MustCompile(`^www\.`).MatchString(urlLikeString) {
		urlLikeString = regexp.MustCompile(`^www\.`).ReplaceAllString(urlLikeString, "")
	}

	return regexp.MustCompile(`([a-z0-9\-]+\.)+[a-z0-9\-]+`).FindString(urlLikeString)
}
