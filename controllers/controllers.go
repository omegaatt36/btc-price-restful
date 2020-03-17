package controllers

import (
	"BTC-price-restful/db"
	"BTC-price-restful/models"
	"BTC-price-restful/utility"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetDefault get a index page for test.
func GetDefault(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "btc-price-restful")
}

// GetServiceMap response all active service
func GetServiceMap(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("userName")
	if !checkQueryLimit(userName) {
		utility.ResponseWithJSON(w, http.StatusForbidden, utility.Response{Message: "over the query limit", Result: utility.ResFailed})
		return
	}
	keys, err := db.RedisKeysByNameSpace(db.NSLatestAPI)
	if err != nil {
		utility.ResponseWithJSON(w, http.StatusInternalServerError, utility.Response{Message: "no active service", Result: utility.ResFailed})
		return
	}
	utility.ResponseWithJSON(w, http.StatusOK, utility.Response{Result: utility.ResSuccess, Data: keys})
	increaseQueryLimit(userName)
}

// GetLatestPrice response latest price which is active from db/redis
func GetLatestPrice(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("userName")
	if !checkQueryLimit(userName) {
		utility.ResponseWithJSON(w, http.StatusForbidden, utility.Response{Message: "over the query limit", Result: utility.ResFailed})
		return
	}
	result := make(map[string]models.Price)
	services := strings.Split(mux.Vars(r)["service"], ",")
	for _, service := range services {
		value, err := db.RedisGet(db.NSLatestAPI + ":" + service)
		if err != nil {
			logrus.Infof("redis get error %s", service)
			continue
		}
		var price models.Price
		err = json.Unmarshal([]byte(value.(string)), &price)
		if err != nil {
			logrus.Infof("json decode error %s", value)
			continue
		}
		if _, ok := result[service]; !ok {
			result[service] = price
		}
	}
	if len(result) == 0 {
		utility.ResponseWithJSON(w, http.StatusBadRequest, utility.Response{Message: "no active service or service not exsit", Result: utility.ResFailed})
		return
	}
	utility.ResponseWithJSON(w, http.StatusOK, utility.Response{Result: utility.ResSuccess, Data: result})
	increaseQueryLimit(userName)
}
