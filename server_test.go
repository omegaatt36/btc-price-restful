package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoClientOptions(t *testing.T) {
	os.Setenv("profile", "prod")
	want := options.Client().ApplyURI("mongodb://mongodb:27017")
	got := mongoClientOptions()
	if !reflect.DeepEqual(got.GetURI(), want.GetURI()) {
		t.Errorf("mongoClientOptions() got = %v, want %v", got, want)
	}
}

func TestMongoClientOptionsNonProdProfile(t *testing.T) {
	os.Setenv("profile", "dev")
	want := options.Client().ApplyURI("mongodb://localhost:27017")

	got := mongoClientOptions()

	if !reflect.DeepEqual(got.GetURI(), want.GetURI()) {
		t.Errorf("mongoClientOptions() got = %v, want %v", got, want)
	}
}

func TestRedisClientOptions(t *testing.T) {
	os.Setenv("profile", "prod")
	want := &redis.Options{
		Addr: "redisdb:6379",
	}

	got := redisClientOptions()

	assert.Equal(t, got.Addr, want.Addr)
}

func TestRedisClientOptionsNonProdProfile(t *testing.T) {
	os.Setenv("profile", "dev")
	want := &redis.Options{
		Addr: "localhost:6379",
	}
	got := redisClientOptions()
	assert.Equal(t, got.Addr, want.Addr)
}
