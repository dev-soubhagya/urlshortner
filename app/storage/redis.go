package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/dev-soubhagya/urlshortner/helpers"
	"github.com/gomodule/redigo/redis"
)

// Shortener represents a URL shortener service with Redis caching.
type Shortener struct {
	pool *redis.Pool
}

func NewShortener(redisAddr string) *Shortener {
	return &Shortener{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", redisAddr)
				if err != nil {
					fmt.Println("redis init dial err :", err)
					return nil, err
				} else {
					return conn, err
				}
			},
		},
	}
}

func (s *Shortener) Get(shortURL string) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	longURL, err := redis.String(conn.Do("GET", shortURL))
	if err != nil && err != redis.ErrNil {
		return "", fmt.Errorf("error getting URL from Redis: %v", err)
	}

	return longURL, nil
}

func (s *Shortener) Set(shortURL, longURL string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", shortURL, longURL)
	if err != nil {
		return fmt.Errorf("error setting URL in Redis: %v", err)
	}

	return nil
}

func (s *Shortener) IncrementCounter(domain string) error {
	conn := s.pool.Get()
	defer conn.Close()
	_, err := conn.Do("INCR", helpers.KEY_PATTER+domain)
	if err != nil {
		return fmt.Errorf("error incrementing domain counter in Redis: %v", err)
	}
	return nil
}

func (s *Shortener) GetkeysByPattern(patter string) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()
	// Retrieve keys matching a pattern
	keys, err := redis.Strings(conn.Do("KEYS", patter+"*"))
	if err != nil {
		fmt.Println("Error retrieving keys from Redis:", err)
		return nil, err
	}
	return keys, nil
}

func (s *Shortener) GetKeysCounter(keys []string) map[string]int {
	conn := s.pool.Get()
	defer conn.Close()

	// Retrieve counter values for each key
	counterValues := make(map[string]int)
	for _, key := range keys {
		value, err := redis.Int(conn.Do("GET", key))
		if err != nil {
			fmt.Printf("Error retrieving value for key %s: %v\n", key, err)
			continue
		}
		domainkey := strings.Split(key, ":")
		counterValues[domainkey[1]] = value
	}
	return counterValues
}
