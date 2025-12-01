package cache

import (
	"github.com/redis/go-redis/v9"
)

func CacheConnection(connectionString string) (*redis.Client, error) {
	connection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		Username: "",
		DB:       0,
		Protocol: 2,
	})

	return connection, nil
}
