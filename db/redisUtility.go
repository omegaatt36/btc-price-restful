package db

import (
	"strings"

	"github.com/garyburd/redigo/redis"
)

// NSLatestAPI is a namespace name to get lastest api
const NSLatestAPI = "latestAPI"

var redisClient redis.Conn

// SetRedisClint initialize client
func SetRedisClint(c redis.Conn) {
	redisClient = c
	if redisClient == nil {
		redisClient = c // <--- NOT THREAD SAFE
	}
}

// RedisSet set key-value into redis
func RedisSet(key string, value string) error {
	_, err := redisClient.Do("SET", key, value)
	return err
}

// RedisGet get value from redis by key
func RedisGet(key string) (string, error) {
	return redis.String(redisClient.Do("GET", key))
}

// RedisKeysByNameSpace get all keys with namespace
func RedisKeysByNameSpace(nameSpace string) (keys []string, err error) {
	keys, err = redis.Strings(redisClient.Do("KEYS", nameSpace+":*"))
	if err != nil {
		return nil, err
	}
	for i, v := range keys {
		keys[i] = strings.Replace(v, nameSpace+":", "", 1)
	}
	return keys, nil
}
