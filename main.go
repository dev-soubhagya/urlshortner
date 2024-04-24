package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dev-soubhagya/urlshortner/router"
	"github.com/dev-soubhagya/urlshortner/storage"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("App running")
	shortener := storage.NewShortener(os.Getenv("REDIS_ADDRESS"))
	router.Routers(shortener)
	log.Println("Server listening on port 8080...")
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil)
}
