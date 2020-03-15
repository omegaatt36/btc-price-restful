package main

import (
	"BTC-price-restful/routes"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
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
