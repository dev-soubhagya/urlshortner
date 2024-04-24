package storage

import (
	"fmt"
	"time"

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
				return redis.Dial("tcp", redisAddr)
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
