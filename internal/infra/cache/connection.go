package cache

import (
	"github.com/redis/go-redis/v9"
)

func CacheConnection(addr string, password string, username string, db int) (*redis.Client, error) {
	connection := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		Username: username,
		DB:       db,
		Protocol: 2,
	})

	return connection, nil
}
