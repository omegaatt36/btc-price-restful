package routes

import (
	"io"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"gopkg.in/go-playground/assert.v1"
)

// GetDefault get a index page for test.
func testRegister(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Register")
}

func TestGetDefault(t *testing.T) {
	r := gofight.New()

	r.GET("/").
		SetDebug(true). // trun on the debug mode.
		Run(NewRouter(), func(res gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "btc-price-restful", res.Body.String())
			assert.Equal(t, http.StatusOK, res.Code)
		})
}

func TestRegister(t *testing.T) {
	register("GET", "/test",
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Register")
		}, nil)
	r := gofight.New()

	r.GET("/test").
		SetDebug(true). // trun on the debug mode.
		Run(NewRouter(), func(res gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "Register", res.Body.String())
			assert.Equal(t, http.StatusOK, res.Code)
		})
}
