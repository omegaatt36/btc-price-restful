package controllers

import (
	"io"
	"net/http"
)

// GetDefault get a index page for test.
func GetDefault(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "btc-price-restful")
}
