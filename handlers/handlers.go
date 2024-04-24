package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

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
	//extract domain from url
	domain := helpers.ExtractDomain(req.URL)
	//increment domain counter
	err = h.Shortener.IncrementCounter(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// to get metrics of top 3 shortend domain
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	// Retrieve keys matching a pattern
	keys, err := h.Shortener.GetkeysByPattern(helpers.KEY_PATTER)
	if err != nil {
		fmt.Println("Error retrieving keys from Redis:", err)
		return
	}
	fmt.Println(keys)
	// Retrieve counter values for each key
	counterValues := h.Shortener.GetKeysCounter(keys)
	fmt.Println("keys counter values: ", counterValues)
	// Sort keys by counter value in descending order
	type keyValue struct {
		Key   string
		Value int
	}
	var sortedKeys []keyValue
	for key, value := range counterValues {
		sortedKeys = append(sortedKeys, keyValue{Key: key, Value: value})
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].Value > sortedKeys[j].Value
	})

	fmt.Println("sorted keys :", sortedKeys)
	// Print top 3 keys with highest counter value
	fmt.Println("Top 3 domains by counter value:")
	for i, kv := range sortedKeys {
		if i >= 3 {
			break
		}
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
		fmt.Fprintf(w, "%s: %d\n", kv.Key, kv.Value)
	}
}
