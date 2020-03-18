package db

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestSetMongoClint(t *testing.T) {
	got, err := mongo.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer got.Disconnect(context.TODO())
	SetMongoClint(got)
	if !reflect.DeepEqual(got, mongoClient) {
		t.Errorf("SetClint() got = %v, want %v", got, mongoClient)
	}
}

func TestSetRedisClint(t *testing.T) {
	got := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err := got.Ping(ctx).Result()
	if err != nil {
		t.Fatal("Connect to redis error", err)
	}
	SetRedisClint(got)
	defer got.Close()
	if !reflect.DeepEqual(got, redisClient) {
		t.Errorf("SetClint() got = %v, want %v", got, redisClient)
	}
}
