package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
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
