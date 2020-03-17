package db

import (
	"github.com/garyburd/redigo/redis"
)

var redisClient redis.Conn

// SetRedisClint initialize client
func SetRedisClint(c redis.Conn) {
	redisClient = c
	if redisClient == nil {
		redisClient = c // <--- NOT THREAD SAFE
	}
}

func RedisSet(key string, value string) error {
	_, err := redisClient.Do("SET", key, value)
	return err
}

func RedisGet(key string) (string, error) {
	return redis.String(redisClient.Do("GET", key))
}
