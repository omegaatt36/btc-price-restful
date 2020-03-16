package main

import (
	"BTC-price-restful/routes"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetDefault(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		SetDebug(true). // trun on the debug mode.
		Run(routes.NewRouter(), func(res gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "btc-price-restful", res.Body.String())
			assert.Equal(t, http.StatusOK, res.Code)
		})
}

func TestClientOptions(t *testing.T) {
	os.Setenv("profile", "prod")
	want := options.Client().ApplyURI("mongodb://db:27017")
	got := clientOptions()
	if !reflect.DeepEqual(got.GetURI(), want.GetURI()) {
		t.Errorf("clientOptions() got = %v, want %v", got, want)
	}
}

func TestClientOptionsNonProdProfile(t *testing.T) {
	os.Setenv("profile", "dev")
	want := options.Client().ApplyURI("mongodb://localhost:27017")

	got := clientOptions()

	if !reflect.DeepEqual(got.GetURI(), want.GetURI()) {
		t.Errorf("clientOptions() got = %v, want %v", got, want)
	}
}
