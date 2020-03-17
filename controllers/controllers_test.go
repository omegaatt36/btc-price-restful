package controllers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetDefault(t *testing.T) {
	r := gofight.New()
	router := mux.NewRouter()
	router.HandleFunc("/", GetDefault).Methods("GET")
	r.GET("/").
		SetDebug(true). // trun on the debug mode.
		Run(router, func(res gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "btc-price-restful", res.Body.String())
			assert.Equal(t, http.StatusOK, res.Code)
		})
}
