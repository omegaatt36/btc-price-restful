package db

import (
	"context"
	"log"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestSetClint(t *testing.T) {
	got, err := mongo.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer got.Disconnect(context.TODO())
	SetClint(got)
	if !reflect.DeepEqual(got, client) {
		t.Errorf("SetClint() got = %v, want %v", got, client)
	}
}
