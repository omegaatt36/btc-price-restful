package db

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// NSLatestAPI is a namespace name to get lastest api
const NSLatestAPI = "latestAPI"

// NSUserToken is a namespace name to get caching JWT
const NSUserToken = "userToken"

// NSUserQueryTimes is a namespace name to get caching user query times
const NSUserQueryTimes = "userQueryTimes"

var redisClient *redis.Client

// SetRedisClint initialize client
func SetRedisClint(c *redis.Client) {
	redisClient = c
	if redisClient == nil {
		redisClient = c // <--- NOT THREAD SAFE
	}
}

// RedisSet set key-value into redis
func RedisSet(key string, value interface{}) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	err := redisClient.Set(ctx, key, value, 0).Err()
	return err
}

// RedisIncr increase value by one
func RedisIncr(key string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	err := redisClient.Incr(ctx, key).Err()
	return err
}

// RedisGet get value from redis by key
func RedisGet(key string) (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return redisClient.Get(ctx, key).Result()
}

// RedisKeysByNameSpace get all keys with namespace
func RedisKeysByNameSpace(nameSpace string) (keys []string, err error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	keys, err = redisClient.Keys(ctx, nameSpace+":*").Result()
	if err != nil {
		return nil, err
	}
	for i, v := range keys {
		keys[i] = strings.Replace(v, nameSpace+":", "", 1)
	}
	return keys, nil
}
