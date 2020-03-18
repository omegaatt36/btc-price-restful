package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"BTC-price-restful/auth"
	"BTC-price-restful/db"
	"BTC-price-restful/models"
	"BTC-price-restful/utility"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const queryLimit = 500

func isUserExist(userName string) bool {
	if singleResult := db.FindOne(db.CollectionUser, models.User{UserName: userName}); singleResult.Err() != nil {
		return false
	}
	return true
}

func checkQueryLimit(userName string) bool {
	key := db.NSUserQueryTimes + ":" + userName
	queryTimes, err := db.RedisGet(key)
	if err != nil {
		db.RedisSet(key, 0)
		return true
	}
	qt, _ := strconv.Atoi(queryTimes)
	if qt >= 500 {
		return false
	}
	return true
}

func increaseQueryLimit(userName string) {
	key := db.NSUserQueryTimes + ":" + userName
	err := db.RedisIncr(key)
	if err != nil {
		logrus.Infof("%s query can't increase", userName)
	}

}

// Register name & cryp(pwd) into DB
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	body, _ := ioutil.ReadAll(r.Body)
	logrus.Info(string(body))
	if err != nil || user.UserName == "" || user.Password == "" {
		utility.ResponseWithJSON(w, http.StatusBadRequest, utility.Response{Message: "bad param", Result: utility.ResFailed})
		return
	}
	if isUserExist(user.UserName) {
		utility.ResponseWithJSON(w, http.StatusBadRequest, utility.Response{Message: "user exist", Result: utility.ResFailed})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utility.ResponseWithJSON(w, http.StatusInternalServerError, utility.Response{Message: "password parse error", Result: utility.ResFailed})
		return
	}
	user.Password = string(hash)
	_, err = db.Create(db.CollectionUser, user)
	if err != nil {
		utility.ResponseWithJSON(w, http.StatusInternalServerError, utility.Response{Message: "DB insert error", Result: utility.ResFailed})
		return
	}
	utility.ResponseWithJSON(w, http.StatusOK, utility.Response{Result: utility.ResSuccess})
}

// Login verify name & pwd and return token
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserName == "" || user.Password == "" {
		utility.ResponseWithJSON(w, http.StatusBadRequest, utility.Response{Message: "bad param", Result: utility.ResFailed})
		return
	}
	singleResult := db.FindOne(db.CollectionUser, models.User{UserName: user.UserName})
	if singleResult.Err() != nil {
		utility.ResponseWithJSON(w, http.StatusBadRequest, utility.Response{Message: "user not exist", Result: utility.ResFailed})
		return
	}
	var userInDB models.User
	err = singleResult.Decode(&userInDB)
	if err != nil {
		utility.ResponseWithJSON(w, http.StatusInternalServerError, utility.Response{Message: "DB decode error", Result: utility.ResFailed})
		return
	}
	eq := bcrypt.CompareHashAndPassword([]byte(userInDB.Password), []byte(user.Password))
	if eq != nil {
		utility.ResponseWithJSON(w, http.StatusInternalServerError, utility.Response{Message: "password error", Result: utility.ResFailed})
		return
	}
	token, _ := auth.GenerateToken(&user)
	utility.ResponseWithJSON(w, http.StatusOK, utility.Response{Result: utility.ResSuccess, Data: models.JwtToken{Token: token}})
}
