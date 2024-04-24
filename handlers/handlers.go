package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dev-soubhagya/urlshortner/helpers"
	"github.com/dev-soubhagya/urlshortner/storage"
	"github.com/dev-soubhagya/urlshortner/utils"
	"github.com/gomodule/redigo/redis"
)

type Handler struct {
	Shortener *storage.Shortener
}

func NewHandler(shortener *storage.Shortener) *Handler {
	return &Handler{
		Shortener: shortener,
	}
}
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req utils.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uniquehash := helpers.GenerateUniqueIdentifier(req.URL)
	fmt.Println("Random code:", uniquehash)
	shortURL := helpers.CodetoUrl(uniquehash)
	fmt.Println("short Url:", shortURL)
	// Check if URL is already shortened in Redis
	longURL, err := h.Shortener.Get(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if longURL == "" {
		// Store short URL and original URL in storage
		err := h.Shortener.Set(shortURL, req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"shortened_url": shortURL})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	geturlcode := r.URL.Path
	geturlcode = geturlcode[1:]
	fmt.Println("short url to redirect", geturlcode)

	shortURL := helpers.CodetoUrl(geturlcode)

	// Retrieve original URL from Redis
	originalURL, err := h.Shortener.Get(shortURL)
	fmt.Println("original url from redis ", originalURL)
	if err == redis.ErrNil {
		fmt.Println("not found", err)
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
