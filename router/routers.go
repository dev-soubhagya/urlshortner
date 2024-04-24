package router

import (
	"net/http"

	"github.com/dev-soubhagya/urlshortner/handlers"
	"github.com/dev-soubhagya/urlshortner/storage"
)

func Router(shortener *storage.Shortener) {
	h := handlers.NewHandler(shortener)
	// Define HTTP routes
	http.HandleFunc("/shorten", h.ShortenURL)
}
