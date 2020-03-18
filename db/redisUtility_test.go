package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type clientRedisMock struct {
	mock.Mock
	client *redis.Client
}

func newRedisMock(client *redis.Client) *clientRedisMock {
	m := new(clientRedisMock)
	m.client = client

	return m
}

func newTestRedis() *clientRedisMock {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return newRedisMock(client)
}

func TestSetRedisClint(t *testing.T) {
	r := newTestRedis()
	SetRedisClint(r.client)
	_, err := redisClient.Ping(context.TODO()).Result()
	assert.Nil(t, err)
}

func TestRedisSet(t *testing.T) {
	key, value := "test", "test"
	r := newTestRedis()
	SetRedisClint(r.client)
	err := RedisSet(key, value)
	want, err := redisClient.Get(context.TODO(), key).Result()
	assert.Equal(t, want, value)
	assert.Nil(t, err)
}

func TestRedisGet(t *testing.T) {
	key, value := "test", "test"
	r := newTestRedis()
	SetRedisClint(r.client)
	err := redisClient.Set(context.TODO(), key, value, 0).Err()
	assert.Nil(t, err)
	want, err := RedisGet(key)
	assert.Equal(t, want, value)
}

func TestRedisIncr(t *testing.T) {
	key, value := "test", int64(1)
	r := newTestRedis()
	SetRedisClint(r.client)
	err := redisClient.Set(context.TODO(), key, value, 0).Err()
	assert.Nil(t, err)
	err = RedisIncr(key)
	assert.Nil(t, err)
	want, err := redisClient.Get(context.TODO(), key).Result()
	assert.Nil(t, err)
	n, err := strconv.ParseInt(want, 10, 64)
	assert.Nil(t, err)
	assert.Equal(t, n, value+1)
}

func TestRedisKeysByNameSpace(t *testing.T) {
	ns := "ns_test"
	key, value := "123", "value"
	r := newTestRedis()
	SetRedisClint(r.client)
	err := redisClient.Set(context.TODO(), ns+":"+key, value, 0).Err()
	assert.Nil(t, err)
	want, err := RedisKeysByNameSpace(ns)
	assert.Equal(t, want, []string{key})
	assert.Nil(t, err)
}
